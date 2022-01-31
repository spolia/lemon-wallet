package internal

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spolia/lemon-wallet/internal/wallet/movement"
	"github.com/spolia/lemon-wallet/internal/wallet/user"
)

type UserService interface {
	Create(ctx context.Context, name, lastName, alias, email string) error
	Get(ctx context.Context, id int64) (user.User, error)
}

type MovementService interface {
	Create(ctx context.Context, movement movement.Movement) (int64, error)
	Search(ctx context.Context, limit, offset uint64, movType, currencyName string, userID int64) ([]movement.Row, error)
}

func API(router *gin.Engine, userService UserService, movementService MovementService) {
	router.POST("/users", createUser(userService))
	router.GET("/users/:id", getUser(userService))
	router.POST("/movements", createMovement(movementService))
	router.GET("/movements/search", searchMovement(movementService))
}