package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

type User struct {
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
}

func TestGet(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/user", req.URL.Path)
		assert.Equal(t, http.MethodGet, req.Method)
		user := User{"Alejandro", "Galue"}
		bytes, _ := json.Marshal(user)
		res.WriteHeader(http.StatusOK)
		res.Write(bytes)
	}))
	defer testServer.Close()

	Instance.URL = testServer.URL
	bytes, err := Instance.Get("/user")

	assert.NilError(t, err)
	var user User
	json.Unmarshal(bytes, &user)
	assert.Equal(t, "Alejandro", user.FirstName)
	assert.Equal(t, "Galue", user.LastName)
}

func TestPost(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/user", req.URL.Path)
		assert.Equal(t, http.MethodPost, req.Method)
		var user User
		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusForbidden)
			return
		}
		json.Unmarshal(bytes, &user)
		assert.Equal(t, "Alejandro", user.FirstName)
		assert.Equal(t, "Galue", user.LastName)
		res.WriteHeader(http.StatusNoContent)
	}))
	defer testServer.Close()

	Instance.URL = testServer.URL
	user := User{"Alejandro", "Galue"}
	bytes, _ := json.Marshal(user)

	err := Instance.Post("/user", bytes)
	assert.NilError(t, err)
}

func TestDelete(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/user", req.URL.Path)
		assert.Equal(t, http.MethodDelete, req.Method)
		res.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	Instance.URL = testServer.URL
	err := Instance.Delete("/user")

	assert.NilError(t, err)
}
