/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package server

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gzuidhof/starcloud/starcloud/handler/npmcdn"
	"github.com/gzuidhof/starcloud/starcloud/logger"
	"github.com/gzuidhof/starcloud/starcloud/middleware"
	"github.com/gzuidhof/starcloud/starcloud/npm"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func CreateCDNApp() (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		CaseSensitive:         true,
		DisableStartupMessage: true,
	})

	// app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(requestid.New(requestid.Config{
		ContextKey: "request_id",
	}))
	app.Use(middleware.NewLoggerMiddleware("starcloud", middleware.LoggerMiddlewareConfig{LogSuccesfulRequests: true}))
	app.Use(etag.New())

	cacheFolderPath := viper.GetString("cache_folder")
	fs := afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("cache_folder"))

	err := fs.MkdirAll(".", 0644)
	if err != nil {
		log.Fatalf("Could not create starcloud cache folder: %v", err)
	}

	// Seed some versions
	_, err = npm.DownloadPackageIntoFolder("starboard-notebook", "0.13.2", cacheFolderPath)
	if err != nil {
		log.Fatalf("Could not download seed starboard-notebook: %v", err)
	}

	cacheFs := afero.NewCacheOnReadFs(fs, afero.NewMemMapFs(), time.Minute*2)

	npmcdn := npmcdn.NewNPMCDNHandler("/npm/", afero.NewHttpFs(cacheFs))

	app.Get("/npm", func(ctx *fiber.Ctx) error {
		return ctx.SendString(`Hello there 🙋‍♂️, this small CDN only serves recent starboard-notebook versions with the correct headers.`)
	})

	app.Get("/npm/*", middleware.AddCommonCDNHeadersMiddleware, npmcdn.Handler)
	app.Head("/npm/*", middleware.AddCommonCDNHeadersMiddleware, npmcdn.Handler)
	app.Options("/npm/*", middleware.AddCommonCDNHeadersMiddleware)

	return app, nil
}

func StartCDNApp() {
	logger.SetupLogger()
	app, err := CreateCDNApp()

	if err != nil {
		log.Fatalf("Failed to start CDN app: %v", err)
	}

	zap.L().Info("Starting Starcloud CDN on :" + viper.GetString("port"))
	log.Fatal(app.Listen(":" + fmt.Sprint(viper.GetInt("port"))))
}
