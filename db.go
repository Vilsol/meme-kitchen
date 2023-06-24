package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"memekitchen/ent"

	_ "github.com/lib/pq"
)

func ConnectToDB() *ent.Client {
	client, err := ent.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			viper.GetString("database.postgres.host"),
			viper.GetInt("database.postgres.port"),
			viper.GetString("database.postgres.user"),
			viper.GetString("database.postgres.db"),
			viper.GetString("database.postgres.pass"),
		),
	)

	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	if err := client.Schema.Create(context.TODO()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}
