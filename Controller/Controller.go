package Controller

import (
	"encoding/json"
	"fmt"
	Repository "github.com/vadim-shalnev/API_Server_dadata/Repository"
	Service "github.com/vadim-shalnev/API_Server_dadata/Service"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type RequestAddressInfo struct {
	Addres string `json:"addres"`
}
type RequestAddressSearch struct {
	Query string `json:"query"`
}
type TokenString struct {
	T string `json:"token"`
}
type NewUser struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

// @title Todo geocode API
// @version 1.0
// @description API Server for search GEOinfo

// @host localhost:8080
// @BasePath /api/address

// Login @Login
// @Summary User login
// @Description User login with JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Success 200 {string} string "User successfully logged in"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/login [get]
func Login(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	valid, err := PrivaseCheker(w, r, tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = w.Write([]byte(valid))
	if err != nil {
		return
	}

}

func PrivaseCheker(w http.ResponseWriter, r *http.Request, Usertoken string) (string, error) {

	req, err := Repository.Login(w, r, Usertoken)
	if err != nil {
		return "", err
	}
	return req, nil

}

// Register @Register
// @Summary Register
// @Tags Reg in service
// @Description Register a new user
// @Accept json
// @Produce json
// @Param input body NewUser true "User object for registration"
// @Success 200 {integer} integer 1
// @Failure 404 {error} http.Error
// @Failure 500 {error} http.Error
// @Router /api/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать запрос", http.StatusBadRequest)
	}
	var regData NewUser
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		http.Error(w, "Не удалось дессериализировать JSON", http.StatusBadRequest)
	}

	tokenString := TokenReqGenerate(w, r, bodyJSON)

	var tokenStr TokenString
	tokenStr.T = tokenString

	tokenJSON, err := json.Marshal(tokenStr)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(tokenJSON)
}
func TokenReqGenerate(w http.ResponseWriter, r *http.Request, User []byte) string {
	req, err := Repository.Register(w, r, User)
	if err != nil {
		log.Fatal(err)
	}

	var tokenstr TokenString

	tokenstr.T = req

	return req

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		valid, err := PrivaseCheker(w, r, tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		_, err = w.Write([]byte(valid))
		if err != nil {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// HandleSearch @Controller
// @Summary QueryGeocode
// @Tags geocode
// @Description create a search query
// @Accept json
// @Produce json
// @Param input body RequestAddressSearch true "query"
// @Success 200 {integer} integer 1
// @Failure 404 {error} http.Error
// @Failure 500 {error} http.Error
// @Router /search [post]
func HandleSearch(w http.ResponseWriter, r *http.Request) {

	resp, err := Service.Handle(w, r)
	if err != nil {
		fmt.Println(err)
	}

	bodyJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(bodyJSON)

}

// HandleGeocode @Controller
// @Summary QueryGeocode
// @Tags geocode
// @Description create a search query
// @Accept json
// @Produce json
// @Param input body RequestAddressSearch true "query"
// @Success 200 {integer} integer 1
// @Failure 404 {error} http.Error
// @Failure 500 {error} http.Error
// @Router /geocode [post]
func HandleGeocode(w http.ResponseWriter, r *http.Request) {

	resp, err := Service.Handle(w, r)
	if err != nil {
		fmt.Println(err)
	}

	bodyJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(bodyJSON)

}
