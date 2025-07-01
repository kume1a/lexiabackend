package shared

import (
	"lexia/ent"
	"log"

	_ "github.com/lib/pq"
)

func IsDatabaseErorNotFound(err error) bool {
	if err == nil {
		return false
	}

	return ent.IsNotFound(err)
}

func InitializeDatabase() (*ent.Client, error) {
	envVars, err := ParseEnv()
	if err != nil {
		log.Fatalf("failed to parse environment variables: %v", err)
		return nil, err
	}

	client, err := ent.Open("postgres", envVars.DbConnectionString)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
		return nil, err
	}

	return client, nil
}
