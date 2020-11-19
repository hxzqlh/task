package auth

import (
	"crypto/rsa"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/dgrijalva/jwt-go/test"
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
)

type Claims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func NewPlugin() plugin.Plugin {
	var pubKey *rsa.PublicKey
	return plugin.NewPlugin(
		plugin.WithName("auth"),
		plugin.WithFlag(
			&cli.StringFlag{
				Name:  "auth_key",
				Usage: "auth key file",
				Value: "./conf/public.key",
			}),
		plugin.WithInit(func(ctx *cli.Context) error {
			pubKeyFile := ctx.String("auth_key")
			pubKey = test.LoadRSAPublicKeyFromDisk(pubKeyFile)
			return nil
		}),
		plugin.WithHandler(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var claims Claims
				token, err := request.ParseFromRequest(r,
					request.AuthorizationHeaderExtractor,
					func(*jwt.Token) (interface{}, error) {
						return pubKey, nil
					},
					request.WithClaims(&claims),
				)

				if err != nil {
					log.Print("token invalid: ", err.Error())
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// token.Valid是否成功，取决于jwt中Claims接口定义的Valid() error方法
				// 本例中我们直接使用了默认Claims实现jwt.StandardClaims提供的方法，实际生产中可以根据需要重写
				if token == nil || !token.Valid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// todo:虽然是有效的token，但并不意味着此用户有权访问所有接口，演示代码省略鉴权细节

				r.Header.Set("userId", claims.UserId)

				h.ServeHTTP(w, r)
			})
		}),
	)
}
