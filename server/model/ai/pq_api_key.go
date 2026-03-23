// 自动生成模板PqApiKey
package ai

import (
	"time"
)

// pqApiKey表 结构体  PqApiKey
type PqApiKey struct {
	Id          *int64     `json:"id" form:"id" gorm:"primarykey;column:id;"`                                       //id字段
	UpdatedAt   *time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;"`                            //updatedAt字段
	DeletedAt   *time.Time `json:"deletedAt" form:"deletedAt" gorm:"column:deleted_at;"`                            //deletedAt字段
	CreatedAt   *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;"`                            //createdAt字段
	UserId      *int64     `json:"userId" form:"userId" gorm:"comment:用户id;column:user_id;"`                        //用户id
	AiModelId   *int64     `json:"aiModelId" form:"aiModelId" gorm:"comment:ai模型;column:ai_model_id;"`              //ai模型
	Key         *string    `json:"key" form:"key" gorm:"comment:密钥;column:key;size:255;"`                           //密钥
	TotalTokens *int64     `json:"totalTokens" form:"totalTokens" gorm:"comment:拥有tokens数;column:total_tokens;"`    //拥有tokens数
	UseTokens   *string    `json:"useTokens" form:"useTokens" gorm:"comment:已消耗tokens;column:use_tokens;size:255;"` //已消耗tokens
	UserKey     *string    `gorm:"column:user_key" json:"UserKey"`                                                  //  外部用户key     version:2026-02-18 17:40
	UserName    string     `gorm:"column:user_name" json:"UserName"`                                                //                  version:2026-02-18 21:21
	Status      int        `gorm:"column:status;" json:"status" form:"status"`                                      //  状态 1 启用 2 禁用    version:2026-02-22 10:29
	Rate        int        `gorm:"column:rate;" json:"rate"`                                                        //  速率                  version:2026-02-23 11:49
}

// TableName pqApiKey表 PqApiKey自定义表名 pq_api_key
func (PqApiKey) TableName() string {
	return "pq_api_key"
}
