package config

import (
	"context"
	"log"

	"entgo.io/ent"

	_ "github.com/lib/pq"
)

func InitializeDatabase() {
	envVars, err := ParseEnv()
	if err != nil {
		log.Fatalf("failed to parse environment variables: %v", err)
		return nil
	}

	client, err := ent.Open("postgres", envVars.DbConnectionString)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
