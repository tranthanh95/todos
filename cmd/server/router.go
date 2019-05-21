package server

import (
	"todo-lists/pkg/http/rest"
	"todo-lists/pkg/login"
	"todo-lists/pkg/user"
	"todo-lists/pkg/todo"
	"todo-lists/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func (ds *dserver) MapRoutes() {

	// Group : v1
	apiV1 := ds.router.Group("api/v1")

	ds.healthRoutes(apiV1)
	ds.loginRoutes(apiV1)
	ds.userRoutes(apiV1)
	ds.todoRoutes(apiV1)
}

func (ds *dserver) healthRoutes(api *gin.RouterGroup) {
	healthRoutes := api.Group("/health")
	{
		h := rest.NewHealthCtrl()
		healthRoutes.GET("/", h.Ping)
	}
}

func (ds *dserver) loginRoutes(api *gin.RouterGroup) {
	var loginSvc login.Service
	ds.cont.Invoke(func(l login.Service) {
		loginSvc = l
	})

	loginRoutes := api.Group("/login")
	{
		f := rest.NewLoginCtrl(ds.logger, loginSvc)
		loginRoutes.POST("/", f.Signin)
	}
}

func (ds *dserver) userRoutes(api *gin.RouterGroup) {
	userRoutes := api.Group("/users")
	{
		var userSvc user.Service
		ds.cont.Invoke(func(u user.Service) {
			userSvc = u
		})

		usr := rest.NewUserCtrl(ds.logger, userSvc)

		userRoutes.GET("/", usr.GetAll)
		userRoutes.POST("/", usr.Store)
		userRoutes.GET("/:id", usr.GetByID)
		userRoutes.PUT("/:id", usr.Update)
		userRoutes.DELETE("/:id", usr.Delete)
	}
}

func (ds *dserver) todoRoutes(api *gin.RouterGroup) {
	var appMiddleware middleware.ApiMiddleware
	todoRoutes := api.Group("/todos").Use(appMiddleware.VerifyToken())
	{
		var todoSvc todo.Service
		ds.cont.Invoke(func(t todo.Service) {
			todoSvc = t
		})

		td := rest.NewTodoCtrl(ds.logger, todoSvc)

		todoRoutes.GET("/", td.GetAll)
		// todoRoutes.POST("/", td.Store)
		// todoRoutes.GET("/:id", td.GetByID)
		// todoRoutes.PUT("/:id", td.Update)
		// todoRoutes.DELETE("/:id", td.Delete)
	}
}