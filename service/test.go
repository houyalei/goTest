package service

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"work/models"
	"work/pkg/logging"

	"github.com/kirinlabs/HttpRequest"
	"github.com/kirinlabs/utils"
	"github.com/kirinlabs/utils/str"
)

const limit = 1

func Combine(inCh1, inCh2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case x, open := <-inCh1:
				if !open {
					logging.Info(open)
					inCh1 = nil
					continue
				}
				logging.Info(x)
				out <- x
			case x, open := <-inCh2:
				if !open {
					logging.Info(open)
					inCh2 = nil
					continue
				}
				logging.Info(x)
				out <- x
			}
			if inCh1 == nil && inCh2 == nil {
				break
			}
		}
	}()

	return out
}

func Test() {
	var wg sync.WaitGroup
	wg.Add(2)
	ticker1 := time.NewTicker(2 * time.Second)
	go func(t *time.Ticker) {
		defer wg.Done()
		for {
			<-t.C
			fmt.Println("get ticker1", time.Now().Format("2006-01-02 15:04:05"))
			logging.Info(time.Now().Format("2006-01-02 15:04:05"))
		}
	}(ticker1)

	wg.Wait()
}

func GetBillDetailInfo() interface{} {
	where := make(map[string]interface{})
	where["is_pull"] = 0
	where["details_json"] = ""
	field := "invoice_code,invoice_number,customer_id,total_amount,billing_date"
	data := models.GetNotDetailBill(where, field, 30)
	if len(data) == 0 {
		logging.Info("没有要获取全票信息的票据")
		return make(map[string]map[string]string, 0)
	}
	sendData := make(map[string]map[string]string)

	for _, v := range data {
		//连接字符串
		var key bytes.Buffer
		key.WriteString(str.String(v["invoice_code"]))
		key.WriteString("-")
		key.WriteString(str.String(v["invoice_number"]))
		key.WriteString("-")
		key.WriteString(str.String(v["customer_id"]))

		//封装发送数据的条件
		arr := make(map[string]string)
		arr["InvoiceCode"] = str.String(v["invoice_code"])
		arr["InvoiceNumber"] = str.String(v["invoice_number"])
		arr["TotalAmount"] = str.String(v["total_amount"])
		arr["BillingDate"] = str.String(v["billing_date"])
		sendData[key.String()] = arr
	}

	url := "http://www.test.yundoubpo.com/api/invoice/verify"
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	res, err := req.Post(url, map[string]interface{}{
		"billData": utils.Json(sendData),
	})
	body, err := res.Body()
	if err != nil {
		logging.Info(err)
	}
	resData := utils.Decode(string(body)).(map[string]interface{})
	for k, v := range resData {
		bill, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("导入数据格式错误")
		}
		if bill["Code"] != "0" {
			logging.Info("返回数据错误")
			return nil
		}
		UpdateDetailInvoice(k, bill["Data"])
	}

	return nil
}

func UpdateDetailInvoice(str1 string, data interface{}) int64 {
	where := make(map[string]interface{})
	update := make(map[string]interface{})
	var slice []string
	//分割key
	slice = strings.Split(str1, "-")
	where["invoice_code"] = str.String(slice[0])
	where["invoice_number"] = str.String(slice[1])
	where["customer_id"] = str.String(slice[2])

	update["details_json"] = data
	update["is_pull"] = 1

	res := models.UpdateInvoiceDetail(where, update)
	return res
}
