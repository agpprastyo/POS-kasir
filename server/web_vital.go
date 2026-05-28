package server

import (
	appmetrics "POS-kasir/pkg/metrics"

	"github.com/gofiber/fiber/v3"
)

// WebVitalsHandler handles incoming Web Vitals data from the frontend.
func WebVitalsHandler() fiber.Handler {
	type webVitalEntry struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}

	return func(c fiber.Ctx) error {
		var entries []webVitalEntry
		if err := c.Bind().JSON(&entries); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
		}

		if len(entries) > 10 {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": "payload too large"})
		}

		for _, entry := range entries {
			if entry.Value < 0 || entry.Value > 100000 { // basic validation (max 100s)
				continue
			}

			switch entry.Name {
			case "LCP":
				appmetrics.WebVitalsLCP.Observe(entry.Value)
			case "FCP":
				appmetrics.WebVitalsFCP.Observe(entry.Value)
			case "CLS":
				appmetrics.WebVitalsCLS.Observe(entry.Value)
			case "INP":
				appmetrics.WebVitalsINP.Observe(entry.Value)
			case "TTFB":
				appmetrics.WebVitalsTTFB.Observe(entry.Value)
			}
		}

		return c.Status(fiber.StatusNoContent).Send(nil)
	}
}
