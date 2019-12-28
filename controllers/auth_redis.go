package ct

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xeuus/gt/pkg/rds"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xeuus/gt/pkg/jwt"
)

type RedisAuthHelper struct {
	JWT     jwt.Authenticator
	Conn    redis.Conn
	REDIS   rds.Redis
	Prefix  string
	Timeout int
	KeyGen  func() string
}

func (r *RedisAuthHelper) Throttle(phoneNumber string, fn func() error) (int, error) {
	key := r.Prefix + "::throttle::" + phoneNumber
	ttl, err := redis.Int(r.Conn.Do("TTL", key))
	if err != nil {
		return 0, err
	}
	if ttl > 0 {
		return ttl, nil
	}
	err = fn()
	if err != nil {
		errRem := r.Conn.Send("DEL", key)
		if errRem != nil {
			return 0, errRem
		}
		return 0, err
	}
	if err := r.Conn.Send("SET", key, 1); err != nil {
		return 0, err
	}
	if err := r.Conn.Send("EXPIRE", key, r.Timeout); err != nil {
		return 0, err
	}
	return r.Timeout, nil
}

func (r *RedisAuthHelper) Get(phoneNumber string) (string, error) {
	return redis.String(r.Conn.Do("GET", r.getKey(phoneNumber)))
}

func (r *RedisAuthHelper) Request(phoneNumber string, send func(string) error) (int, error) {
	key := r.getKey(phoneNumber)
	ttl, err := redis.Int(r.Conn.Do("TTL", key))
	if err != nil {
		return 0, err
	}
	if ttl > 0 {
		code, _ := redis.String(r.Conn.Do("GET", key))
		log.Println("code for user " + phoneNumber + " is " + code)
		return ttl, nil
	}
	code := r.KeyGen()
	log.Println("code for user " + phoneNumber + " is " + code)
	err = send(code)
	if err != nil {
		errRem := r.Remove(phoneNumber)
		if errRem != nil {
			return 0, errRem
		}
		return 0, err
	}
	if err := r.Conn.Send("SET", key, code); err != nil {
		return 0, err
	}
	if err := r.Conn.Send("EXPIRE", key, r.Timeout); err != nil {
		return 0, err
	}
	return r.Timeout, nil

}

func (r *RedisAuthHelper) Remove(phoneNumber string) error {
	key := r.getKey(phoneNumber)
	err := r.Conn.Send("DEL", key)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAuthHelper) Create(userID string, clientName string, gateway string, duration time.Duration, flags int64) (string, error) {
	expAt := time.Now().Add(duration)
	if clientName == "" || gateway == "" {
		return "", fmt.Errorf("invalid headers")
	}

	token, err := r.JWT.Create(&jwt.AuthClaims{
		UserId:     userID,
		Gateway:    gateway,
		ClientName: clientName,
		ExpiresAt:  expAt,
		Flags:      flags,
	})
	if err != nil {
		return "", err
	}
	key := r.getToken(userID)
	if err := r.Conn.Send("HSET", key, clientName, expAt.Format(time.RFC3339)); err != nil {
		return "", err
	}
	return token, nil
}

func (r *RedisAuthHelper) Check(token, clientName, gateway string, flag int64) (int64, *jwt.AuthClaims, error) {
	claims, err := r.JWT.Parse(token)
	if err != nil {
		return 0, nil, fmt.Errorf("parse failed, %s", err.Error())
	}
	if claims.ClientName != clientName || claims.Gateway != gateway {
		return 0, nil, fmt.Errorf("invalid headers")
	}
	key := r.getToken(claims.UserId)
	tokens, err := redis.Strings(r.Conn.Do("HGETALL", key))
	if err != nil {
		return 0, nil, fmt.Errorf("failed getting token, %s", err.Error())
	}
	found := false
	for i := 0; i < len(tokens); i++ {
		clientName := tokens[i]
		i++
		expAt, err := time.Parse(time.RFC3339, tokens[i])
		if err != nil || expAt.Before(time.Now()) {
			err = r.Conn.Send("HDEL", key, clientName)
			if err != nil {
				return 0, nil, fmt.Errorf("failed loading cache, %s", err.Error())
			}
			continue
		}
		if claims.ClientName == clientName {
			found = true
		}
	}
	if found {
		if flag != 0 && (claims.Flags&flag) == 0 {
			return 0, nil, fmt.Errorf("no right to access")
		}

		id, err := strconv.ParseInt(claims.UserId, 10, 64)
		if err != nil {
			return 0, nil, fmt.Errorf("failed parsing user, %s", err.Error())
		}
		return id, claims, nil
	}
	return 0, nil, fmt.Errorf("invalid token")
}

func (r *RedisAuthHelper) getKey(phoneNumber string) string {
	return r.Prefix + "::otp::" + phoneNumber
}

func (r *RedisAuthHelper) getToken(phoneNumber string) string {
	return r.Prefix + "::token::" + phoneNumber
}
func (r *RedisAuthHelper) Middleware(ctx *gin.Context) {
	bearer := ctx.GetHeader("Authorization")
	if len(bearer) > 7 && bearer[0:7] == "Bearer " {
		bearer = bearer[7:]
		clientName := ctx.GetHeader("Client-Name")
		gateway := ctx.GetHeader("App-Gateway")
		if clientName == "" || gateway == "" {
			panic("headers not set")
		}
		userID := int64(0)
		err := r.REDIS.Action(func(c redis.Conn) error {
			r.Conn = c
			var err error
			userID, _, err = r.Check(bearer, clientName, gateway, 0)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic("token validation failed, " + err.Error())
		}
		ctx.Set("userID", userID)
		ctx.Set("gateway", gateway)
		ctx.Next()
		return
	}
	panic("invalid token")
}
