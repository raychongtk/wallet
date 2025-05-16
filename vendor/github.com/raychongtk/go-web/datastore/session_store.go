package datastore

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/raychongtk/go-web/util"
	"net/http"
)

func ProvideSessionStore() redis.Store {
	var config *SessionConfig
	util.Load("dev", ".", &config)
	store, _ := redis.NewStore(10, "tcp", config.SessionAddress, config.SessionPassword, []byte(config.SessionKey))
	// add secure cookie if you use https
	// for demo purpose, i only have http protocol
	store.Options(sessions.Options{HttpOnly: true, SameSite: http.SameSiteLaxMode, Domain: "localhost", Path: "/", MaxAge: 3600})
	return store
}

type SessionConfig struct {
	SessionAddress  string `mapstructure:"SESSION_ADDRESS"`
	SessionPassword string `mapstructure:"SESSION_PASSWORD"`
	SessionKey      string `mapstructure:"SESSION_KEY"`
}
