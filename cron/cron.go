package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"super-signature/util/conf"
	"super-signature/util/tools"
)

const spec = "0 0 0 * * ?" // 每天0点执行
//const spec = "*/10 * * * * ?" //10秒执行一次，用于测试

func Init() {
	c := cron.New(cron.WithSeconds()) //支持到秒级别
	_, err := c.AddFunc(spec, GoRun)
	if err != nil {
		zap.S().Errorf("定时任务开启失败 %s", err)
	}
	c.Start()
	zap.L().Debug("定时任务已开启")
	select {}
}

func GoRun() {
	zap.L().Debug("----开始清理无用文件----")
	err := tools.Command("/bin/bash", "-c", fmt.Sprintf(
		`find %s -mtime +1 -name "*.*" -exec rm -rf {} \;`,
		conf.Config.ApplePath.TemporaryDownloadPath))
	if err != nil {
		zap.L().Error(err.Error())
	}
}
