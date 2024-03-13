package db

import (
	"context"
	"fmt"
)

func (sqlStore *SQLStore) execTex(ctx context.Context, fn func(*Queries) error) error {
	tx, err := sqlStore.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	queries := New(tx)
	err = fn(queries)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transictation err: %v, roll back err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
