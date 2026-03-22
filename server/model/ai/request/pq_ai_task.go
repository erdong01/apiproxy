package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type PqAiTaskSearch struct {
	request.PageInfo
	Status         string `json:"status" form:"status"`
	Key            string `json:"key" form:"key"`
	GenerateTaskId string `json:"generateTaskId" form:"generateTaskId"`
}
