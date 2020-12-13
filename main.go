package main

// "users"
// "net/http"
import (
	"covid/covid"
	"covid/users"
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()
	// store := cookie.NewStore([]byte("covid19"))
	r.Use(corsMiddleware())
	// r.Use(sessions.Sessions("session_covid", store))
	r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "Corona Data",
		Key:             []byte("key_covid19"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     users.IdentityKey,
		PayloadFunc:     users.Payload,
		IdentityHandler: users.Identity,
		Authenticator:   users.Login,
		Authorizator:    users.Authorizator,
		Unauthorized:    users.Unauthorized,
		TokenLookup:     "header : Authorization, query : token, cookie : jwt",
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {

			
			// session := sessions.Default(c)
			// cookies := c

			// if err != nil {
			// 	fmt.Printf("error cookie : %+v", err)
			// }
			// type map_value struct {
			// 	kd_rs     int
			// 	nama_user string
			// 	nm_rs     string
			// 	userid    int
			// 	written   bool
			// 	writer    http.ResponseWriter
			// }
			// data := reflect.ValueOf(c.Keys["github.com/gin-contrib/sessions"]).Elem()
			// {name:session_covid request:0xc000244100 store:0xc000006708 session:0xc0002ba0f0 written:false writer:0xc0002b8000
			// check_tipe := reflect.TypeOf(data)
			// fmt.Println(data)
			// export_name := data.FieldByName("name")
			// export_name = reflect.NewAt(export_name.Type(), unsafe.Pointer(export_name.UnsafeAddr())).Elem()
			// export_request := data.FieldByName("request")
			// export_request = reflect.NewAt(export_request.Type(), unsafe.Pointer(export_request.UnsafeAddr())).Elem()
			// export_store := data.FieldByName("store")
			// export_store = reflect.NewAt(export_store.Type(), unsafe.Pointer(export_store.UnsafeAddr())).Elem()
			// export_session := data.FieldByName("session")
			// export_session = reflect.NewAt(export_session.Type(), unsafe.Pointer(export_session.UnsafeAddr())).Elem()

			// fmt.Printf("tes name %+v", export_name)
			// fmt.Printf("tes request %+v\n", export_request)
			// sesi := reflect.ValueOf(export_session).Elem()
			// sesi = reflect.ValueOf(sesi).Elem()
			// fmt.Printf("tes sesi %+v\n", sesi)
			// mv := reflect.TypeOf(map_value).Elem()
			// values := reflect.Indirect(export_session).FieldByName("Values")
			// values = reflect.NewAt(values.Type(), unsafe.Pointer(values.UnsafeAddr())).Elem()
			// check := reflect.ValueOf(values).Interface().(reflect.Value)

			// var kd_rs int
			// if check.Kind() == reflect.Map {
			// 	for _, e := range check.MapKeys() {
			// 		v := check.MapIndex(e)
			// 		switch t := v.Interface().(type) {
			// 		case int:
			// 			if e.Interface().(string) == "kd_rs" {
			// 				// fmt.Println("int ", e, t)
			// 				kd_rs = t
			// 			}
			// 		case string:
			// 			fmt.Println("string ", e, t)
			// 		case bool:
			// 			fmt.Println("bool ", e, t)
			// 		default:
			// 			// fmt.Println("not found")

			// 		}
			// 	}
			// }

			// fmt.Printf("values %+v", kd_rs)

			c.JSON(code, gin.H{
				"code":   code,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
				// "kd_rs":  kd_rs,
			})
		},
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		fmt.Printf("JWT error" + err.Error())
	}

	// // fmt.Println(gin.Context.FullPath())
	r.POST("/covid/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		fmt.Printf("Noroute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/covid")
	auth.GET("/security_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", hello)
		auth.POST("/data_board", covid.Get_data_board)
		auth.GET("/get_user", hello)
	}

	r.Run(":8088")

}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// AllowOrigins:  []string{"http://192.168.3.62:3000, http://rsd-covid.bee-health.com:92"},
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"POST, OPTIONS"},
		AllowHeaders:  []string{"Origin, Authorization, Content-Type, Cookie, Set-Cookie"},
		ExposeHeaders: []string{"Content-Length, Set-Cookie, Cookie"},
	})
	// return func(c *gin.Context) {
	// 	// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.3.62:3000")
	// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Set-Cookie, Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// 	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(200)
	// 	} else {
	// 		c.Next()
	// 	}
	// }
}

func hello(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(users.IdentityKey)

	c.JSON(200, gin.H{
		"code":     200,
		"message":  "bisa login",
		"userid":   claims[users.IdentityKey],
		"username": user.(*users.User).Username,
	})
}
