package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DataSource *pgxpool.Pool

func Connect(postgresURL string) {
	var err error
	DataSource, err = pgxpool.Connect(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	log.Println("Connected to database")
}

func CloseDatabase() {
	DataSource.Close()
}
