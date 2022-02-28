package tests

import (
	"github.com/gavv/httpexpect/v2"
	v1 "github.com/sigit14ap/go-commerce/internal/delivery/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	Data interface{} `json:"services"`
}

var h = v1.Handler{}
var router = h.Init()

func TestSignUp(t *testing.T) {

	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	payload := map[string]interface{}{
		"name":     "foo",
		"email":    "foo@gmail.com",
		"password": "bar",
	}

	e.POST("/api/v1/users/auth/sign-up").
		WithJSON(payload).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().ContainsKey("weight").ValueEqual("weight", 100)

	//test := supertest.NewSuperTest(router, t)
	//
	//payload := map[string]interface{}{
	//	"name":     "foo",
	//	"email":    "foo@gmail.com",
	//	"password": "bar",
	//}
	//
	//test.Post("/api/v1/users/auth/sign-up")
	//test.Send(payload)
	//test.Set("Content-Type", "application/json")
	//test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {
	//	//var services Response
	//	//json.Unmarshal(rr.Body.Bytes(), &services)
	//
	//	assert.Equal(t, http.StatusCreated, rr.Code)
	//})
}
