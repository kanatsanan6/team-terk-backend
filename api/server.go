package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kanatsanan6/go-test/db"
	dbConn "github.com/kanatsanan6/go-test/db/sqlc"
	"github.com/kanatsanan6/go-test/utils"
)

type Server struct {
	queries *dbConn.Queries
	router  *gin.Engine
}

func NewServer() (*Server, error) {
	server := &Server{
		queries: dbConn.New((db.DB)),
	}

	server.setupRouter()
	server.setupCors()

	return server, nil
}

func (server *Server) setupRouter() {
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
	v1.POST("/sign_up", server.SignUp)

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
	v1.POST("/sign_in", server.SignIn)

	authRoutes := v1.Use(AuthMiddleware())

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
	authRoutes.GET("/me", server.Me)

	server.router = r
}

func (server *Server) setupCors() {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	server.router.Use(cors.New(config))
}

func (server *Server) Start() {
	server.router.Run()
}

func dataResponse(response interface{}) gin.H {
	return gin.H{"data": response}
}

func errorResponse(err error) gin.H {
	return gin.H{"errors": err.Error()}
}

func validationErrorsResponse(err error) gin.H {
	var errors []map[string]string
	for _, err := range err.(validator.ValidationErrors) {
		field := utils.ToSnakeCase(err.Field())
		errorField := map[string]string{
			"field":   field,
			"message": fmt.Sprintf("%s is invalid with %s tag", field, err.Tag()),
		}
		errors = append(errors, errorField)
	}

	return gin.H{"errors": errors}
}
