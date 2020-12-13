package users

// "users"
// "net/http"
import (
	"context"
	"covid/config"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type User struct {
	Id          int              `json:"Id"`
	Created_at  pgtype.Timestamp `json:"Created_at"`
	Email       pgtype.Varchar   `json:"Email"`
	Password    pgtype.Varchar   `json:"Password"`
	Updated_at  pgtype.Timestamp `json:"Updated_at"`
	Username    string           `json:"Username"`
	Authy_id    pgtype.Varchar   `json:"Authy_id"`
	Phone       pgtype.Varchar   `json:"Phone"`
	Verified_at pgtype.Timestamp `json:"Verified_at"`
	Otp         pgtype.Varchar   `json:"Otp"`
	Employee_id pgtype.Varchar   `json:"Employee_id"`
}

type User_hospital struct {
	hospital_id   int
	nama_hospital string
}

var connect *pgx.Conn

func init() {
	var section string
	section = "database_user"

	connect = config.Connect(section)
}

func Get_user(c *gin.Context) {
	// var nama_user pgtype.Varchar
	var data User

	check := data_user(data.Username)

	fmt.Printf("Check : %+v\n", check[0])

	c.JSON(200, gin.H{
		"message": check,
	})

}

func data_user(username string) []User {

	// connect := connect()

	// rows, _ := connect.Query(context.Background(), "select id, created_at, email, password, updated_at, username, authy_id, phone, verified_at, otp, employee_id from users")
	rows, _ := connect.Query(
		context.Background(),
		`select id,
				created_at,
				email,
				password,
				updated_at,
				username,
				authy_id,
				phone,
				verified_at,
				otp,
				employee_id
		   from users
		  where username = coalesce($1, username)`, username)
	// rows, _ := connect.Query(context.Background(), "select * from users where username = coalesce($1, username)", username)
	// err := connect.QueryRow(context.Background(), "select id, created_at, email, password, updated_at, username, authy_id, phone, verified_at, otp, employee_id from users").Scan(&user.id, &user.created_at, &user.email, &user.password, &user.updated_at, &user.username, &user.authy_id, &user.phone, &user.verified_at, &user.otp, &user.employee_id)
	defer rows.Close()

	data := []User{}
	for rows.Next() {
		var user_row User
		// err := rows.Scan(&user_row)
		// err := rows.Scan(&user_row.id, &user_row.created_at, &user_row.email, &user_row.password, &user_row.updated_at, &user_row.username, &user_row.authy_id, &user_row.phone, &user_row.verified_at, &user_row.otp, &user_row.employee_id)
		err := rows.Scan(&user_row.Id, &user_row.Created_at, &user_row.Email, &user_row.Password, &user_row.Updated_at, &user_row.Username, &user_row.Authy_id, &user_row.Phone, &user_row.Verified_at, &user_row.Otp, &user_row.Employee_id)

		if err != nil {
			fmt.Printf("Query failed: %v\n", err)
		}

		data = append(data, user_row)
	}

	// fmt.Printf("%+v\n", data)

	return data
}

func data_user_hospital(user_id int) User_hospital {

	rows := connect.QueryRow(context.Background(),
		`select users_hospitals.hospital_id,
				hospitals.name
		   from users_hospitals, hospitals
		  where users_hospitals.hospital_id = hospitals.id
		    and users_hospitals.user_id = $1`, user_id)

	var user_row User_hospital
	err := rows.Scan(&user_row.hospital_id, &user_row.nama_hospital)

	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
	}

	return user_row
}
