package internal

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spolia/lemon-wallet/internal/wallet/movement"
	"github.com/spolia/lemon-wallet/internal/wallet/user"
)

type Service interface {
	CreateUser(ctx context.Context, name, lastName, alias, email string) (int64, error)
	GetUser(ctx context.Context, id int64) (user.User, error)
	CreateMovement(ctx context.Context, movement movement.Movement) (int64, error)
	SearchMovement(ctx context.Context, userID int64, limit, offset uint64, movType, currencyName string) ([]movement.Row, error)
}

func API(router *gin.Engine, service Service) {
	router.POST("/users", createUser(service))
	router.GET("/users/:id", getUser(service))
	router.POST("/movements", createMovement(service))
	router.GET("/movements/search", searchMovement(service))
}
