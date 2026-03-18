package ai

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	PqAiTaskRouter
	PqAiModelRouter
	PqApiKeyRouter
}

var (
	pqAiTaskApi  = api.ApiGroupApp.AiApiGroup.PqAiTaskApi
	pqAiModelApi = api.ApiGroupApp.AiApiGroup.PqAiModelApi
	pqApiKeyApi  = api.ApiGroupApp.AiApiGroup.PqApiKeyApi
)
