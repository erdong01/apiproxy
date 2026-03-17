package schedule

import (
	"context"
	"fmt"

	"time"

	"github.com/erdong01/kit"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	arkruntimeModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

var Schedule *kit.Schedule

func Init(ctx context.Context) {
	Schedule = kit.NewSchedule()
	Schedule.Run(ctx)

	var seedanceTask SeedanceTask
	Schedule.Add(&seedanceTask, time.Second*10, true)
}

type SeedanceTask struct {
}

func (that *SeedanceTask) OnTimer() {
	var aiTaskData []model.AiTask
	var aiTaskCount int64

	global.GVA_DB.Model(&model.AiTask{}).Where("status IS NULL").Count(&aiTaskCount)
	var limit int = 10
	var size int

	for aiTaskCount > int64(limit*size) {
		offset := size * limit
		global.GVA_DB.Model(&model.AiTask{}).
			Offset(offset).
			Limit(int(limit)).
			Where("status IS NULL").
			Find(&aiTaskData)
		for i := range aiTaskData {
			fmt.Println("aiTaskData[i].Key:", aiTaskData[i].Key)
			fmt.Println("aiTaskData[i].GenerateTaskId:", aiTaskData[i].GenerateTaskId)
			client := arkruntime.NewClientWithApiKey(aiTaskData[i].Key)
			ctx := context.Background()
			req := arkruntimeModel.GetContentGenerationTaskRequest{
				ID: aiTaskData[i].GenerateTaskId,
			}
			resp, err := client.GetContentGenerationTask(ctx, req)
			if err != nil {
				fmt.Printf("get content generation task error: %v\n", err)
				continue
			}
			if resp.Status != arkruntimeModel.StatusQueued && resp.Status != arkruntimeModel.StatusRunning {
				aiTaskData[i].Status = resp.Status
				aiTaskData[i].CompletionTokens = &resp.Usage.CompletionTokens
				aiTaskData[i].TotalTokens = &resp.Usage.TotalTokens
				aiTaskData[i].ErrorMessage = resp.Error.Code + " " + resp.Error.Message
				global.GVA_DB.Model(&model.AiTask{}).
					Where("id = ?", aiTaskData[i].Id).
					Updates(&aiTaskData[i])
			}
			fmt.Printf("%+v\n", resp)
		}
		size++

	}
}
