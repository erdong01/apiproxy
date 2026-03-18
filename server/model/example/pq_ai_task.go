
// 自动生成模板PqAiTask
package example
import (
	"time"
	"gorm.io/datatypes"
)

// pqAiTask表 结构体  PqAiTask
type PqAiTask struct {
  Id  *int64 `json:"id" form:"id" gorm:"primarykey;column:id;"`  //id字段
  CreatedAt  *time.Time `json:"createdAt" form:"createdAt" gorm:"comment:创建时间;column:created_at;"`  //创建时间
  UpdatedAt  *time.Time `json:"updatedAt" form:"updatedAt" gorm:"comment:更新时间;column:updated_at;"`  //更新时间
  DeletedAt  *time.Time `json:"deletedAt" form:"deletedAt" gorm:"column:deleted_at;"`  //deletedAt字段
  GenerateTaskId  *string `json:"generateTaskId" form:"generateTaskId" gorm:"comment:供应商任务id;column:generate_task_id;size:64;"`  //供应商任务id
  UserId  *int64 `json:"userId" form:"userId" gorm:"comment:用户id;column:user_id;"`  //用户id
  ModelId  *int64 `json:"modelId" form:"modelId" gorm:"comment:ai模型id;column:model_id;"`  //ai模型id
  Status  *string `json:"status" form:"status" gorm:"comment:状态;column:status;size:64;"`  //状态
  Params  datatypes.JSON `json:"params" form:"params" gorm:"comment:存储原始生成参数 (如 seed, motion_bucket_id 等);column:params;" swaggertype:"object"`  //存储原始生成参数 (如 seed, motion_bucket_id 等)
  ErrorMessage  *string `json:"errorMessage" form:"errorMessage" gorm:"comment:失败原因;column:error_message;"`  //失败原因
  CompletionTokens  *int32 `json:"completionTokens" form:"completionTokens" gorm:"comment:输出视频花费的 token 数;column:completion_tokens;"`  //输出视频花费的 token 数
  TotalTokens  *int32 `json:"totalTokens" form:"totalTokens" gorm:"comment:本次请求消耗的总 token 数;column:total_tokens;"`  //本次请求消耗的总 token 数
  Key  *string `json:"key" form:"key" gorm:"column:key;size:128;"`  //key字段
}


// TableName pqAiTask表 PqAiTask自定义表名 pq_ai_task
func (PqAiTask) TableName() string {
    return "pq_ai_task"
}





