package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestAddressSearch struct {
	Query string `json:"query"`
}
type TokenString struct {
	T string `json:"token"`
}

var token string

type NewUser struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

func main() {
	client := &http.Client{}
	// Отправляем запрос на регистрацию пользователя и получаем токен
	if token == "" {
		newUser := NewUser{Username: "user123", Password: "123"}
		userJSON, err := json.Marshal(newUser)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp, err := client.Post("http://localhost:8080/api/register", "application/json", bytes.NewReader(userJSON))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		var response TokenString

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Println(err)
			return
		}

		token = response.T
	}
	fmt.Println(token)
	req, err := http.NewRequest("GET", "http://localhost:8080/api/login", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyByte))

	// Устанавливаем токен в запрос
	query := RequestAddressSearch{Query: "мск садовая 20"}
	queryJSON, err := json.Marshal(query)
	if err != nil {
		fmt.Println(err)
	}
	req, err = http.NewRequest("POST", "http://localhost:8080/api/address/geocode", bytes.NewReader(queryJSON))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Отправляем запрос на сервер
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	bodyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var datainfo RequestAddressSearch
	err = json.Unmarshal(bodyJSON, &datainfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(datainfo)
}
