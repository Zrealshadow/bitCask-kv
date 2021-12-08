package global

import (
	bc "github.com/bitCaskKV/bitCask"
	"github.com/bitCaskKV/pkg/logger"
	"github.com/bitCaskKV/pkg/setting"
)

var (
	ServerSetting        *setting.ServerSettingS
	BitCaskSetting       *setting.BitCaskSettingS
	Logger               *logger.Logger
	LoggerSetting        *setting.LoggerSettingS
	DefaultBitCaskEngine *bc.BitCaskEngine
)
