package users

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	// "github.com/gin-gonic/gin"
)

type user_login struct {
	User     string `json:"user"   binding:"required"`
	Password string `json:"password" binding:"required"`
}

var IdentityKey = "id"

func Identity(c *gin.Context) interface{} {

	claims := jwt.ExtractClaims(c)

	return &User{
		Username: claims[IdentityKey].(string),
	}
}

func Login(c *gin.Context) (interface{}, error) {
	var data user_login

	if err := c.ShouldBindJSON(&data); err != nil {
		// c.JSON(400, gin.H{"error": err.Error()})
		return "", jwt.ErrMissingLoginValues
	}

	check := data_user(data.User)
	password := []byte(data.Password)
	err_compare := bcrypt.CompareHashAndPassword([]byte(check[0].Password.String), password)
	// fmt.Printf("user = %v\n", data.User)
	// fmt.Printf("check = %v\n", check)
	// fmt.Printf("password = %v\n", password)
	// fmt.Printf("err_compare = %v\n", err_compare)

	if err_compare == nil {
		// session := sessions.Default(c)
		// userid := check[0].Id

		// user_hospital := data_user_hospital(userid)
		// fmt.Printf("session : tes")
		// session.Set("userid", userid)
		// session.Set("nama_user", data.User)
		// session.Set("kd_rs", user_hospital.hospital_id)
		// session.Set("nm_rs", user_hospital.nama_hospital)
		// session.Save()
		// c.JSON(200, gin.H{
		// 	"code":    200,
		// 	"message": "you are logged in",
		// })
		return check, nil
	} else {
		// c.JSON(401, gin.H{
		// 	"code":    401,
		// 	"message": "unauthorized",
		// })
		return nil, jwt.ErrFailedAuthentication
	}

}

func Authorizator(data interface{}, c *gin.Context) bool {
	// v, ok := data.(*User)
	// fmt.Printf("Authorizator data : %+v", data)
	// fmt.Printf("Authorizator v : %+v", v)
	// fmt.Printf("Authorizator ok : %v", ok)
	// fmt.Printf("data : %+v", data)
	// data := c.Keys["JWT_PAYLOAD"]
	get_data := c.Keys["JWT_PAYLOAD"].(jwt.MapClaims)
	// fmt.Printf("get_data : %+v", get_data)
	// fmt.Printf("get_data : %v", get_data["id"])

	// fmt.Printf("c : %+v", c.Keys["JWT_PAYLOAD"].string["id"])
	if v, ok := data.(*User); ok && v.Username == get_data["id"] {
		// fmt.Printf("Authorizator true ok : %v", ok)
		return true
	}
	// fmt.Printf("Authorizator ok : %v", ok)
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func Payload(data interface{}) jwt.MapClaims {

	if v, ok := data.([]User); ok {
		return jwt.MapClaims{
			IdentityKey: v[0].Username,
		}
	}

	return jwt.MapClaims{}
}
