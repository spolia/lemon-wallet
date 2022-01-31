package internal

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spolia/lemon-wallet/internal/wallet/movement"
	"github.com/spolia/lemon-wallet/internal/wallet/user"
)

/*
1-hacer el endpoint de search
2-testeart todos los casos y meter comentarios
3-agregar test en todos los niveles con testdata
5-levantar con docker
4-agregar readme

*/
func createUser(userService UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userRequest user.User
		if err := ctx.ShouldBindJSON(&userRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := userService.Create(ctx, userRequest.FirstName, userRequest.LastName, userRequest.Alias, userRequest.Email); err != nil {
			if err == user.ErrorAlreadyExist {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusCreated, "ok")
	}
}

func getUser(userService UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		userResult, err := userService.Get(ctx, userID)
		if err != nil {
			if err == user.ErrorUserNotFound {
				ctx.JSON(http.StatusNotFound, err.Error())
				return
			}

			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, userResult)
	}
}

func createMovement(movementService MovementService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var movementRequest movement.Movement
		if err := ctx.ShouldBindJSON(&movementRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		movementID, err := movementService.Create(ctx, movementRequest)
		if err != nil {
			if err == movement.ErrorWrongCurrency || err == movement.ErrorWrongUser || err == movement.ErrorInsufficientBalance {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}

			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, movementID)
	}
}

func searchMovement(movementService MovementService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, _ := strconv.ParseUint(ctx.DefaultQuery("limit", "0"), 10, 64)
		fmt.Print("limit", limit)
		offset, _ := strconv.ParseUint(ctx.DefaultQuery("offset", "0"), 10, 0)
		fmt.Print("offset", offset)
		userID, err := strconv.ParseInt(ctx.Query("userid"), 10, 64)
		fmt.Print("useridInt", userID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		typeMov := ctx.Query("type")
		fmt.Print("type", typeMov)
		currency := ctx.Query("currencyname")
		fmt.Print("currency", currency)
		ctx.Query("currencyname")
		var searchRequest = struct {
			Limit        uint64
			Offset       uint64
			Type         string
			CurrencyName string
			UserID       int64
		}{
			Limit:        limit,
			Offset:       offset,
			Type:         typeMov,
			CurrencyName: currency,
			UserID:       userID,
		}

		movementsResult, err := movementService.Search(ctx, searchRequest.Limit, searchRequest.Offset, searchRequest.Type,
			searchRequest.CurrencyName, searchRequest.UserID)
		if err != nil {
			if err == movement.ErrorNoMovements {
				ctx.JSON(http.StatusNotFound, err.Error())
				return
			}

			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, movementsResult)
	}
}
