package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kanatsanan6/go-test/db"
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

func SignUp(ctx *gin.Context) {
	var req signUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
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

	queries := dbConn.New(db.DB)
	insertedUser, err := queries.CreateUser(ctx, args)
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

func SignIn(ctx *gin.Context) {
	var req SignInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
	}

	queries := dbConn.New(db.DB)
	user, err := queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"errors": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	err = utils.ValidatePassword([]byte(req.Password), user.EncryptedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
		return
	}

	jwtToken, jwtPayload, err := utils.GenerateJwtToken(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	response := signInResponse(jwtToken, jwtPayload.ExpiresAt)

	ctx.JSON(http.StatusOK, response)
}

type MeResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func Me(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*utils.CustomClaims)

	queries := dbConn.New(db.DB)
	user, err := queries.GetUserByEmail(ctx, authPayload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"errors": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	response := MeResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}
