package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	dbConn "github.com/kanatsanan6/go-test/db/sqlc"
	"github.com/kanatsanan6/go-test/utils"
)

type signUpRequest struct {
	FirstName string `json:"first_name" binding:"required,min=1"`
	LastName  string `json:"last_name" binding:"required,min=1"`
	Email     string `json:"email" binding:"required,min=1"`
	Password  string `json:"password" binding:"required,min=1"`
}

type userResponse struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func signUpResponse(user dbConn.User) userResponse {
	return userResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) SignUp(ctx *gin.Context) {
	var req signUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, validationErrorsResponse(err))
		return
	}

	hashPassword, err := utils.GeneratePassword([]byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	args := dbConn.CreateUserParams{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		EncryptedPassword: hashPassword,
	}

	insertedUser, err := server.queries.CreateUser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	response := signUpResponse(insertedUser)

	ctx.JSON(http.StatusCreated, response)
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=1"`
}

type tokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expired_at"`
}

func signInResponse(token string, expiresAt int64) tokenResponse {
	return tokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

func (server *Server) SignIn(ctx *gin.Context) {
	var req SignInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, validationErrorsResponse(err))
		return
	}

	user, err := server.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.ValidatePassword([]byte(req.Password), user.EncryptedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	jwtToken, jwtPayload, err := utils.GenerateJwtToken(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := signInResponse(jwtToken, jwtPayload.ExpiresAt)

	ctx.JSON(http.StatusOK, dataResponse(response))
}

type MeResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (server *Server) Me(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*utils.CustomClaims)

	user, err := server.queries.GetUserByEmail(ctx, authPayload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := MeResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	ctx.JSON(http.StatusOK, dataResponse(response))
}
