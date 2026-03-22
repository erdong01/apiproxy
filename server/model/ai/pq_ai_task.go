// 自动生成模板PqAiTask
package ai

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// pqAiTask表 结构体  PqAiTask
type PqAiTask struct {
	Id                    int64          `gorm:"column:id;primaryKey" json:"id"`                              //type:int64            comment:                                                                          version:2026-02-21 12:58
	CreatedAt             *time.Time     `gorm:"column:created_at" json:"createdAt"`                          //type:*time.Time       comment:创建时间                                                                  version:2026-02-21 12:58
	UpdatedAt             *time.Time     `gorm:"column:updated_at" json:"updatedAt"`                          //type:*time.Time       comment:更新时间                                                                  version:2026-02-21 12:58
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`                          //type:gorm.DeletedAt   comment:                                                                          version:2026-02-21 12:58
	GenerateTaskId        string         `gorm:"column:generate_task_id" json:"generateTaskId"`               //type:string           comment:供应商任务id                                                              version:2026-02-21 12:58
	UserId                int64          `gorm:"column:user_id" json:"userId"`                                //type:int64            comment:用户id                                                                    version:2026-02-21 12:58
	Model                 string         `gorm:"column:model" json:"model"`                                   //type:int64            comment:ai模型id                                                                  version:2026-02-21 12:58
	Status                string         `gorm:"column:status" json:"status"`                                 //type:string           comment:状态                                                                      version:2026-02-21 12:58
	TaskCreatedAt         int64          `gorm:"column:task_created_at" json:"taskCreatedAt"`                 //type:int64            comment:任务创建时间的 Unix 时间戳（秒）。                                        version:2026-02-21 12:58
	TaskUpdatedAt         int64          `gorm:"column:task_updated_at" json:"taskUpdatedAt"`                 //type:int64            comment:任务当前状态更新时间的 Unix 时间戳（秒）。                                version:2026-02-21 12:58
	Content               datatypes.JSON `gorm:"column:content" json:"content"`                               //type:datatypes.JSON   comment:视频生成任务的输出内容。                                                  version:2026-02-21 12:58
	Seed                  int            `gorm:"column:seed" json:"seed"`                                     //type:int              comment:本次请求使用的种子整数值。                                                version:2026-02-21 12:58
	Resolution            string         `gorm:"column:resolution" json:"resolution"`                         //type:string           comment:生成视频的分辨率                                                          version:2026-02-21 12:58
	Ratio                 string         `gorm:"column:ratio" json:"ratio"`                                   //type:string           comment:生成视频的宽高比。                                                        version:2026-02-21 12:58
	Duration              int            `gorm:"column:duration" json:"duration"`                             //type:int              comment:生成视频的时长，单位：秒。                                                version:2026-02-21 12:58
	Frames                int            `gorm:"column:frames" json:"frames"`                                 //type:int              comment:生成视频的帧数。                                                          version:2026-02-21 12:58
	FramesPerSecond       int            `gorm:"column:frames_per_second;" json:"framesPerSecond"`            //type:int              comment:生成视频的帧率。                                                          version:2026-02-22 09:26
	GenerateAudio         bool           `gorm:"column:generate_audio" json:"generateAudio"`                  //type:bool             comment:生成的视频是否包含与画面同步的声音。仅 Seedance 1.5 pro 会返回该参数。    version:2026-02-21 12:58
	Draft                 bool           `gorm:"column:draft" json:"draft"`                                   //type:bool             comment:生成的视频是否为 Draft 视频。仅 Seedance 1.5 pro 会返回该参数。           version:2026-02-21 12:58
	DraftTaskId           string         `gorm:"column:draft_task_id" json:"draftTaskId"`                     //type:string           comment:Draft 视频任务 ID。基于 Draft 视频生成正式视频时，会返回该参数。          version:2026-02-21 12:58
	ServiceTier           string         `gorm:"column:service_tier" json:"serviceTier"`                      //type:string           comment:实际处理任务使用的服务等级。                                              version:2026-02-21 12:58
	ExecutionExpiresAfter int            `gorm:"column:execution_expires_after" json:"executionExpiresAfter"` //type:int              comment:任务超时阈值，单位：秒。                                                  version:2026-02-21 12:58
	Params                datatypes.JSON `gorm:"column:params" json:"params"`                                 //type:datatypes.JSON   comment:存储原始生成参数 (如 seed, motion_bucket_id 等)                           version:2026-02-21 12:58
	ErrorMessage          string         `gorm:"column:error_message" json:"errorMessage"`                    //type:string           comment:失败原因                                                                  version:2026-02-21 12:58
	CompletionTokens      int64          `gorm:"column:completion_tokens" json:"completionTokens"`            //type:int64            comment:输出视频花费的 token 数                                                   version:2026-02-21 12:58
	TotalTokens           int64          `gorm:"column:total_tokens" json:"totalTokens"`                      //type:int64            comment:本次请求消耗的总 token 数                                                 version:2026-02-21 12:58
	Key                   string         `gorm:"column:key" json:"key"`                                       //type:string           comment:供应商key                                                                 version:2026-02-21 12:58
	RequestId             string         `gorm:"column:request_id" json:"requestId"`                          //type:string           comment:                                                                          version:2026-02-21 12:58
	VendorAmount          float64        `gorm:"column:vendor_amount;" json:"vendorAmount"`                   //type:float64          comment:成本                                                                      version:2026-02-22 09:26
	RetailAmount          float64        `gorm:"column:retail_amount;" json:"retailAmount"`                   //type:float64          comment:零售                                                                      version:2026-02-22 09:26

}

// TableName pqAiTask表 PqAiTask自定义表名 pq_ai_task
func (PqAiTask) TableName() string {
	return "pq_ai_task"
}
