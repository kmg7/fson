package server

import (
	"github.com/kmg7/fson/internal/server/config"
	"github.com/kmg7/fson/internal/server/transfer"
)

func ConfigServerStart() {
	config.StartConfigServer("localhost:8080")
}

func TransferServerStart() {
	transfer.Start("localhost:8081")
}
