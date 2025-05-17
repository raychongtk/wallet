//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/raychongtk/wallet/datastore"
	"github.com/raychongtk/wallet/repository"
	"github.com/raychongtk/wallet/service"
)

import "github.com/google/wire"

func injectRoutes(ctx context.Context) (*gin.Engine, error) {
	panic(wire.Build(
		datastore.WireSet,
		repository.WireSet,
		service.WireSet,
	))
}
