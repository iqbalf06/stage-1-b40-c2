package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() { //huruf besar agar bisa dipanggil diluar folder

	//user:password@host:port/database_name
	databaseUrl := "postgres://postgres:admin@localhost:5432/personal_web_b40"

	var err error
	Conn, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil { //kondisi error. jika terjadi error, maka akan tampil.
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success connect to database")
}
