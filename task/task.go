package task

import (
	"log"
	"net/http"
	"net/url"
	"work/pkg/logging"
	"work/pkg/setting"

	"github.com/robfig/cron"
)

func Task() {
	specForVlan := GetTaskConfig("Cert")
	go NewCronTask(CheckCertResult, specForVlan)

	specForVlan = GetTaskConfig("Input")
	go NewCronTask(GuoPiaoInputData, specForVlan)

	specForVlan = GetTaskConfig("Detail")
	go NewCronTask(GuoPiaoGetDetailInfo, specForVlan)

}

func NewCronTask(funcName func(), specForVlan string) *cron.Cron {
	c := cron.New()
	c.AddFunc(specForVlan, func() {
		funcName()
	})
	c.Start()
	return c
}

func CheckCertResult() {
	logging.Info("CheckCertResult start......")
	host := GetTaskConfig("Host")
	posturl := host + "/task/CheckCertResult"
	http.PostForm(posturl, url.Values{})
}

func GuoPiaoInputData() {
	logging.Info("GuoPiaoInputData start......")
	host := GetTaskConfig("Host")
	posturl := host + "/task/GuoPiaoInputData"
	http.PostForm(posturl, url.Values{})
}

func GuoPiaoGetDetailInfo() {
	logging.Info("GuoPiaoGetDetailInfo start......")
	host := GetTaskConfig("Host")
	posturl := host + "/task/GuoPiaoGetDetailInfo"
	http.PostForm(posturl, url.Values{})
}

func GetTaskConfig(task string) string {
	sec, err := setting.Cfg.GetSection("task")
	if err != nil {
		log.Fatal(2, "Fail to get section 'task':%v", err)
	}
	conf := sec.Key(task).String()
	return conf
}
