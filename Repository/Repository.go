package Repository

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth/v5"
	"io/ioutil"
	"net/http"
	"strings"
)

type TokenString struct {
	T string `json:"token"`
}
type NewUser struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

var AuthUser map[string]NewUser
var UserToken map[string]TokenString
var tokenAuth *jwtauth.JWTAuth

func Login(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	fmt.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка подписи методом HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if err != nil {
		http.Error(w, "Ошибка разбора токена", http.StatusBadRequest)
		return
	}
	CheckForUsername := ""
	CheckForPassword := ""
	// Проверка, успешно ли разобран токен
	if token.Valid {
		// Получение клеймов из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.Write([]byte("Ошибка получения клеймов из токена"))
			return
		}

		// Вывод клеймов
		for key, value := range claims {

			if key == "Username" {
				username, _ := value.(string)
				CheckForUsername = username
			}
			if key == "Password" {
				password, _ := value.(string)
				CheckForPassword = password
			}
		}

	} else {
		w.Write([]byte("Некорректный токен"))
	}
	user, _ := AuthUser[CheckForUsername]
	if user.Username != CheckForUsername || user.Password != CheckForPassword {
		w.Write([]byte("Неправильный пароль или имя пользователя"))
		return
	}
	userToken, _ := UserToken[CheckForUsername]
	if userToken.T != tokenString {
		w.Write([]byte("неправильный токен"))
		return
	}

	w.Write([]byte("Вы успешно авторизованы"))
}
func Register(w http.ResponseWriter, r *http.Request) {
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var regData NewUser
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		fmt.Println(err)
	}
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": regData.Username, "Password": regData.Password})
	if err != nil {
		fmt.Println(err)
	}
	var tokenStr TokenString
	tokenStr.T = tokenString
	AuthUser = make(map[string]NewUser)
	AuthUser[regData.Username] = regData
	UserToken = make(map[string]TokenString)
	UserToken[regData.Username] = tokenStr

	tokenJSON, err := json.Marshal(tokenStr)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(tokenJSON)
}
