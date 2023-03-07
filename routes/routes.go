package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kanatsanan6/go-test/controllers"
)

func Router() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	/*
		endpoint: POST /api/v1/sign_up
		description: use register to the website
		headers: not required
		params:
		 	first_name string
			last_name  string
			email      string
			password   string
		response:
		 	first_name string
			last_name  string
			email      string
			created_at string

	*/
	v1.POST("/sign_up", controllers.SignUp)

	/*
		endpoint: POST /api/v1/sign_in
		description: use to signin to the website
		headers: not required
		params:
			email      string
			password   string
		response:
			token      string
			expires_at int64
	*/
	v1.POST("/sign_in", controllers.SignIn)

	authRoutes := v1.Use(controllers.AuthMiddleware())

	/*
		endpoint: GET /api/v1/me
		description: returns current user information
		headers: authentication required
		params: -
		response:
			email:      string
			first_name: string
			last_name:  string
	*/
	authRoutes.GET("/me", controllers.Me)

	return r
}
