package api

import (
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateUserAPI(t *testing.T) {

}

func randomUser(t *testing.T) (db.User, string) {
	pass := util.RandomString(6)
	hashedPass, err := util.HashPassword(pass)
	require.NoError(t, err)
	user := db.User{
		Username:       util.RandomOwnerName(false),
		HashedPassword: hashedPass,
		FullName:       util.RandomOwnerName(true),
		Email:          util.RandomEmail(),
	}
	return user, hashedPass
}
