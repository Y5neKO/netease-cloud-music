package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	appImpl := application.New(application.Options{
		Name:        "NCM GUI",
		Description: "网易云音乐 GUI 客户端",
		Assets: application.AssetOptions{
			FS: assets,
		},
		Bind: []any{
			&App{},
			&AuthService{},
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	appImpl.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "NCM GUI",
		Width:  960,
		Height: 640,
		URL:    "/",
	})

	if err := appImpl.Run(); err != nil {
		log.Fatal(err)
	}
}
