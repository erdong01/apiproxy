
// 自动生成模板PqAiModel
package ai
import (
	"time"
)

// pqAiModel表 结构体  PqAiModel
type PqAiModel struct {
  Id  *int64 `json:"id" form:"id" gorm:"primarykey;column:id;"`  //id字段
  CreatedAt  *time.Time `json:"createdAt" form:"createdAt" gorm:"comment:创建时间;column:created_at;"`  //创建时间
  UpdatedAt  *time.Time `json:"updatedAt" form:"updatedAt" gorm:"comment:更新时间;column:updated_at;"`  //更新时间
  DeletedAt  *time.Time `json:"deletedAt" form:"deletedAt" gorm:"column:deleted_at;"`  //deletedAt字段
  Name  *string `json:"name" form:"name" gorm:"comment:名称;column:name;size:255;"`  //名称
  Provider  *string `json:"provider" form:"provider" gorm:"comment:供应商;column:provider;size:64;"`  //供应商
  Version  *string `json:"version" form:"version" gorm:"comment:版本;column:version;size:64;"`  //版本
}


// TableName pqAiModel表 PqAiModel自定义表名 pq_ai_model
func (PqAiModel) TableName() string {
    return "pq_ai_model"
}





