package main

import (
	"red-feed/internal/repository"
	"red-feed/internal/repository/dao"
	"red-feed/internal/service"
	"red-feed/internal/web"
	"red-feed/internal/web/middleware"
	"red-feed/pkg/ginx/middlewares/ratelimit"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	server := initWebServer() // 初始化web server

	u := initUser(db)
	u.RegisterRoutes(server)

	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 1).Build())
	server.Use(cors.New(cors.Config{
		// AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
