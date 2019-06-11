package routers

import (
	"work/pkg/setting"

	"work/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	bill := r.Group("/api")
	{
		//获取销项数据
		bill.POST("/test1", controller.AddInvoiceAllInfo)

		bill.POST("/add", controller.AddBill)

		bill.POST("/test", controller.Test)

	}

	return r
}
