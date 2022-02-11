package tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/restuwahyu13/go-supertest/supertest"
	v1 "github.com/sigit14ap/go-commerce/internal/delivery/http/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Method     string      `json:"method"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

var router = v1.NewHandler()

func TestPostMethod(t *testing.T) {
	test := supertest.NewSuperTest(router, t)

	payload := gin.H{
		"name": "tes",
	}

	test.Post("/")
	test.Send(payload)
	test.Set("Content-Type", "application/json")
	test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

		var response Response
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, req.Method, req.Method)
		assert.Equal(t, "fetch request using post method", response.Message)
	})
}
