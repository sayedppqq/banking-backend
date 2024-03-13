package db

import (
	"context"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	GetTransfer(ctx context.Context, id int64) (Transfer, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
}

// (nil) is used to initialize a nil pointer of type *Queries.
// It's essentially a placeholder value for the pointer,
// indicating that we are not interested in a specific instance of Queries but just the type information.
var _ Querier = (*Queries)(nil)
