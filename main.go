package main

import (
	"embed"
	"runtime"

	appsvc "github.com/DlnKot/arc/internal/app"
	"github.com/DlnKot/arc/internal/config"
	"github.com/DlnKot/arc/internal/launchers"
	"github.com/DlnKot/arc/internal/logging"
	"github.com/DlnKot/arc/internal/network"
	"github.com/DlnKot/arc/internal/store"
	"github.com/DlnKot/arc/internal/updater"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logger, err := logging.New(config.AppName)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	storeSvc := store.New(config.AppName, logger)
	networkSvc := network.New(logger)
	launcherSvc := launchers.New(logger)
	updaterSvc := updater.New(logger)
	app := appsvc.New(storeSvc, networkSvc, launcherSvc, updaterSvc, logger)

	appOptions := &options.App{
		Width:     1440,
		Height:    900,
		MinWidth:  1180,
		MinHeight: 760,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 242, G: 243, B: 245, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	}

	if runtime.GOOS == "darwin" {
		appOptions.Title = config.AppName
		appOptions.Mac = &mac.Options{
			TitleBar: mac.TitleBarHidden(),
		}
	} else {
		appOptions.Title = config.AppName
	}

	err = wails.Run(appOptions)

	if err != nil {
		println("Error:", err.Error())
	}
}
