package server

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func SetupFrontend(app *App) {
	distPath := "./web/dist"

	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		app.Logger.Warn("Frontend dist directory not found, skipping frontend setup. Run 'cd web && npm run build' to generate it.")
		return
	}

	app.FiberApp.Use("/", static.New(distPath, static.Config{
		Compress: true,
		Browse:   false,
	}))

	app.FiberApp.Get("/*", func(c fiber.Ctx) error {
		path := c.Path()

		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/swagger/") ||
			path == "/healthz" {
			return c.Next()
		}

		return c.SendFile(distPath + "/index.html")
	})
}
