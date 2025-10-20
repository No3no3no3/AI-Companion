package websocket

import (
	"github.com/ai-companion/backend/internal/infrastructure/websocket/service"
	"github.com/ai-companion/backend/internal/pkg/config"
)

func StartWebSocketServer() {

	config.Load()

	service.StartWebSocket()

}
