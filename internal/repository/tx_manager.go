package repository

import (
	"context"
)

type CtxKey struct{}

type TxManager interface {
	Do(context.Context, func(context.Context) error) error
}
