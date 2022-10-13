package middleware

import (
	"github.com/gin-gonic/gin"
	"ipashare/internal/model"
	"ipashare/pkg/caches"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

var (
	authKeyMu = make(map[string]*sync.Mutex)
	mu        sync.Mutex
)

func lockAuthKey(authKey string) {
	mu.Lock()
	defer mu.Unlock()
	if am, ok := authKeyMu[authKey]; ok {
		am.Lock()
	} else {
		m := sync.Mutex{}
		authKeyMu[authKey] = &m
		m.Lock()
	}
}

func unLockAuthKey(authKey string) {
	mu.Lock()
	defer mu.Unlock()
	if am, ok := authKeyMu[authKey]; ok {
		am.Unlock()
	}
}

func Key(store *model.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !caches.GetKeyInfo().EnableKey {
			c.Next()
			return
		}
		authKey := c.Request.Header.Get("Authorization")
		key, err := store.Key.Query(authKey)
		if err != nil || key == nil || key.Num <= 0 {
			c.Header("WWW-Authenticate", "Basic realm="+strconv.Quote("Authorization Required"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("authKey", authKey)
	}
}

func VerifyKey(store *model.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !caches.GetKeyInfo().EnableKey {
			c.Next()
			return
		}
		args := struct {
			AuthKey string `uri:"authKey" binding:"required"`
		}{}
		err := c.ShouldBindUri(&args)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		authKey, err := url.QueryUnescape(args.AuthKey)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		lockAuthKey(authKey)
		defer unLockAuthKey(authKey)
		key, err := store.Key.Query(authKey)
		if err != nil || key == nil || key.Num <= 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = store.Key.SubNum(key.Username)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
