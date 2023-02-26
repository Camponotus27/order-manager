package repository

import (
	"context"
)

type Cookie interface {
	Current(ctx context.Context) (int, error)
	SellOne(ctx context.Context) (int, error)
}
