package models

import (
	"time"

	"work/pkg/logging"

	"github.com/kirinlabs/utils/str"
)

type InvoiceAllInfo struct {
	Billing_date      string `json:"billing_date"`
	CheckCode6        string `json:"check_code_6"`
	InvoiceCode       string `json:"invoice_code"`
	InvoiceNumber     string `json:"invoice_number"`
	TaxNo             string `json:"tax_no"`
	TotalAmount       string `json:"total_amount"`
	ClientDetailsJson string `json:"client_details_json"`
	SendDatetime      string `json:"send_datetime"`
	OpenType          int    `json:"open_type"`
}

type InvoiceInfo struct {
	Id                int    `json:"id"`
	Billing_date      string `json:"billing_date"`
	CheckCode6        string `json:"check_code_6"`
	InvoiceCode       string `json:"invoice_code"`
	InvoiceNumber     string `json:"invoice_number"`
	TaxNo             string `json:"tax_no"`
	TotalAmount       string `json:"total_amount"`
	ClientDetailsJson string `json:"client_details_json"`
	SendDatetime      string `json:"send_datetime"`
	OpenType          int    `json:"open_type"`
	OpenDatetime      string `json:"open_datetime"`
	OpenDetailsJson   string `json:"open_details_json"`
}

func GetReceiveInvoice(taxNo string, invoiceCode string, invoice_number string) *InvoiceAllInfo {
	i := &InvoiceAllInfo{}
	err := db.Table("invoice_all_info").Where("tax_no", taxNo).Where("invoice_code", invoiceCode).Where("invoice_number", invoice_number).First(i)
	if err != nil {
		logging.Info(err)
	}
	return i
}

func GetReceiveInvoiceAsMap(taxNo string, invoiceCode string, invoice_number string) map[string]interface{} {
	list, err := db.Table("invoice_all_info").Where("tax_no", taxNo).Where("state", 0).Where("invoice_code", invoiceCode).Where("invoice_number", invoice_number).Fetch()
	if err != nil {
		logging.Info(err)
	}
	return list
}

func AddReceiveInvoice(maps map[string]interface{}) int64 {
	list := &InvoiceAllInfo{
		str.String(maps["billingDate"]),
		str.String(maps["checkCode_6"]),
		str.String(maps["invoiceCode"]),
		str.String(maps["invoiceNumber"]),
		str.String(maps["taxNo"]),
		str.String(maps["totalAmount"]),
		str.String(maps["clientDetailsJson"]),
		time.Now().Format("2006-01-02 15:04:05"),
		3,
	}
	num, err := db.Table("invoice_all_info").Insert(list)

	if err != nil {
		logging.Info(err)
	}
	return num

}

func EditReceiveInvoice(taxNo string, invoiceCode string, invoice_number string, u map[string]interface{}) int64 {
	num, err := db.Table("invoice_all_info").Where("tax_no", taxNo).Where("invoice_code", invoiceCode).Where("invoice_number", invoice_number).Update(u)
	if err != nil {
		logging.Info(err)
	}
	return num
}

func GetReceiveList(tax_no string, offset int) (invoiceList []*InvoiceInfo) {
	err := db.Table("invoice_all_info").Where("tax_no", tax_no).Limit(offset, 500).Find(&invoiceList)
	if err != nil {
		logging.Info(err)
	}
	return
	// list, err := db.Table("invoice_all_info").Where("tax_no", tax_no).Where("state", 0).Limit(offset, 500).FetchAll()
	// logging.Info(list)
	// if err != nil {
	// 	logging.Info(err)
	// }
	// return list
}
