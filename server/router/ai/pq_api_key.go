package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PqApiKeyRouter struct {}

// InitPqApiKeyRouter 初始化 pqApiKey表 路由信息
func (s *PqApiKeyRouter) InitPqApiKeyRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	pqApiKeyRouter := Router.Group("pqApiKey").Use(middleware.OperationRecord())
	pqApiKeyRouterWithoutRecord := Router.Group("pqApiKey")
	pqApiKeyRouterWithoutAuth := PublicRouter.Group("pqApiKey")
	{
		pqApiKeyRouter.POST("createPqApiKey", pqApiKeyApi.CreatePqApiKey)   // 新建pqApiKey表
		pqApiKeyRouter.DELETE("deletePqApiKey", pqApiKeyApi.DeletePqApiKey) // 删除pqApiKey表
		pqApiKeyRouter.DELETE("deletePqApiKeyByIds", pqApiKeyApi.DeletePqApiKeyByIds) // 批量删除pqApiKey表
		pqApiKeyRouter.PUT("updatePqApiKey", pqApiKeyApi.UpdatePqApiKey)    // 更新pqApiKey表
	}
	{
		pqApiKeyRouterWithoutRecord.GET("findPqApiKey", pqApiKeyApi.FindPqApiKey)        // 根据ID获取pqApiKey表
		pqApiKeyRouterWithoutRecord.GET("getPqApiKeyList", pqApiKeyApi.GetPqApiKeyList)  // 获取pqApiKey表列表
	}
	{
	    pqApiKeyRouterWithoutAuth.GET("getPqApiKeyPublic", pqApiKeyApi.GetPqApiKeyPublic)  // pqApiKey表开放接口
	}
}
