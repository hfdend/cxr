package models

import (
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/jinzhu/gorm"
)

type OrderAddress struct {
	Model
	OrderID       string `json:"order_id"`
	UserID        int    `json:"user_id"`
	AddressID     int    `json:"address_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	DistrictCode  string `json:"district_code"`
	DistrictName  string `json:"district_name"`
	DetailAddress string `json:"detail_address"`
	OptUserID     int    `json:"opt_user_id"`
	Created       int64  `json:"created"`
}

var OrderAddressDefault OrderAddress

func (OrderAddress) TableName() string {
	return "order_address"
}

func (addr *OrderAddress) Insert(db *gorm.DB) error {
	addr.Created = time.Now().Unix()
	return db.Create(addr).Error
}

func (OrderAddress) GetByOrderID(orderID string) (*OrderAddress, error) {
	var data OrderAddress
	if err := cli.DB.Where("order_id = ?", orderID).Last(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (OrderAddress) GetByID(id int) (*OrderAddress, error) {
	var data OrderAddress
	if err := cli.DB.Where("id = ?", id).Find(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (OrderAddress) GetByOrderIDs(orderIDs []string) (list []*OrderAddress, err error) {
	err = cli.DB.Where("order_id in (?)", orderIDs).Find(&list).Error
	return
}
