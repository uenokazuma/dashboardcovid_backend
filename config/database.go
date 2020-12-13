package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"gopkg.in/ini.v1"
	// "github.com/gin-gonic/gin"
)

func config_db(section string) (config_db string) {

	config, err_cf := ini.Load("config.ini")

	if err_cf != nil {
		fmt.Printf("Gagal membaca file config : %v", err_cf)
		os.Exit(1)
	}

	host := config.Section(section).Key("host").String()
	port := config.Section(section).Key("port").String()
	user := config.Section(section).Key("user").String()
	password := config.Section(section).Key("password").String()
	database := config.Section(section).Key("database").String()

	config_db = "host=" + host + " port=" + port + " user=" + user + " password=" + password + " database=" + database

	return config_db
}

func Connect(section string) *pgx.Conn {

	connect, err_db := pgx.Connect(context.Background(), config_db(section))

	if err_db != nil {
		fmt.Printf("Gagal koneksi : %v", err_db)
		os.Exit(1)
	}

	return connect
}
