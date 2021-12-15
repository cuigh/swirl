package security

import (
	"errors"
	"strings"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/misc"
	"github.com/golang-jwt/jwt"
)

var ErrNoNeedRefresh = errors.New("no need to refresh")

type JWT struct {
	Schema      string
	Sources     data.Options
	KeyFunc     jwt.Keyfunc
	Identifier  func(token *jwt.Token) web.User
	tokenExpiry int64
	logger      log.Logger
}

func NewIdentifier() web.Filter {
	logger := log.Get("security")
	key := misc.Options.TokenKey
	expiry := misc.Options.TokenExpiry
	if key == "" {
		key = "swirl"
		logger.Warnf("Swirl is using default token key as token_key isn't configured, this may cause security problems")
	}
	if expiry == 0 {
		expiry = 30 * time.Minute
	}

	return &JWT{
		logger:      logger,
		tokenExpiry: int64(expiry.Seconds()),
		Schema:      "Bearer",
		Sources: data.Options{
			{Name: "header", Value: web.HeaderAuthorization},
		},
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			// TODO: use user salt as key?
			return []byte(key), nil
		},
		Identifier: func(token *jwt.Token) web.User {
			claims := token.Claims.(jwt.MapClaims)
			return security.NewUser(cast.ToString(claims["sub"]), cast.ToString(claims["name"]))
		},
	}
}

func (j *JWT) Apply(next web.HandlerFunc) web.HandlerFunc {
	if j.KeyFunc == nil {
		panic("KeyFunc is required")
	}
	if j.Schema == "" {
		j.Schema = "Bearer"
	}
	if len(j.Sources) == 0 {
		j.Sources = data.Options{
			{Name: "header", Value: web.HeaderAuthorization},
		}
	}
	if j.Identifier == nil {
		j.Identifier = func(token *jwt.Token) web.User {
			claims := token.Claims.(jwt.MapClaims)
			return security.NewUser(cast.ToString(claims["sub"]), cast.ToString(claims["name"]))
		}
	}

	return func(ctx web.Context) error {
		ts := j.extractToken(ctx)
		if ts != "" {
			token, err := jwt.Parse(ts, j.KeyFunc)
			if err != nil {
				j.logger.Debugf("failed to parse token: %s", err)
			} else {
				user := j.Identifier(token)
				ctx.SetUser(user)
				if ts, err = j.refreshToken(user, token); err == nil {
					ctx.SetHeader(web.HeaderAuthorization, ts)
				} else if err != ErrNoNeedRefresh {
					j.logger.Errorf("failed to refresh token: %s", err)
				}
			}
		}

		return next(ctx)
	}
}

func (j *JWT) extractToken(ctx web.Context) (token string) {
	for _, src := range j.Sources {
		switch src.Name {
		case "header":
			token = ctx.Header(src.Value)
			if strings.HasPrefix(token, j.Schema) {
				return token[len(j.Schema)+1:]
			}
		case "cookie":
			if cookie, err := ctx.Cookie(src.Value); err == nil {
				token = cookie.Value
			}
		case "form":
			token = ctx.Form(src.Value)
		case "query":
			token = ctx.Query(src.Value)
		}
		if token != "" {
			return
		}
	}
	return
}

func (j *JWT) refreshToken(user web.User, token *jwt.Token) (string, error) {
	claims := token.Claims.(jwt.MapClaims)
	expiry := cast.ToInt64(claims["exp"])
	now := time.Now().Unix()
	// refresh token when remaining expiry is less than 5 minutes
	if (expiry - now) < 5*60 {
		ts, err := j.CreateToken(user.ID(), user.Name())
		if err != nil {
			return "", err
		}
		return ts, nil
	}
	return "", ErrNoNeedRefresh
}

func (j *JWT) CreateToken(id, name string) (string, error) {
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"name": name,
		"sub":  id,
		"iat":  now,
		"exp":  now + j.tokenExpiry,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key, err := j.KeyFunc(token)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}
