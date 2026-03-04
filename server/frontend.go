package server

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

// SetupFrontend serves the SPA static files from the web/dist directory.
// It handles:
//  1. Static assets (JS, CSS, images, etc.) via static middleware
//  2. SPA catch-all: any route not matching /api/*, /swagger/*, or /healthz
//     returns index.html so client-side routing works.
func SetupFrontend(app *App) {
	distPath := "./web/dist"

	// Check if the dist folder exists (it won't exist during dev without a build)
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		app.Logger.Warn("Frontend dist directory not found, skipping frontend setup. Run 'cd web && npm run build' to generate it.")
		return
	}

	// Serve static files (JS, CSS, images, fonts, etc.)
	app.FiberApp.Use("/", static.New(distPath, static.Config{
		Compress: true,
		Browse:   false,
	}))

	// SPA catch-all: return index.html for any route that is not an API route,
	// swagger, or health check. This allows client-side routing to handle the path.
	app.FiberApp.Get("/*", func(c fiber.Ctx) error {
		path := c.Path()

		// Skip API, swagger, and health check routes
		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/swagger/") ||
			path == "/healthz" {
			return c.Next()
		}

		return c.SendFile(distPath + "/index.html")
	})
}
