package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	db "github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/routes"
	"github.com/vtv-us/kahoot-backend/internal/services"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

func main() {
	c, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatal("can't load config", err)
	}

	fmt.Println(c)

	conn, err := sql.Open(c.DB_DRIVER, c.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := services.NewServer(store, &c)
	route := routes.InitRoutes(server)

	err = route.Run(c.SERVER_ADDRESS)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
