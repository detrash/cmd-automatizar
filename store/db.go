package store

import (
	"context"
	"fmt"

	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreatePGXPool(user string, pass string, host string, port string, database string, app string) *pgxpool.Pool {

	connString := "postgres://" + user + ":" + pass + "@" + host + "/" + database + "?sslmode=disable" + "&" + "application_name=" + app
	configPool, errConfPool := pgxpool.ParseConfig(connString)
	if errConfPool != nil {
		panic(errConfPool)
	}
	configPool.MinConns = 1
	configPool.MaxConns = 4

	poolConn, errPool := pgxpool.ConnectConfig(context.Background(), configPool)
	if errPool != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", errPool)
		os.Exit(1)
	}
	return poolConn
}
