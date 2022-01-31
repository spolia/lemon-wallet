package api

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spolia/lemon-wallet/cmd/api/internal"
	"github.com/spolia/lemon-wallet/internal/wallet"
	"github.com/spolia/lemon-wallet/internal/wallet/movement"
	"github.com/spolia/lemon-wallet/internal/wallet/user"
)

func StarApplication() () {
	log.Println("starting")
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", "root", "root", "127.0.0.1:3306", "wallet","parseTime=true")
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userService := wallet.NewUserService(user.New(db), movement.New(db))
	movementService := wallet.NewMovementService(user.New(db), movement.New(db))
	log.Println("database successfully configured")
	router := gin.Default()
	internal.API(router, userService, movementService)

	router.Run("localhost:8080")
	log.Println("listening")
}
