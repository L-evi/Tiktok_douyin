package logset

import (
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"train-tiktok/common/tool"
)

func Handler(debug string, logConf logx.LogConf) {
	if debug == "true" {
		logConf.Level = "debug"
		logConf.Mode = "console"
	} else {
		logx.DisableStat()

		logConf.Level = "info"
		logConf.Mode = "file"
		logConf.KeepDays = 60
		logConf.Rotation = "daily"
		logConf.Encoding = "json"
		logConf.Path = "/app/logs"

		if err := tool.CheckPathOrCreate(logConf.Path); err != nil {
			log.Panicf("unable to create log path: %v", err)
		}
	}

	logx.MustSetup(logConf)
}
