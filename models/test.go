package models

import (
	"work/pkg/logging"

	"github.com/kirinlabs/utils/str"
)

type LogTag struct {
	Name     string `json:"name"`
	State    int    `json:"state"`
	CreateBy string `json:"created_by"`
}

func Test() map[string]interface{} {
	//l := &LogTag{}
	//err := db.Table("blog_tag").Where("id", ">", 1).First(l)

	list, err := db.Table("blog_tag").Where("id", ">=", 1).Fetch()
	if err != nil {
		logging.Info(err)
	}
	return list
}

func AddBill(maps map[string]interface{}) int64 {
	num, err := db.Table("blog_tag").Insert(maps)
	if err != nil {
		logging.Info(err)
	}

	return num
}

func AddStruct(maps map[string]interface{}) int64 {
	list := &LogTag{
		str.String(maps["name"]),
		maps["state"].(int),
		str.String(maps["created_by"]),
	}
	num, err := db.Table("blog_tag").Insert(list)
	if err != nil {
		logging.Info(err)
	}

	return num
}

func GetNotDetailBill(where map[string]interface{}, field string, limit int) []map[string]interface{} {
	list, err := db.Table("bill_expending_invoice").Where(where).Limit(limit).Fields(field).FetchAll()
	if err != nil {
		logging.Info(err)
	}
	return list
}

func UpdateInvoiceDetail(where map[string]interface{}, update map[string]interface{}) int64 {
	num, err := db.Table("bill_expending_invoice").Where(where).Update(update)
	if err != nil {
		logging.Info(err)
	}
	return num
}
