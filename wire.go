//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"red-feed/internal/repository"
	"red-feed/internal/repository/cache"
	"red-feed/internal/repository/dao"
	"red-feed/internal/service"
	"red-feed/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitRedis, ioc.InitDB,

		// 初始化DAO层 和 Cache层
		dao.NewUserDAO,
		cache.NewUserCache,
		cache.NewCodeCache,

		// 初始化Repo层
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// 初始化Service层
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSMSService,

		// web handler
		web.NewUserHandler,

		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	return new(gin.Engine)
}