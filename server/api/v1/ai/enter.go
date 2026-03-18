package ai

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	PqAiTaskApi
	PqAiModelApi
	PqApiKeyApi
}

var (
	pqAiTaskService  = service.ServiceGroupApp.AiServiceGroup.PqAiTaskService
	pqAiModelService = service.ServiceGroupApp.AiServiceGroup.PqAiModelService
	pqApiKeyService  = service.ServiceGroupApp.AiServiceGroup.PqApiKeyService
)
