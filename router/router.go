package router

import (
	"github.com/julienschmidt/httprouter"
	"go.uber.org/dig"
	"rezvin-pro-bot/config"
	"rezvin-pro-bot/internal/logger"
)

type IRouter interface {
	GetHttpRouter() *httprouter.Router
}

type routerDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`
	Config config.IConfig `name:"Config"`
}

type router struct {
	httpRouter *httprouter.Router

	logger logger.ILogger
	config config.IConfig
}

func NewRouter(deps routerDependencies) *router {
	r := &router{
		httpRouter: httprouter.New(),
		logger:     deps.Logger,
		config:     deps.Config,
	}

	r.registerRoutes()

	return r
}

func (router *router) registerRoutes() {
}

func (router *router) GetHttpRouter() *httprouter.Router {
	return router.httpRouter
}
