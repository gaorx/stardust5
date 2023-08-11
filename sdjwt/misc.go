package sdjwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/lo"

	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
)

var (
	signingMethod = jwt.SigningMethodHS256
)

func Encode(secret string, payload any) (string, error) {
	payload1, err := sdjson.StructToObject(payload)
	if err != nil {
		return "", sderr.Wrap(err, "struct to claims error")
	}
	rawToken := jwt.NewWithClaims(signingMethod, jwt.MapClaims(payload1))
	signedToken, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", sderr.Wrap(err, "encode jwt token error")
	}
	return signedToken, nil
}

func Decode[T any](secret string, signedToken string) (T, error) {
	var payload0 = map[string]any{}
	_, err := jwt.ParseWithClaims(signedToken, jwt.MapClaims(payload0), func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, sderr.NewWith("unexpected signing method", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return lo.Empty[T](), sderr.Wrap(err, "decode jwt token error")
	}
	return sdjson.ObjectToStruct[T](payload0)
}
