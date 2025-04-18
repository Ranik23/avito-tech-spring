package repository

import (
	"context"
)



type CtxManager interface {
    Default(context.Context) Transaction
    ByKey(context.Context, CtxKey) Transaction
	CtxKey() CtxKey
}
