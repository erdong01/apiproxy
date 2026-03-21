// 自动生成模板PqAiModel
package ai

import (
	"time"

	"github.com/shopspring/decimal"
)

// pqAiModel表 结构体  PqAiModel
type PqAiModel struct {
	Id        *int64     `json:"id" form:"id" gorm:"primarykey;column:id;"`                            //id字段
	CreatedAt *time.Time `json:"createdAt" form:"createdAt" gorm:"comment:创建时间;column:created_at;"`    //创建时间
	UpdatedAt *time.Time `json:"updatedAt" form:"updatedAt" gorm:"comment:更新时间;column:updated_at;"`    //更新时间
	DeletedAt *time.Time `json:"deletedAt" form:"deletedAt" gorm:"column:deleted_at;"`                 //deletedAt字段
	Name      *string    `json:"name" form:"name" gorm:"comment:名称;column:name;size:255;"`             //名称
	Provider  *string    `json:"provider" form:"provider" gorm:"comment:供应商;column:provider;size:64;"` //供应商
	Version   *string    `json:"version" form:"version" gorm:"comment:版本;column:version;size:64;"`     //版本
}

// TableName pqAiModel表 PqAiModel自定义表名 pq_ai_model
func (PqAiModel) TableName() string {
	return "pq_ai_model"
}

type ModelPrice struct {
	Price float64
}

var ModelPriceMap = map[string]map[string]ModelPrice{
	"doubao-seedance-2-0-260128": map[string]ModelPrice{
		"text": {
			Price: 0.046, //文本 千tokens  价格
		},
		"image_url": {
			Price: 0.046, //图文文本 千tokens  价格
		},
		"draft_task": {
			Price: 0.028, //参考视频 千tokens  价格
		},
	},
	"doubao-seedance-2-0-fast-260128": map[string]ModelPrice{
		"text": {
			Price: 0.037, //文本 千tokens  价格
		},
		"image_url": {
			Price: 0.037, //图文文本 千tokens  价格
		},
		"draft_task": { // 有参考视频
			Price: 0.022, //参考视频 千tokens  价格
		},
	},
}

// 计算每秒平均成本
func Calculate(model string, contentType string, tokens int64) (price float64) {
	if modelPrice, ok := ModelPriceMap[model][contentType]; ok {
		dTokens := decimal.NewFromInt(tokens)
		dThousand := decimal.NewFromInt(1000)
		dPrice := decimal.NewFromFloat(modelPrice.Price)

		totalPrice, _ := dTokens.Div(dThousand).Mul(dPrice).Truncate(2).Float64()
		return totalPrice
	}

	return
}

type PanQuModelPrice struct {
	Price float64
}

var PanQuModelPriceMap = map[string]map[string]PanQuModelPrice{
	"doubao-seedance-2-0-260128": map[string]PanQuModelPrice{
		"720p": {
			Price: 0.994,
		},
		"720p_draft_task": {
			Price: 0.5,
		},
		"480p": {
			Price: 0.462,
		},
		"480p_draft_task": {
			Price: 0.22,
		},
	},
	"doubao-seedance-2-0-fast-260128": map[string]PanQuModelPrice{
		"720p": {
			Price: 0.8,
		},
		"720p_draft_task": {
			Price: 0.37,
		},
		"480p": {
			Price: 0.372,
		},
		"480p_draft_task": {
			Price: 0.17,
		},
	},
}

func PanQuModelPriceCalculate(model string, resolution string, draft_task_id string, duration int64) (price float64) {
	if draft_task_id != "" {
		model = model + "_draft_task"
	}
	if modelPrice, ok := PanQuModelPriceMap[model][resolution]; ok {
		duration := decimal.NewFromInt(duration)
		dPrice := decimal.NewFromFloat(modelPrice.Price)

		totalPrice, _ := duration.Mul(dPrice).Truncate(2).Float64()
		return totalPrice
	}

	return
}
