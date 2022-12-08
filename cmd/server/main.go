package main

import (
	"database/sql"
	"fmt"
	"log"

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

	conn, err := sql.Open(c.DBDriver, c.DBUrl)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := services.NewServer(store, &c)

	socket := services.InitSocketServer()

	routes.InitGoth(&c)
	route := routes.InitRoutes(server, socket, &c)
	go func() {
		if err := socket.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer socket.Close()

	address := fmt.Sprintf(":%v", c.Port)
	err = route.Run(address)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
