package jwt

import (
	"exam3/config"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenJWT(m map[interface{}]interface{}) (string, string, error) {

	var (
		accesToken, refreshToken *jwt.Token
		claims                   jwt.MapClaims
	)

	accesToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)

	claims = accesToken.Claims.(jwt.MapClaims)
	rCaims := refreshToken.Claims.(jwt.MapClaims)

	for k, v := range m {
		claims[k.(string)] = v
		rCaims[k.(string)] = v
	}

	claims["iss"] = "user"
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().AddDate(0, 0, 1).Unix()

	rCaims["iss"] = "user"
	rCaims["iat"] = time.Now().Unix()
	rCaims["exp"] = time.Now().AddDate(0, 0, 10).Unix()

	accesTokenString, err := accesToken.SignedString(config.SignedKey)

	if err != nil {
		err = fmt.Errorf("acces_token generating error:%s", err)
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString(config.SignedKey)
	if err != nil {
		err = fmt.Errorf("refersh_token generating error:%s", err)

		return "", "", err
	}

	return accesTokenString, refreshTokenString, nil
}
func ExtractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return config.SignedKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}
