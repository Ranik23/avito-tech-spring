package repository

import (
	"context"

)


type Transaction interface {
    Commit(context.Context) error 
    Rollback(context.Context) error 
    Transaction() interface{}
}
