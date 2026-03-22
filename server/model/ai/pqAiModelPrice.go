package ai

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// PqAiModelPrice  模型价格。
type PqAiModelPrice struct {
	Id              int64          `gorm:"column:id;primaryKey;" json:"Id"`                 //主键
	CreatedAt       *time.Time     `gorm:"column:created_at" json:"CreatedAt"`              //
	UpdatedAt       *time.Time     `gorm:"column:updated_at" json:"UpdatedAt"`              //
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"DeletedAt"`              //
	PqAiModelId     int64          `gorm:"column:pq_ai_model_id" json:"PqAiModelId"`        // 模型id
	Resolution      string         `gorm:"column:resolution" json:"Resolution"`             // 分辨率 1080P 720P 480P
	VendorPrice     float64        `gorm:"column:vendor_price" json:"VendorPrice"`          // 供应商价格
	VendorUnit      string         `gorm:"column:vendor_unit" json:"VendorUnit"`            // 供应商单位 1 千tokens  2 张
	RetailPrice     float64        `gorm:"column:retail_price" json:"RetailPrice"`          // 零售价格
	RetailUnit      string         `gorm:"column:retail_unit" json:"RetailUnit"`            // 零售单位 1 单秒  2 单次
	GenerationModes string         `gorm:"column:generation_modes;" json:"GenerationModes"` // 生成模式                       version:2026-02-21 15:44

}

// TableName 表名:pq_ai_model_price，模型价格。
// 说明:
func (*PqAiModelPrice) TableName() string {
	return "pq_ai_model_price"
}

func (that *PqAiModelPrice) Calculate(v int64, duration int64) (vendorAmount float64, retailAmount float64) {
	//成本
	switch that.VendorUnit {
	case "1":
		dTokens := decimal.NewFromInt(v)
		dThousand := decimal.NewFromInt(1000)
		dPrice := decimal.NewFromFloat(that.VendorPrice)
		dVendor := dTokens.Div(dThousand).Mul(dPrice).Truncate(2)
		vendorAmount, _ = dVendor.Float64()
	case "2":

	}

	//零售
	switch that.RetailUnit {
	case "1":
		retailAmount, _ = decimal.NewFromInt(duration).Mul(decimal.NewFromFloat(that.RetailPrice)).Truncate(2).Float64()
	case "2":
		retailAmount = that.RetailPrice
	}

	return
}
