package gapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/pb"
	"github.com/sayedppqq/banking-backend/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid password")
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.GenerateToken(user.Username, user.Role, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.GenerateToken(user.Username, user.Role, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create refresh token")
	}

	meta := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID: pgtype.UUID{
			Bytes: refreshTokenPayload.ID,
			Valid: true,
		},
		Username:     user.Username,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		ExpiresAt: pgtype.Timestamptz{
			Time:  refreshTokenPayload.ExpiredAt,
			Valid: true,
		},
		UserAgent: meta.userAgent,
		ClientIp:  meta.clientIP,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create session")
	}

	resp := pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             uuidToString(session.ID),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}
	return &resp, nil
}

func uuidToString(uuid pgtype.UUID) string {
	uuidString := uuid.Bytes[:]
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuidString[0:4], uuidString[4:6], uuidString[6:8], uuidString[8:10], uuidString[10:])
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
