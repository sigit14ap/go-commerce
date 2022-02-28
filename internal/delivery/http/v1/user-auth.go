package v1

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/pkg/auth"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) initUserAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-in", h.userSignIn)
		auth.POST("/sign-up", h.userSignUp)
		auth.POST("/refresh", h.userRefresh)
	}
}

// UserSignIn godoc
// @Summary  User sign-in
// @Tags     user-auth
// @Accept   json
// @Produce  json
// @Param    user  body      dto.SignInDTO  true  "user credentials"
// @Success  200   {object}  auth.AuthDetails
// @Failure  400   {object}  failure
// @Failure  401   {object}  failure
// @Failure  404   {object}  failure
// @Failure  500   {object}  failure
// @Router   /users/auth/sign-in [post]
func (h *Handler) userSignIn(context *gin.Context) {
	var signInDTO dto.SignInDTO

	err := context.BindJSON(&signInDTO)
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	user, err := h.services.Users.FindByCredentials(context, signInDTO)

	log.Error(h.services.Users.CheckPasswordHash(signInDTO.Password, user.Password))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ErrorResponse(context, http.StatusUnauthorized, "Email not found")
		} else {
			ErrorResponse(context, http.StatusInternalServerError, err.Error())
		}
		return
	}

	matchPassword := h.services.Users.CheckPasswordHash(signInDTO.Password, user.Password)
	if matchPassword == false {
		ErrorResponse(context, http.StatusUnauthorized, "Password does not match")
		return
	}

	userClaims := jwt.MapClaims{"userID": user.ID}
	authDetails, err := h.tokenProvider.CreateJWTSession(auth.CreateSessionInput{
		Fingerprint: signInDTO.Fingerprint,
		Claims:      userClaims,
	})

	if err != nil {
		ErrorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}
	successResponse(context, authDetails)
}

// UserSignUp godoc
// @Summary  User sign-up
// @Tags     user-auth
// @Accept   json
// @Produce  json
// @Param    user  body      dto.SignUpDTO  true  "user services"
// @Success  200   {object}  domain.UserInfo
// @Failure  400   {object}  failure
// @Failure  401   {object}  failure
// @Failure  404   {object}  failure
// @Failure  500   {object}  failure
// @Router   /users/auth/sign-up [post]
func (h *Handler) userSignUp(context *gin.Context) {
	var signUpDTO dto.SignUpDTO

	err := context.ShouldBindJSON(&signUpDTO)
	if err != nil {

		for _, fieldErr := range err.(validator.ValidationErrors) {
			ErrorResponse(context, http.StatusUnprocessableEntity, fmt.Sprintf(fieldErr.Error()))
			return // exit on first error
		}
	}

	user, err := h.services.Users.Create(context, dto.CreateUserDTO{
		Name:     signUpDTO.Name,
		Email:    signUpDTO.Email,
		Password: signUpDTO.Password,
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ErrorResponse(context, http.StatusInternalServerError,
				fmt.Sprintf("user with email %s already exists", signUpDTO.Email))
		} else {
			ErrorResponse(context, http.StatusInternalServerError, err.Error())
		}
		return
	}

	createdResponse(context, domain.UserInfo{
		Name:  user.Name,
		Email: user.Email,
	})
	return
}

// UserRefresh godoc
// @Summary  User refresh token
// @Tags     user-auth
// @Accept   json
// @Produce  json
// @Param    refreshInput  body      auth.RefreshInput  true  "user refresh services"
// @Success  200           {object}  auth.AuthDetails
// @Failure  400           {object}  failure
// @Failure  401           {object}  failure
// @Failure  404           {object}  failure
// @Failure  500           {object}  failure
// @Router   /users/auth/refresh [post]
func (h *Handler) userRefresh(context *gin.Context) {
	h.refreshToken(context)
}

func (h *Handler) verifyUser(context *gin.Context) {
	h.verifyToken(context, "userID")
}
