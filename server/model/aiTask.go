package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AiTask  AI 任务表。

type AiTask struct {
	Id               *int64          `gorm:"column:id;primaryKey" json:"Id"`                   //           comment:                                                   version:2026-02-15 21:40
	CreatedAt        *time.Time      `gorm:"column:created_at" json:"CreatedAt"`               //创建时间                                           version:2026-02-15 21:40
	UpdatedAt        *time.Time      `gorm:"column:updated_at" json:"UpdatedAt"`               //更新时间                                           version:2026-02-15 21:40
	DeletedAt        *gorm.DeletedAt `gorm:"column:deleted_at" json:"DeletedAt"`               //DeletedAt   comment:                                                   version:2026-02-15 21:40
	GenerateTaskId   string          `gorm:"column:generate_task_id" json:"GenerateTaskId"`    //供应商任务id                                       version:2026-02-15 21:40
	UserId           *int64          `gorm:"column:user_id" json:"UserId"`                     //用户id                                             version:2026-02-15 21:40
	ModelId          *int64          `gorm:"column:model_id" json:"ModelId"`                   //ai模型id                                           version:2026-02-15 21:40
	Status           string          `gorm:"column:status" json:"Status"`                      //状态                                               version:2026-02-15 21:40
	Params           datatypes.JSON  `gorm:"column:params" json:"Params"`                      //存储原始生成参数 (如 seed, motion_bucket_id 等)    version:2026-02-15 21:40
	ErrorMessage     string          `gorm:"column:error_message" json:"ErrorMessage"`         //失败原因                                           version:2026-02-15 21:40
	CompletionTokens *int            `gorm:"column:completion_tokens" json:"CompletionTokens"` //输出视频花费的 token 数                            version:2026-02-15 21:40
	TotalTokens      *int            `gorm:"column:total_tokens" json:"TotalTokens"`           //本次请求消耗的总 token 数                          version:2026-02-15 21:40
	Key              string          `gorm:"column:key" json:"Key"`                            //                                                   version:2026-02-16 20:08

}

// TableName 表名:ai_task，AI 任务表。
// 说明:
func (*AiTask) TableName() string {
	return "pq_ai_task"
}
