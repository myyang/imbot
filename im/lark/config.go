package lark

import (
	"os"

	lCfg "github.com/larksuite/oapi-sdk-go/core/config"
	lConstants "github.com/larksuite/oapi-sdk-go/core/constants"
	lLog "github.com/larksuite/oapi-sdk-go/core/log"
)

var larkConfig = lCfg.NewConfigWithDefaultStore(
	lConstants.DomainLarkSuite,
	lCfg.NewInternalAppSettings(
		os.Getenv("LARK_APP_ID"),
		os.Getenv("LARK_APP_SECRET"),
		os.Getenv("LARK_MSG_TOKEN"),
		os.Getenv("LARK_ENCRYPT_TOKEN"),
	),
	lLog.NewDefaultLogger(),
	lLog.LevelInfo,
)
