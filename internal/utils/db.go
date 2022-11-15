package utils

import (
	"fmt"
	"log"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

func MigrateDB(config Config) error {
	opts, err := pg.ParseURL(config.DBUrl)
	if err != nil {
		log.Fatal("can't parse db url", err)
	}

	conn := pg.Connect(opts)
	defer conn.Close()

	collection := migrations.NewCollection()
	err = collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return fmt.Errorf("cannot discover migrations: %w", err)
	}

	collection.Run(conn, "init")
	collection.Run(conn, "up")

	return err
}
