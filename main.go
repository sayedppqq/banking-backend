package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sayedppqq/banking-backend/api"
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/gapi"
	"github.com/sayedppqq/banking-backend/pb"
	"github.com/sayedppqq/banking-backend/util"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	conn, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal("can not create pg connection pool", err)
	}

	store := db.NewStore(conn)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runGrpcServer(ctx, waitGroup, store, config)
	runGatewayServer(ctx, waitGroup, store, config)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal("error from wait group")
	}
}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	store db.Store,
	config util.Config,
) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create grpc server", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterBankingBackendServer(grpcServer, server)

	listener, err := net.Listen("tcp", config.GrpcHostAddress)
	if err != nil {
		log.Fatal("grpc server listen failed", err)
	}

	waitGroup.Go(func() error {
		log.Println("starting grpc server at port", config.GrpcHostAddress)
		err := grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			return fmt.Errorf("grpc server start failed: %w", err)
		}
		return nil
	})
	waitGroup.Go(func() error {
		<-ctx.Done() // if ctrl+c is pressed
		log.Println("graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Println("gRPC server is stopped")

		return nil
	})
}

func runGatewayServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	store db.Store,
	config util.Config,
) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create gateway server: %w", err)
	}
	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	err = pb.RegisterBankingBackendHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register grpc gateway handler", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	//listener, err := net.Listen("tcp", config.HTTPHostAddress)
	//if err != nil {
	//	log.Fatal("http gateway server listen failed", err)
	//}

	httpServer := http.Server{
		Addr:    config.HTTPHostAddress,
		Handler: mux,
	}

	waitGroup.Go(func() error {
		log.Println("starting grpc gateway server at port", config.HTTPHostAddress)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return fmt.Errorf("grpc gateway server start failed: %w", err)
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done() // if ctrl+c is pressed
		log.Println("graceful shutdown grpc gateway server")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("failed to shutdown HTTP gateway server: %w", err)
		}
		log.Println("HTTP gateway server is stopped")

		return nil
	})
}

func runGinServer(store db.Store, config util.Config) error {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("can not setup new server", err)
	}

	err = server.Run(config.HTTPHostAddress)
	return err
}
