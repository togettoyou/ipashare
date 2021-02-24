package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"super-signature/pkg/setting"
	"super-signature/pkg/util"
)

const spec = "0 0 0 * * ?" // 每天0点执行
//const spec = "*/10 * * * * ?" //10秒执行一次，用于测试

func Init() {
	c := cron.New(cron.WithSeconds()) //支持到秒级别
	_, err := c.AddFunc(spec, GoRun)
	if err != nil {
		log.Fatalf("定时任务开启失败 %s", err)
	}
	c.Start()
	log.Printf("定时任务已开启")
	select {}
}

func GoRun() {
	log.Println("----开始清理无用文件----")
	err := util.RunCmd(fmt.Sprintf(
		`find %s -mtime +1 -name "*.*" -exec rm -rf {} \;`,
		setting.PathSetting.TemporaryDownloadPath))
	if err != nil {
		log.Printf("%s", err.Error())
	}
}
