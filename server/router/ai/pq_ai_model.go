package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PqAiModelRouter struct {}

// InitPqAiModelRouter 初始化 pqAiModel表 路由信息
func (s *PqAiModelRouter) InitPqAiModelRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	pqAiModelRouter := Router.Group("pqAiModel").Use(middleware.OperationRecord())
	pqAiModelRouterWithoutRecord := Router.Group("pqAiModel")
	pqAiModelRouterWithoutAuth := PublicRouter.Group("pqAiModel")
	{
		pqAiModelRouter.POST("createPqAiModel", pqAiModelApi.CreatePqAiModel)   // 新建pqAiModel表
		pqAiModelRouter.DELETE("deletePqAiModel", pqAiModelApi.DeletePqAiModel) // 删除pqAiModel表
		pqAiModelRouter.DELETE("deletePqAiModelByIds", pqAiModelApi.DeletePqAiModelByIds) // 批量删除pqAiModel表
		pqAiModelRouter.PUT("updatePqAiModel", pqAiModelApi.UpdatePqAiModel)    // 更新pqAiModel表
	}
	{
		pqAiModelRouterWithoutRecord.GET("findPqAiModel", pqAiModelApi.FindPqAiModel)        // 根据ID获取pqAiModel表
		pqAiModelRouterWithoutRecord.GET("getPqAiModelList", pqAiModelApi.GetPqAiModelList)  // 获取pqAiModel表列表
	}
	{
	    pqAiModelRouterWithoutAuth.GET("getPqAiModelPublic", pqAiModelApi.GetPqAiModelPublic)  // pqAiModel表开放接口
	}
}
