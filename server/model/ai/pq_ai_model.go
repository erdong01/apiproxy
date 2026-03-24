// 自动生成模板PqAiModel
package ai

import (
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/shopspring/decimal"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

// pqAiModel表 结构体  PqAiModel
type PqAiModel struct {
	Id             *int64           `json:"id" form:"id" gorm:"column:id;primarykey;"`                            //id字段
	CreatedAt      *time.Time       `json:"createdAt" form:"createdAt" gorm:"comment:创建时间;column:created_at;"`    //创建时间
	UpdatedAt      *time.Time       `json:"updatedAt" form:"updatedAt" gorm:"comment:更新时间;column:updated_at;"`    //更新时间
	DeletedAt      *time.Time       `json:"deletedAt" form:"deletedAt" gorm:"column:deleted_at;"`                 //deletedAt字段
	Name           string           `json:"name" form:"name" gorm:"comment:名称;column:name;size:255;"`             //名称
	Provider       *string          `json:"provider" form:"provider" gorm:"comment:供应商;column:provider;size:64;"` //供应商
	Version        *string          `json:"version" form:"version" gorm:"comment:版本;column:version;size:64;"`     //版本
	PqAiModelPrice []PqAiModelPrice `json:"PqAiModelPrice" gorm:"foreignKey:PqAiModelId;references:Id"`
}

// TableName pqAiModel表 PqAiModel自定义表名 pq_ai_model
func (PqAiModel) TableName() string {
	return "pq_ai_model"
}
func (that *PqAiModel) ModelId() string {
	return that.Name + *that.Version
}
func SplitModelId(modelId string) (name, version string) {
	idx := strings.LastIndex(modelId, "-")
	if idx <= 0 || idx == len(modelId)-1 {
		return modelId, ""
	}

	return modelId[:idx], modelId[idx:]

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

// 计算火山费用
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

func PanQuModelPriceCalculate(model string, resolution string, draft_task_id string, duration int64, draftVideoDuration int64) (price float64) {
	if draft_task_id != "" {
		draftModelPrice := resolution + "_draft_task"
		var totalPrice float64
		var totalPrice2 decimal.Decimal

		if modelPrice, ok := PanQuModelPriceMap[model][resolution]; ok {
			dDraftVideoDuration := decimal.NewFromInt(duration)
			dPrice := decimal.NewFromFloat(modelPrice.Price)

			totalPrice2 = dDraftVideoDuration.Mul(dPrice).Truncate(2)
		}

		if modelPrice, ok := PanQuModelPriceMap[model][draftModelPrice]; ok {
			dDraftVideoDuration := decimal.NewFromInt(draftVideoDuration)
			dPrice := decimal.NewFromFloat(modelPrice.Price)

			totalPrice, _ = dDraftVideoDuration.Mul(dPrice).Add(totalPrice2).Truncate(2).Float64()
		}

		return totalPrice
	} else {
		if modelPrice, ok := PanQuModelPriceMap[model][resolution]; ok {
			duration := decimal.NewFromInt(duration)
			dPrice := decimal.NewFromFloat(modelPrice.Price)

			totalPrice, _ := duration.Mul(dPrice).Truncate(2).Float64()
			return totalPrice
		}
	}

	return
}

// PriceIndex 将模型价格列表按 分辨率 -> 生成模式 建立索引，便于业务侧快速查询。
func (that *PqAiModel) PriceIndex() map[string]map[string]PqAiModelPrice {
	res := make(map[string]map[string]PqAiModelPrice)
	for _, price := range that.PqAiModelPrice {
		if _, ok := res[price.Resolution]; !ok {
			res[price.Resolution] = make(map[string]PqAiModelPrice)
		}
		res[price.Resolution][price.GenerationModes] = price
	}

	return res
}

// ResolveGenerationModes 根据任务返回的特征字段推导候选生成模式，按优先级返回。
func ResolveGenerationModes(generateAudio *bool, draft *bool, draftTaskID *string) []string {
	modes := make([]string, 0, 5)
	if generateAudio != nil {
		if *generateAudio {
			modes = append(modes, "generate_audio_true")
		} else {
			modes = append(modes, "generate_audio_false")
		}
	}
	if draftTaskID != nil && *draftTaskID != "" {
		modes = append(modes, "draft_task")
	}
	if draft != nil {
		if *draft {
			modes = append(modes, "draft_true")
		} else {
			modes = append(modes, "draft_false")
		}
	}

	return append(modes, "", "default")
}

// FindModelPrice 根据模型、分辨率和任务特征查找对应价格，兼容默认模式兜底。
func FindModelPrice(
	modelPrices map[string]map[string]map[string]PqAiModelPrice,
	resp *model.GetContentGenerationTaskResponse,
) (PqAiModelPrice, bool) {
	resolutionPrices, ok := modelPrices[resp.Model][*resp.Resolution]
	if !ok {
		return PqAiModelPrice{}, false
	}

	for _, mode := range ResolveGenerationModes(resp.GenerateAudio, resp.Draft, resp.DraftTaskID) {
		if price, ok := resolutionPrices[mode]; ok {
			return price, true
		}
	}

	if len(resolutionPrices) == 1 {
		for _, price := range resolutionPrices {
			return price, true
		}
	}

	return PqAiModelPrice{}, false
}

// ModelPriceInit 从数据库加载模型价格，并构建 模型 -> 分辨率 -> 生成模式 的三级索引。
func ModelPriceInit() (res map[string]map[string]map[string]PqAiModelPrice) {
	var pqAiModel []PqAiModel
	global.GVA_DB.Model(&PqAiModel{}).
		Preload("PqAiModelPrice").Find(&pqAiModel)
	res = make(map[string]map[string]map[string]PqAiModelPrice)
	for i := range pqAiModel {
		res[pqAiModel[i].Name] = pqAiModel[i].PriceIndex()
	}

	return
}
