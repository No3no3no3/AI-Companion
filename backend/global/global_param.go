package global

import (
	"github.com/ai-companion/backend/internal/pkg/config"
	"github.com/google/uuid"
)

var (
	Cfg  *config.Config //全局配置变量
	UUID uuid.UUID
)

func init() {
	Cfg = config.Load()
	UUID = uuid.New()
}
