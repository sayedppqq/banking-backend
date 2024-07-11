package api

import (
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func randomUser(t *testing.T) (db.User, string) {
	pass := util.RandomString(6)
	hashedPass, err := util.HashPassword(pass)
	require.NoError(t, err)
	user := db.User{
		Username:       util.RandomOwnerName(),
		HashedPassword: hashedPass,
		FullName:       util.RandomOwnerName(),
		Email:          util.RandomEmail(),
	}
	return user, hashedPass
}
