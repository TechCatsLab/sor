/*
 * Revision History:
 *     Initial: 2018/08/15        Shi Ruitao
 */

package filter

import (
	"errors"
	"fmt"
	"strings"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/sor/base/constants"
)

const (
	claimUID = "uid"
	claimAdmin = "admin"
)

func GetUID(ctx *server.Context, tokenHMACKey string) (uint, error) {
	claims, err := checkJWT(ctx, tokenHMACKey)
	if err != nil {
		return constants.InvalidUID, err
	}

	rawUID := uint(claims[claimUID].(float64))

	return rawUID, nil
}

func GetAdmin(ctx *server.Context, tokenHMACKey string) (bool, error) {
	claims, err := checkJWT(ctx, tokenHMACKey)
	if err != nil {
		return constants.InvalidAdmin, err
	}

	isAdmin := claims[claimAdmin].(bool)

	return isAdmin, nil
}

// checkJWT check whether the token is valid, it returns claims if valid.
func checkJWT(ctx *server.Context, tokenHMACKey string) (jwtgo.MapClaims, error) {
	var (
		err error
	)

	authString := ctx.Request().Header.Get("Authorization")
	kv := strings.Split(authString, " ")

	if len(kv) != 2 || kv[0] != "Bearer" {
		err = errors.New("invalid token authorization string")
		return nil, err
	}

	tokenString := kv[1]

	token, _ := jwtgo.Parse(tokenString, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenHMACKey), nil
	})

	claims, ok := token.Claims.(jwtgo.MapClaims)

	if !ok || !token.Valid {
		err = errors.New("invalid token")
		return nil, err
	}

	return claims, nil
}
