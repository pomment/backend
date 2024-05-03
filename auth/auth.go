package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"pomment-go/common"
	"pomment-go/config"
	"time"
)

func FindUserByName(username interface{}) (*common.PommentConfigAdminUser, error) {
	var user *common.PommentConfigAdminUser

	// 查询用户是否存在
	for _, e := range config.Content.Admin.User {
		if e.Name == username {
			user = &e
			break
		}
	}
	if user == nil {
		return nil, errors.New("unknown user")
	}
	return user, nil
}
func CheckPassword(username string, password string) error {
	user, err := FindUserByName(username)
	if err != nil {
		return err
	}

	// 校验密码
	hash := &user.Password
	err = bcrypt.CompareHashAndPassword([]byte(*hash), []byte(password))
	return err
}

func GenerateToken(username string) (token string, err error) {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": username,
		"logged":   "true",
		"time":     time.Now().Unix(),
	})

	token, err = tokenObj.SignedString([]byte(config.Content.Admin.Salt))
	return token, err
}

func ValidateToken(token string) (username string, err error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.Content.Admin.Salt), nil
	})

	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if !(ok && tokenObj.Valid) {
		return "", err
	}

	_, err = FindUserByName(claims["username"])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", claims["username"]), nil
}
