package models

import (
	"strings"
	"time"

	"fmt"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
	"github.com/jinzhu/gorm"
)

type Product struct {
	Model
	Name string `json:"name"`
	// 类型
	Type string `json:"type"`
	// 味道
	Taste string `json:"taste"`
	// 商品规格
	Unit string `json:"unit"`
	// 售价
	Price float64 `json:"price"`
	// 图片
	Image string `json:"image"`
	// 备注
	Mark string `json:"mark"`
	// 介绍
	Intro   string `json:"intro"`
	IsDel   Sure   `json:"is_del"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`

	ImageSrc string `json:"image_src" gorm:"-"`
}

var ProductDefault Product

func (Product) TableName() string {
	return "product"
}

func (p *Product) Save() error {
	if p.Created == 0 {
		p.Created = time.Now().Unix()
	}
	p.Updated = time.Now().Unix()
	return cli.DB.Save(p).Error
}

func (Product) GetByID(id int) (*Product, error) {
	var data Product
	if err := cli.DB.Where("id = ? and is_del = ?", id, SureNo).Find(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	data.SetImageSrc()
	return &data, nil
}

func (Product) DelByID(id int) error {
	data := map[string]interface{}{
		"is_del":  SureYes,
		"updated": time.Now().Unix(),
	}
	return cli.DB.Where("id = ?", id).Update(data).Error
}

func (Product) GetList(pager *Pager) (list []*Product, err error) {
	db := cli.DB.Model(Product{}).Where("is_del = ?", SureNo)
	if pager != nil {
		if db, err = pager.Exec(db); err != nil {
			return
		} else if pager.Count == 0 {
			return
		}
	}
	err = db.Find(&list).Error
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (p *Product) SetImageSrc() {
	if p.Image == "" {
		return
	}
	c := conf.Config.Aliyun.OSS
	p.ImageSrc = fmt.Sprintf("%s/%s", strings.TrimRight(c.Domain, "/"), strings.TrimLeft(p.Image, "/"))
	return
}
