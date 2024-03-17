package Repository

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
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

func Login(w http.ResponseWriter, r *http.Request, usrToken string) (string, error) {
	fmt.Println(usrToken)
	token, err := jwt.Parse(usrToken, func(token *jwt.Token) (interface{}, error) {
		// Проверка подписи методом HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if err != nil {
		http.Error(w, "Ошибка разбора токена", http.StatusBadRequest)
		return "", err
	}
	CheckForUsername := ""
	CheckForPassword := ""
	// Проверка, успешно ли разобран токен
	if token.Valid {
		// Получение клеймов из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return "", fmt.Errorf("ошибка получения клеймов из токена")
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
		return "", fmt.Errorf("некорректный токен")
	}
	user, _ := AuthUser[CheckForUsername]
	if user.Username != CheckForUsername || user.Password != CheckForPassword {
		return "", fmt.Errorf("неправильный пароль или имя пользователя")
	}
	userToken, _ := UserToken[CheckForUsername]
	if userToken.T != usrToken {
		return "", fmt.Errorf("неправильный токен")
	}

	return "Вы успешно авторизованы", nil
}
func Register(w http.ResponseWriter, r *http.Request, User []byte) (string, error) {
	var regData NewUser
	err := json.Unmarshal(User, &regData)
	if err != nil {
		return "", fmt.Errorf("неправильный формат учетных данных")
	}
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": regData.Username, "Password": regData.Password})
	if err != nil {
		return "", fmt.Errorf("ошибка генерации токена")
	}
	var tokenStr TokenString
	tokenStr.T = tokenString
	AuthUser = make(map[string]NewUser)
	AuthUser[regData.Username] = regData
	UserToken = make(map[string]TokenString)
	UserToken[regData.Username] = tokenStr

	return tokenString, nil
}
