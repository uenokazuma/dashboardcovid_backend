package config

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	"github.com/jackc/pgx/v4"
// 	"gopkg.in/ini.v1"
// 	// "github.com/gin-gonic/gin"
// )

// func allow_origin() (allow string) {

// 	config, err_cf := ini.Load("config.ini")

// 	if err_cf != nil {
// 		fmt.Printf("Gagal membaca file config : %v", err_cf)
// 		os.Exit(1)
// 	}

// 	host := config.Section(section).Key("allow-origin").String()

// 	config_db = "host=" + host + " port=" + port + " user=" + user + " password=" + password + " database=" + database

// 	return allow
// }

// func Connect(section string) *pgx.Conn {

// 	connect, err_db := pgx.Connect(context.Background(), config_db(section))

// 	if err_db != nil {
// 		fmt.Printf("Gagal koneksi : %v", err_db)
// 		os.Exit(1)
// 	}

// 	return connect
// }
