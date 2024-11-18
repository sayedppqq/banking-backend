package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"time"
)

type renewAccessTokenParams struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var params renewAccessTokenParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(params.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(payload.ExpiredAt) {
		err := errors.New("token is expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	session, err := server.store.GetSession(ctx, pgtype.UUID{Bytes: payload.ID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != params.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := errors.New("session is locked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != payload.Username {
		err := errors.New("invalid user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt.Time) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.GenerateToken(payload.Username, payload.Role, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, response)
}
