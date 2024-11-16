package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-starter-kit/internal/pkg/database"
	"go.uber.org/zap"
)

func Tx(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := InitCtx(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)
		defer func() {
			EndCtx(ctx, logger)
		}()
		c.Next()
	}
}

func InitCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, database.TransactionCtxKey, database.TransactionCtx{})
}

func EndCtx(ctx context.Context, logger *zap.Logger) {
	if transactionCtx, ok := ctx.Value(database.TransactionCtxKey).(*database.TransactionCtx); ok {
		if tx := transactionCtx.Conn; tx != nil {
			if err := recover(); err != nil {
				if err := tx.Rollback(); err != nil {
					logger.Error(fmt.Sprintf("tx rollback failed: %s", err))
				} else {
					logger.Info("tx rollbacked")
				}
				panic(err)
			} else {
				if err = tx.Commit(); err != nil {
					logger.Error(fmt.Sprintf("commit transaction failed: %s", err))
				}
			}
		}
	}
}
