package main

import (
	"context"
	"os"
	"sync"

	"github.com/chaunsin/netease-cloud-music/api"
	"github.com/chaunsin/netease-cloud-music/api/weapi"
	"github.com/chaunsin/netease-cloud-music/config"
	"github.com/chaunsin/netease-cloud-music/pkg/log"
)

var (
	globalApp     *App
	globalAppOnce sync.Once
)

type App struct {
	client *api.Client
	api    *weapi.Api
	cfg    *config.Config
	l      *log.Logger
}

func ensureApp() *App {
	globalAppOnce.Do(func() {
		globalApp = &App{}
		globalApp.init()
	})
	return globalApp
}

func (a *App) init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	cfg := config.GetDefault()
	cfg.ReplaceMagicVariables("HOME", home)
	cfg.Log.Stdout = false
	cfg.Log.Level = "info"
	a.cfg = cfg

	logger := log.New(cfg.Log)
	log.Default = logger
	a.l = logger

	client, err := api.NewClient(cfg.Network, logger)
	if err != nil {
		log.Error("init api client: %s", err)
		return
	}
	a.client = client
	a.api = weapi.New(client)
}

func (a *App) Close() {
	if a.client != nil {
		_ = a.client.Close(context.Background())
	}
	if a.l != nil {
		_ = a.l.Close()
	}
}
