package apiserver

import (
	"eastwh/internal/store"
	"errors"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not autenticated")
)

var (
	hmacSampleSecret = "8a046a6b436496d9c7af3e196a73ee9948677eb30b251706667ad59d6261bd78d2f6f501a6dea0118cfb3b0dcd62d6c9eb88142e2c24c2c686133a935cd65651"
)

type server struct {
	router *gin.Engine
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: gin.Default(),
		store:  store,
	}

	s.router.Use()

	s.router.Use(func(ctx *gin.Context) {
		fmt.Println("Запрос с сайта: ", ctx.Request.Header.Get("Origin"))
		ctx.Next()
	})
	confCors := cors.DefaultConfig()
	confCors.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	confCors.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	confCors.ExposeHeaders = []string{"Content-Length"}
	confCors.AllowCredentials = true
	confCors.MaxAge = 12 * time.Hour
	confCors.AllowOriginFunc = func(origin string) bool {
		return true
	}

	s.router.Use(cors.New(confCors))

	s.configureRouter()

	return s
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *server) configureRouter() {
	apiGroup := s.router.Group("/api/v1")
	{
		userGroup := apiGroup.Group("/user", s.AuthMW)
		{
			userGroup.POST("/logout/", s.Logout)
			userGroup.POST("/update", s.UpdateUser)
			userGroup.POST("/update/pass", s.UpdatePassword)
			userGroup.POST("/:id/block/", s.BlockedUser)
		}

		usersGroup := apiGroup.Group("/users")
		{
			usersGroup.POST("", s.AddUser)
			usersGroup.POST("/login", s.Login)
			usersGroup.GET("", s.GetUsers)
			usersGroup.POST("/password/restore", s.SetUserTemporaryPassword)
		}
	}
}

func (s *server) AddUser(ctx *gin.Context) {

}

func (s *server) UpdateUser(ctx *gin.Context) {

}

func (s *server) Login(ctx *gin.Context) {

}

func (s *server) AddLoginUser(ctx *gin.Context) {

}

func (s *server) SetUserTemporaryPassword(ctx *gin.Context) {

}

func (s *server) GetUsers(ctx *gin.Context) {

}

func (s *server) Logout(ctx *gin.Context) {

}

func (s *server) UpdatePassword(ctx *gin.Context) {

}

func (s *server) BlockedUser(ctx *gin.Context) {

}
