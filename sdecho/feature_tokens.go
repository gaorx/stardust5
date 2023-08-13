package sdecho

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjwt"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type Token struct {
	UID  string `json:"uid,omitempty"`
	From string `json:"form,omitempty"`
	At   int64  `json:"at,omitempty"`
}

func tokenEncode(t Token, secret string) string {
	if lo.IsEmpty(t) {
		return ""
	}
	return lo.Must(sdjwt.Encode(secret, t))
}

func tokenDecode(s string, secrets []string) (Token, bool) {
	for _, secret := range secrets {
		t, err := sdjwt.Decode[Token](secret, s)
		if err == nil {
			return t, true
		}
	}
	return Token{}, false
}

type Tokens struct {
	Secrets    []string
	GetEncoded func(echo.Context) string
	IsExpired  func(echo.Context, Token) bool
}

const (
	keyTokens = "sdecho.tokens"
)

func (tt Tokens) Apply(app *echo.Echo) error {
	if len(tt.Secrets) <= 0 {
		return sderr.New("no tokens secret")
	}

	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			ec.Set(keyTokens, &tt)
			return next(ec)
		}
	}
	app.Use(middleware)
	return nil
}

func TokenDecode(_ context.Context, ec echo.Context) (Token, error) {
	tt := MustGet[*Tokens](ec, keyTokens)
	var encoded string
	if tt.GetEncoded != nil {
		encoded = tt.GetEncoded(ec)
	} else {
		encoded = ec.QueryParam("_token")
	}
	if encoded != "" {
		t, ok := tokenDecode(encoded, tt.Secrets)
		if !ok {
			return Token{}, sderr.WithStack(ErrDecodeToken)
		}
		if tt.IsExpired != nil {
			if tt.IsExpired(ec, t) {
				return t, sderr.WithStack(ErrTokenExpired)
			}
		}
		return t, nil
	} else {
		return Token{}, nil
	}
}

func TokenEncode(_ context.Context, ec echo.Context, t Token) string {
	tt := MustGet[*Tokens](ec, keyTokens)
	return tokenEncode(t, tt.Secrets[0])
}
