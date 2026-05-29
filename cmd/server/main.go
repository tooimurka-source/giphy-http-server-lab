package main

import (
	"fmt"

	"lab3giphybot/internal/api/giphy"
	"lab3giphybot/internal/server"
	"lab3giphybot/pkg/config"
	"lab3giphybot/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Error(err)
		return
	}

	if cfg.GiphyKey == "" {
		logger.Error(fmt.Errorf("GIPHY_API_KEY is required"))
		return
	}

	router := server.NewRouter(giphy.NewClient(cfg.GiphyKey))
	addr := ":" + cfg.Port

	logger.Info("HTTP server started on " + addr)
	if err := router.Run(addr); err != nil {
		logger.Error(err)
	}
}
