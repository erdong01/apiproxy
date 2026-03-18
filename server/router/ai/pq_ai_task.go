package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PqAiTaskRouter struct {}

// InitPqAiTaskRouter 初始化 pqAiTask表 路由信息
func (s *PqAiTaskRouter) InitPqAiTaskRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	pqAiTaskRouter := Router.Group("pqAiTask").Use(middleware.OperationRecord())
	pqAiTaskRouterWithoutRecord := Router.Group("pqAiTask")
	pqAiTaskRouterWithoutAuth := PublicRouter.Group("pqAiTask")
	{
		pqAiTaskRouter.POST("createPqAiTask", pqAiTaskApi.CreatePqAiTask)   // 新建pqAiTask表
		pqAiTaskRouter.DELETE("deletePqAiTask", pqAiTaskApi.DeletePqAiTask) // 删除pqAiTask表
		pqAiTaskRouter.DELETE("deletePqAiTaskByIds", pqAiTaskApi.DeletePqAiTaskByIds) // 批量删除pqAiTask表
		pqAiTaskRouter.PUT("updatePqAiTask", pqAiTaskApi.UpdatePqAiTask)    // 更新pqAiTask表
	}
	{
		pqAiTaskRouterWithoutRecord.GET("findPqAiTask", pqAiTaskApi.FindPqAiTask)        // 根据ID获取pqAiTask表
		pqAiTaskRouterWithoutRecord.GET("getPqAiTaskList", pqAiTaskApi.GetPqAiTaskList)  // 获取pqAiTask表列表
	}
	{
	    pqAiTaskRouterWithoutAuth.GET("getPqAiTaskPublic", pqAiTaskApi.GetPqAiTaskPublic)  // pqAiTask表开放接口
	}
}
