package main

import (
	"context"
	"crypto/subtle"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/ccremer/web-powercycle/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v2"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func StartWeb(c *cli.Context) error {

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		Skipper:     skipAccessLogs,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, c.Request().Method,
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, c.Request().Method,
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})
	e.POST("/execute", func(c echo.Context) error {
		logger.Info("SHUTTING DOWN")
		return c.Render(http.StatusOK, "done.html", nil)
	})
	publicFs := getFileSystem()
	assetHandler := http.FileServer(publicFs)
	e.GET("/", echo.WrapHandler(assetHandler))

	logger.Info("Starting server", "port", ":7443")
	return e.Start(":7443")
}

var publicRoutes = map[string]bool{
	"/favicon.ico": true,
	"/robots.txt":  true,
}

func skipAccessLogs(ctx echo.Context) bool {
	// given an exact known key, lookups in maps are faster than iterating over slices.
	_, exists := publicRoutes[ctx.Request().URL.Path]
	return exists
}

func getFileSystem() http.FileSystem {
	if _, err := os.Stat("templates"); err == nil {
		return http.Dir("templates")
	}
	return http.FS(templates.PublicFs)
}
