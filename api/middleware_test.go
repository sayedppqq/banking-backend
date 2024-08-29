package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sayedppqq/banking-backend/token"
	"github.com/sayedppqq/banking-backend/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func addAuthorization(t *testing.T, req *http.Request, tokenMaker token.Maker, authorizationType string,
	name, role string, duration time.Duration) {
	token, payload, err := tokenMaker.GenerateToken(name, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotNil(t, token)

	req.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationType, token))
}

func TestAuthMiddleware(t *testing.T) {
	username := util.RandomOwnerName(false)
	role := util.DepositorRole

	tests := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, username, role, time.Hour)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "unsupported", username, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "", username, role, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			server.router.GET(
				"/auth",
				authMiddleware(server.tokenMaker),
				func(context *gin.Context) {
					context.JSON(http.StatusOK, gin.H{})
				},
			)
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)

			test.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			test.checkResponse(t, recorder)
		})
	}
}
