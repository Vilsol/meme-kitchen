package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger() fiber.Handler {
	sublog := log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		msg := "Request"
		if err != nil {
			msg = err.Error()
		}

		code := c.Response().StatusCode()

		dumplogger := sublog.With().
			Int("status", code).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Str("latency", time.Since(start).String()).
			Str("user-agent", c.Get(fiber.HeaderUserAgent)).
			Logger()

		switch {
		case code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError:
			dumplogger.Warn().Msg(msg)
		case code >= http.StatusInternalServerError:
			dumplogger.Error().Msg(msg)
		default:
			dumplogger.Info().Msg(msg)
		}

		return err
	}
}
