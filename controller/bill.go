package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"work/models"
	"work/pkg/e"
	"work/pkg/gredis"
	"work/pkg/logging"

	mcache "github.com/OrlovEvgeny/go-mcache"

	"github.com/gin-gonic/gin"
	"github.com/kirinlabs/utils"
	"github.com/kirinlabs/utils/str"
)

func AddInvoiceAllInfo(c *gin.Context) {
	data := make(map[string]interface{})
	err := c.BindJSON(&data)

	code := e.SUCCESS
	if err != nil {
		code = e.ERROR
		logging.Info(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": e.GetMsg(code),
			"code":    code,
			"data":    err,
		})
		return
	}

	//获取state字段的值
	message, result := getState(data)
	fmt.Println(message)
	return
	if !result {
		code = e.ERROR
		c.JSON(http.StatusNotFound, gin.H{
			"message": e.GetMsg(code),
			"code":    code,
			"data":    message,
		})
		return
	}

	//插入前判断数据是否存在，存在的话标记老数据
	tax := str.String(data["taxNo"])
	icode := str.String(data["invoiceCode"])
	number := str.String(data["invoiceNumber"])
	state := message
	key := fmt.Sprintf("%s-%s", icode, number)
	//判断key存在不存在
	if isIn := gredis.Exists(key); isIn {
		//key存在的话比较数据是否有变动
		fmt.Println(1111111111)
		oldValue, error := gredis.Get(key)
		if error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": error,
			})
			return
		}

		//新旧值如果一致，直接返回，不做任何操作
		if oldValue == state {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "数据重复",
			})
			return
		} else {
			fmt.Println(2222222222)
			//新旧值不一致,插入新值，更改旧值状态
			maps := make(map[string]interface{})
			maps["state"] = 1
			models.EditReceiveInvoice(tax, icode, number, maps)
			//插入新数据
			num := models.AddReceiveInvoice(data)
			//更改缓存值
			gredis.Set(key, state, 20)
			c.JSON(http.StatusOK, gin.H{
				"data":    num,
				"code":    code,
				"message": e.GetMsg(code),
			})
		}

	} else {
		//key不存在的话，获取数据库是否有相同值，
		res := models.GetReceiveInvoiceAsMap(tax, icode, number)
		if res != nil {
			fmt.Println(33333333333)
			//转map 取state
			clientJson := res["client_details_json"].(string)
			clientDetail := utils.Decode(clientJson).(map[string]interface{})
			msg := clientDetail["Data"].(map[string]interface{})
			if msg["State"] == state {
				//如果值相等  记录缓存
				gredis.Set(key, state, 20)
				c.JSON(http.StatusNotFound, gin.H{
					"message": "数据重复",
				})
				return
			}
			//有更改更新数据
			maps := make(map[string]interface{})
			maps["state"] = 1
			models.EditReceiveInvoice(tax, icode, number, maps)
			//插入新数据
			num := models.AddReceiveInvoice(data)
			gredis.Set(key, state, 20)
			c.JSON(http.StatusOK, gin.H{
				"data":    num,
				"code":    code,
				"message": e.GetMsg(code),
			})
		} else {
			fmt.Println(44444444)
			//如果没有直接插入，并记录缓存
			num := models.AddReceiveInvoice(data)
			gredis.Set(key, state, 20)
			c.JSON(http.StatusOK, gin.H{
				"data":    num,
				"code":    code,
				"message": e.GetMsg(code),
			})
		}
	}
}

func getState(data map[string]interface{}) (string, bool) {
	v := data["clientDetailsJson"]
	fmt.Println(v)
	return "", false
	detail, ok := v.(map[string]interface{})

	if !ok {
		return "", false
	}
	vv, exists := detail["Data"]
	if !exists {
		return "参数缺少Data字段", false
	}

	Data, ok := vv.(map[string]interface{})
	if !ok {
		return "Data参数格式解析错误", false
	}
	state := Data["State"].(string)
	return state, true
}

func Bill(c *gin.Context) {
	data := make(map[string]interface{})
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err)
	}
	v := data["clientDetailsJson"]
	detail, ok := v.(map[string]interface{})
	if !ok {
		logging.Info(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "参数格式解析错误",
		})
		return
	}
	vv, exists := detail["Data"]
	if !exists {
		logging.Info(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "clientDetailsJson->Data不存在",
		})
		return
	}

	Data, ok := vv.(map[string]interface{})
	if !ok {
		logging.Info(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "参数格式解析错误",
		})
		return
	}
	state := Data["State"].(string)
	iCode := data["invoiceCode"].(string)
	iNum := data["invoiceNumber"].(string)
	key := fmt.Sprintf("%s-%s", iCode, iNum)

	//判断key存在不存在
	if isIn := gredis.Exists(key); isIn {
		//key存在的话比较数据是否有变动
		oldValue, error := gredis.Get(key)
		if error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": error,
			})
			return
		}

		//新旧值如果一致，直接返回，不做任何操作
		if oldValue == state {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "数据重复",
			})
			return
		} else {
			//新旧值不一致,插入新值，更改旧值状态
			fmt.Println(22222)
		}

	} else {
		fmt.Println(11111)
		//key不存在的话，获取数据库是否有相同值，如果没有直接插入，如果有比较state是否一致，如果一致直接返回并写入缓存，如果不一致插入新值，并更改旧值状态
		gredis.Set(key, state, 20)
	}

}

func AddBill(c *gin.Context) {
	MCache := mcache.StartInstance()
	key := "custom_key1"
	user := "lilei"
	if pointer, ok := MCache.GetPointer(key); ok {
		fmt.Println("获取缓存success", pointer)
	} else {
		err := MCache.SetPointer(key, user, time.Second*10)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("设置缓存success")
	}
}

type Xdata struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func testStruct(xcdata *[]Xdata) {
	fmt.Println(*xcdata)
}

type InvoiceInfo struct {
	Id                int
	Billing_date      string
	CheckCode6        string
	InvoiceCode       string
	InvoiceNumber     string
	TaxNo             string
	TotalAmount       string
	ClientDetailsJson string
	SendDatetime      string
	OpenType          int
	OpenDatetime      string
	OpenDetailsJson   string
}

func Test(c *gin.Context) {
	data := make(map[string]interface{})
	err := c.BindJSON(&data)
	code := e.SUCCESS
	if err != nil {
		code = e.ERROR
		logging.Info(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": e.GetMsg(code),
			"code":    code,
			"data":    err,
		})
		return
	}

	tax := str.String(data["taxNo"])
	offset := int(data["offset"].(float64))
	list := models.GetReceiveList(tax, offset)
	// newList := make(map[int]map[string]interface{})
	// new := make(map[string]interface{})
	// for key, val := range list {
	// new["id"] = val["id"]
	// new["clientDetailsJson"] = val["client_details_json"]
	// new["sendDatetime"] = val["send_datetime"]
	// new["openDetailsJson"] = val["open_details_json"]
	// new["openDatetime"] = val["open_datetime"]
	// new["openType"] = val["open_type"]
	// new["taxNo"] = val["tax_no"]
	// new["invoiceCode"] = val["invoice_code"]
	// new["invoiceNumber"] = val["invoice_number"]
	// new["billingDate"] = val["billing_date"]
	// new["checkCode_6"] = val["check_code_6"]
	// new["totalAmount"] = val["total_amount"]
	// newList[key] = new
	// new = make(map[string]interface{})
	// }

	c.JSON(http.StatusNotFound, gin.H{
		"message": "success",
		"code":    0,
		"data":    list,
	})
}
