package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/ccremer/web-powercycle/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v3"
)

var (
	indexHtml   = "index.html"
	executeHtml = "execute.html"
	cancelHtml  = "cancel.html"
)

type WebCommand struct {
	DryRunMode    bool
	SkipSudo      bool
	InsecureHttp  bool
	ListenAddress string
	AuthUser      string
	AuthPass      string
	CertFilePath  string
	CertKeyPath   string
}

type Renderer struct {
	templates map[string]*template.Template
}

type Values struct {
	Hostname     string
	ErrorMessage string
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return r.templates[name].Execute(w, data)
}

func (c *WebCommand) StartWeb(_ context.Context, _ *cli.Command) error {

	server := echo.New()
	defer func(server *echo.Echo) {
		err := server.Close()
		if err != nil {
			panic(err)
		}
	}(server)
	server.HideBanner = true
	server.HidePort = true
	server.TLSServer.ErrorLog = stdLogger
	server.StdLogger = stdLogger

	if c.AuthPass == "" || c.AuthUser == "" {
		return fmt.Errorf("required flags \"%s\" or \"%s\" not set", newAuthUserFlag(nil).Name, newAuthPassFlag(nil).Name)
	}

	server.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(c.AuthUser)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(c.AuthPass)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		Skipper:     skipAccessLogs,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slogger.LogAttrs(context.Background(), slog.LevelInfo, c.Request().Method,
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slogger.LogAttrs(context.Background(), slog.LevelError, c.Request().Method,
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	renderer := &Renderer{
		templates: map[string]*template.Template{
			indexHtml:   getTemplate(indexHtml),
			executeHtml: getTemplate(executeHtml),
			cancelHtml:  getTemplate(cancelHtml),
		},
	}
	server.Renderer = renderer

	var shutdownHandler ShutdownHandler = &ExecutableShutdown{
		Logger:   slogger,
		SkipSudo: c.SkipSudo,
	}
	if c.DryRunMode {
		slogger.Debug("Using Dry-run shutdown handler")
		shutdownHandler = &DryRunShutdown{
			Logger: slogger,
		}
	}

	server.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, indexHtml, getData())
	})
	server.POST("/execute", func(c echo.Context) error {
		data := getData()
		err := shutdownHandler.ShutDownDelayed(1)
		if err != nil {
			slogger.Error(err.Error())
			data.ErrorMessage = err.Error()
		}
		return c.Render(http.StatusOK, executeHtml, data)
	})

	server.GET("/cancel", func(c echo.Context) error {
		data := getData()
		err := shutdownHandler.CancelShutdown()
		if err != nil {
			slogger.Error(err.Error())
			data.ErrorMessage = err.Error()
		}
		return c.Render(http.StatusOK, cancelHtml, data)
	})

	if c.InsecureHttp {
		if !c.DryRunMode {
			return fmt.Errorf("insecure-http flag is only allowed in dry-run mode")
		}
		slogger.Info("Starting server", "address", c.ListenAddress)
		return server.Start(c.ListenAddress)
	}
	slogger.Info("Starting server", "address", c.ListenAddress, "cert", c.CertFilePath, "key", c.CertKeyPath)
	return server.StartTLS(c.ListenAddress, c.CertFilePath, c.CertKeyPath)
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

func getData() Values {
	hostname, _ := os.Hostname()
	return Values{
		Hostname: hostname,
	}
}

func getTemplate(name string) *template.Template {
	return template.Must(template.ParseFS(templates.PublicFs, "layout.html", name))
}
