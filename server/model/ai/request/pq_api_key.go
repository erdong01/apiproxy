package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type PqApiKeySearch struct {
	request.PageInfo
	Status *int `json:"status" form:"status"`
}
