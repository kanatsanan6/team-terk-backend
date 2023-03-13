package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	dbConn "github.com/kanatsanan6/go-test/db/sqlc"
	"github.com/kanatsanan6/go-test/utils"
)

type signUpRequest struct {
	FirstName            string `json:"first_name" binding:"required,min=1"`
	LastName             string `json:"last_name" binding:"required,min=1"`
	CompanyName          string `json:"company_name" binding:"required"`
	Email                string `json:"email" binding:"required,min=1"`
	Password             string `json:"password" binding:"required,min=1"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

func (server *Server) SignUp(ctx *gin.Context) {
	var req signUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ValidationErrorsResponse(err))
		return
	}

	if req.Password != req.PasswordConfirmation {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Password does not match with confirmation"})
		return
	}

	result, err := server.store.SignUpTx(ctx, dbConn.SignUpTxParams{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		CompanyName: req.CompanyName,
		Email:       req.Email,
		Password:    req.Password,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": "duplicated email"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, result)
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
		ctx.JSON(http.StatusBadRequest, utils.ValidationErrorsResponse(err))
		return
	}

	user, err := server.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	err = utils.ValidatePassword([]byte(req.Password), user.EncryptedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	jwtToken, jwtPayload, err := utils.GenerateJwtToken(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	response := signInResponse(jwtToken, jwtPayload.ExpiresAt)

	ctx.JSON(http.StatusOK, utils.DataResponse(response))
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
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	response := MeResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	ctx.JSON(http.StatusOK, utils.DataResponse(response))
}
