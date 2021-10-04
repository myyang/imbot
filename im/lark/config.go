package lark

import (
	lCfg "github.com/larksuite/oapi-sdk-go/core/config"
	lConstants "github.com/larksuite/oapi-sdk-go/core/constants"
	lLog "github.com/larksuite/oapi-sdk-go/core/log"

	botCfg "github.com/myyang/imbot/config"
)

var larkConfig = lCfg.NewConfigWithDefaultStore(
	lConstants.DomainLarkSuite,
	lCfg.NewInternalAppSettings(
		botCfg.Config.GetString("LARK_APP_ID"),
		botCfg.Config.GetString("LARK_APP_SECRET"),
		botCfg.Config.GetString("LARK_MSG_TOKEN"),
		botCfg.Config.GetString("LARK_ENCRYPT_TOKEN"),
	),
	lLog.NewDefaultLogger(),
	lLog.LevelInfo,
)
