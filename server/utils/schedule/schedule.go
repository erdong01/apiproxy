package schedule

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/erdong01/kit"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	arkruntimeModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

var Schedule *kit.Schedule

func Init(ctx context.Context) {
	Schedule = kit.NewSchedule()
	Schedule.Run(ctx)

	var seedanceTask SeedanceTask
	Schedule.Add(&seedanceTask, time.Minute*10, true)
}

type SeedanceTask struct {
}

func (that *SeedanceTask) OnTimer() {
	var limit int = 10
	var lastID int64 = 0
	modelPrice := ai.ModelPriceInit()

	for {
		var aiTaskData []model.AiTask
		// 使用基于 ID 的游标分页，避免数据变更导致的偏移问题
		result := global.GVA_DB.Model(&model.AiTask{}).
			Where("status IS NULL AND id > ?", lastID).
			Order("id ASC").
			Limit(limit).
			Find(&aiTaskData)
		if result.Error != nil {
			fmt.Printf("查询 AI 任务失败: %v\n", result.Error)
			break
		}
		// 没有更多数据，退出循环
		if len(aiTaskData) == 0 {
			break
		}
		for i := range aiTaskData {
			// 记录当前处理的最大 ID，用于下一页游标
			if aiTaskData[i].Id != 0 && aiTaskData[i].Id > lastID {
				lastID = aiTaskData[i].Id
			}

			client := arkruntime.NewClientWithApiKey(aiTaskData[i].Key)
			ctx := context.Background()
			req := arkruntimeModel.GetContentGenerationTaskRequest{
				ID: aiTaskData[i].GenerateTaskId,
			}
			resp, err := client.GetContentGenerationTask(ctx, req)
			if err != nil {
				// fmt.Printf("获取内容生成任务失败: %v\n", err)
				continue
			}
			if resp.Status != arkruntimeModel.StatusQueued && resp.Status != arkruntimeModel.StatusRunning {
				//模型id
				aiTaskData[i].Model = resp.Model
				aiTaskData[i].Status = resp.Status
				aiTaskData[i].TaskCreatedAt = resp.CreatedAt
				aiTaskData[i].TaskUpdatedAt = resp.UpdatedAt
				aiTaskData[i].CompletionTokens = int64(resp.Usage.CompletionTokens)
				aiTaskData[i].TotalTokens = int64(resp.Usage.TotalTokens)
				if content, err := json.Marshal(resp.Content); err == nil {
					aiTaskData[i].Content = content
				}
				aiTaskData[i].Seed = resp.Seed
				aiTaskData[i].Resolution = resp.Resolution
				aiTaskData[i].Ratio = resp.Ratio
				aiTaskData[i].Duration = resp.Duration
				aiTaskData[i].Frames = resp.Frames
				aiTaskData[i].FramesPerSecond = resp.FramesPerSecond
				aiTaskData[i].GenerateAudio = resp.GenerateAudio
				aiTaskData[i].Draft = resp.Draft
				aiTaskData[i].DraftTaskId = resp.DraftTaskID
				aiTaskData[i].ServiceTier = resp.ServiceTier
				aiTaskData[i].ExecutionExpiresAfter = resp.ExecutionExpiresAfter
				// 安全访问 Error，避免空指针
				if resp.Error != nil {
					aiTaskData[i].ErrorMessage = resp.Error.Message
					aiTaskData[i].ErrorCode = resp.Error.Code
				}
				if resp.Usage.TotalTokens > 0 && resp.Resolution != nil && resp.Duration != nil {
					if price, ok := ai.FindModelPrice(
						modelPrice,
						&resp,
					); ok {
						aiTaskData[i].VendorAmount, aiTaskData[i].RetailAmount = price.Calculate(int64(resp.Usage.TotalTokens), *resp.Duration)
					}
				}
				if err := global.GVA_DB.Model(&model.AiTask{}).
					Where("id = ?", aiTaskData[i].Id).
					Updates(&aiTaskData[i]).Error; err != nil {
					fmt.Printf("更新 AI 任务失败 (id=%v): %v\n", aiTaskData[i].Id, err)
				}
			}
		}
	}
}
