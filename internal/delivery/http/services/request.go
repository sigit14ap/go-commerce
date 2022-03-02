package services

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetIdFromPath(c *gin.Context, paramName string) (primitive.ObjectID, error) {
	idString := c.Param(paramName)
	if idString == "" {
		return primitive.ObjectID{}, errors.New("empty id param")
	}

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}

func GetIdFromRequestContext(context *gin.Context, paramName string) (primitive.ObjectID, error) {
	idString, ok := context.Get(paramName)
	if !ok {
		return primitive.ObjectID{}, errors.New("not authenticated")
	}

	id, err := primitive.ObjectIDFromHex(idString.(string))
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}

func GetDataFromContext(context *gin.Context, paramName string) (interface{}, error) {
	data, ok := context.Get(paramName)

	if !ok {
		return data, errors.New("Data from context not found")
	}

	return data, nil
}

func GetIdFromRequest(paramName string) (primitive.ObjectID, error) {
	if paramName == "" {
		return primitive.ObjectID{}, errors.New("empty id param")
	}

	id, err := primitive.ObjectIDFromHex(paramName)
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}
