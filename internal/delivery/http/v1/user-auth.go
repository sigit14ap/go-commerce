package v1

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	"github.com/sigit14ap/go-commerce/pkg/auth"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate *validator.Validate = validator.New()

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
		errorResponse(context, http.StatusBadRequest, "invalid input body")
		return
	}

	user, err := h.services.Users.FindByCredentials(context, signInDTO)

	log.Error(h.services.Users.CheckPasswordHash(signInDTO.Password, user.Password))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			errorResponse(context, http.StatusUnauthorized, "Email not found")
		} else {
			errorResponse(context, http.StatusInternalServerError, err.Error())
		}
		return
	}

	matchPassword := h.services.Users.CheckPasswordHash(signInDTO.Password, user.Password)
	if matchPassword == false {
		errorResponse(context, http.StatusUnauthorized, "Password does not match")
		return
	}

	userClaims := jwt.MapClaims{"userID": user.ID}
	authDetails, err := h.tokenProvider.CreateJWTSession(auth.CreateSessionInput{
		Fingerprint: signInDTO.Fingerprint,
		Claims:      userClaims,
	})

	if err != nil {
		errorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}
	successResponse(context, authDetails)
}

// UserSignUp godoc
// @Summary  User sign-up
// @Tags     user-auth
// @Accept   json
// @Produce  json
// @Param    user  body      dto.SignUpDTO  true  "user data"
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
			log.Error(fieldErr)
			errorResponse(context, http.StatusBadRequest, fmt.Sprintf(fieldErr.Error()))
			return // exit on first error
		}
	}

	var cartID primitive.ObjectID
	cartIDHex, err := context.Cookie("cartID")
	if err == nil {
		cartID, err = primitive.ObjectIDFromHex(cartIDHex)
		if err != nil {
			log.Warnf("failed to convert cart id %s to object id", cartIDHex)
		}
	}

	user, err := h.services.Users.Create(context, dto.CreateUserDTO{
		Name:     signUpDTO.Name,
		Email:    signUpDTO.Email,
		Password: signUpDTO.Password,
		CartID:   cartID,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			errorResponse(context, http.StatusInternalServerError,
				fmt.Sprintf("user with email %s already exists", signUpDTO.Email))
		} else {
			errorResponse(context, http.StatusInternalServerError, err.Error())
		}
		return
	}

	createdResponse(context, domain.UserInfo{
		Name:  user.Name,
		Email: user.Email,
	})
}

// UserRefresh godoc
// @Summary  User refresh token
// @Tags     user-auth
// @Accept   json
// @Produce  json
// @Param    refreshInput  body      auth.RefreshInput  true  "user refresh data"
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
