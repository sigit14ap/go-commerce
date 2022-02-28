package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/go-commerce/pkg/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func extractAuthToken(context *gin.Context) (string, error) {
	authHeader := context.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Auth header is empty")
	}

	headerParts := strings.Split(authHeader, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}

func (h *Handler) refreshToken(context *gin.Context) {
	var input auth.RefreshInput
	err := context.BindJSON(&input)
	if err != nil {
		ErrorResponse(context, http.StatusBadRequest, "Invalid request body")
		return
	}

	authDetails, err := h.tokenProvider.Refresh(auth.RefreshInput{
		RefreshToken: input.RefreshToken,
		Fingerprint:  input.Fingerprint,
	})

	if err != nil {
		ErrorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	successResponse(context, authDetails)
}

func (h *Handler) verifyToken(context *gin.Context, idName string) {
	tokenString, err := extractAuthToken(context)
	if err != nil {
		ErrorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	tokenClaims, err := h.tokenProvider.VerifyToken(tokenString)
	if err != nil {
		ErrorResponse(context, http.StatusUnauthorized, err.Error())
		return
	}

	id, ok := tokenClaims[idName]
	if !ok {
		ErrorResponse(context, http.StatusForbidden, "this endpoint is forbidden")
		return
	}

	context.Set(idName, id)
}

func (h *Handler) extractIdFromAuthHeader(context *gin.Context, idName string) (primitive.ObjectID, error) {
	tokenString, err := extractAuthToken(context)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	tokenClaims, err := h.tokenProvider.VerifyToken(tokenString)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	idHex, ok := tokenClaims[idName]
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("failed to extract %s from auth header", idName)
	}

	id, err := primitive.ObjectIDFromHex(idHex.(string))
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to convert %s to objectId", idHex)
	}

	return id, nil
}
