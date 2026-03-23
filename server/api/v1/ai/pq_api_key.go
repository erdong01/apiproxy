package ai

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/lib/apisix"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PqApiKeyApi struct{}

// CreatePqApiKey 创建pqApiKey表
// @Tags PqApiKey
// @Summary 创建pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ai.PqApiKey true "创建pqApiKey表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /pqApiKey/createPqApiKey [post]
func (pqApiKeyApi *PqApiKeyApi) CreatePqApiKey(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pqApiKey ai.PqApiKey
	err := c.ShouldBindJSON(&pqApiKey)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 生成一个给外部用户使用的key
	userKey := strings.ReplaceAll(uuid.New().String(), "-", "")
	pqApiKey.UserKey = &userKey

	err = pqApiKeyService.CreatePqApiKey(ctx, &pqApiKey)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	apisix.POSTConsumers(pqApiKey.UserName, *pqApiKey.Key, *pqApiKey.UserKey, pqApiKey.Rate)
	response.OkWithMessage("创建成功", c)
}

// DeletePqApiKey 删除pqApiKey表
// @Tags PqApiKey
// @Summary 删除pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ai.PqApiKey true "删除pqApiKey表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /pqApiKey/deletePqApiKey [delete]
func (pqApiKeyApi *PqApiKeyApi) DeletePqApiKey(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	id := c.Query("id")
	apiKeyData, err := pqApiKeyService.GetPqApiKey(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("数据不存在!", zap.Error(err))
		response.FailWithMessage("数据不存在:"+err.Error(), c)
		return
	}
	err = pqApiKeyService.DeletePqApiKey(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	if err := apisix.DELETEConsumers(apiKeyData.UserName); err != nil {
		global.GVA_LOG.Error("删除 APISIX consumer 失败!", zap.Error(err), zap.String("username", apiKeyData.UserName))
	}
	response.OkWithMessage("删除成功", c)
}

// DeletePqApiKeyByIds 批量删除pqApiKey表
// @Tags PqApiKey
// @Summary 批量删除pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /pqApiKey/deletePqApiKeyByIds [delete]
func (pqApiKeyApi *PqApiKeyApi) DeletePqApiKeyByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ids := c.QueryArray("ids[]")
	err := pqApiKeyService.DeletePqApiKeyByIds(ctx, ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdatePqApiKey 更新pqApiKey表
// @Tags PqApiKey
// @Summary 更新pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ai.PqApiKey true "更新pqApiKey表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /pqApiKey/updatePqApiKey [put]
func (pqApiKeyApi *PqApiKeyApi) UpdatePqApiKey(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var pqApiKey ai.PqApiKey
	err := c.ShouldBindJSON(&pqApiKey)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = pqApiKeyService.UpdatePqApiKey(ctx, pqApiKey)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}

	if pqApiKey.Status == 1 {
		if err := apisix.POSTConsumers(pqApiKey.UserName, *pqApiKey.Key, *pqApiKey.UserKey, pqApiKey.Rate); err != nil {
			global.GVA_LOG.Error("创建/更新 APISIX consumer 失败!", zap.Error(err), zap.String("username", pqApiKey.UserName))
		}
	} else if pqApiKey.Status == 2 {
		if err := apisix.DELETEConsumers(pqApiKey.UserName); err != nil {
			global.GVA_LOG.Error("删除 APISIX consumer 失败!", zap.Error(err), zap.String("username", pqApiKey.UserName))
		}
	}

	response.OkWithMessage("更新成功", c)
}

// FindPqApiKey 用id查询pqApiKey表
// @Tags PqApiKey
// @Summary 用id查询pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query int true "用id查询pqApiKey表"
// @Success 200 {object} response.Response{data=ai.PqApiKey,msg=string} "查询成功"
// @Router /pqApiKey/findPqApiKey [get]
func (pqApiKeyApi *PqApiKeyApi) FindPqApiKey(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	id := c.Query("id")
	repqApiKey, err := pqApiKeyService.GetPqApiKey(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(repqApiKey, c)
}

// GetPqApiKeyList 分页获取pqApiKey表列表
// @Tags PqApiKey
// @Summary 分页获取pqApiKey表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query aiReq.PqApiKeySearch true "分页获取pqApiKey表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /pqApiKey/getPqApiKeyList [get]
func (pqApiKeyApi *PqApiKeyApi) GetPqApiKeyList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo aiReq.PqApiKeySearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pqApiKeyService.GetPqApiKeyInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetPqApiKeyPublic 不需要鉴权的pqApiKey表接口
// @Tags PqApiKey
// @Summary 不需要鉴权的pqApiKey表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /pqApiKey/getPqApiKeyPublic [get]
func (pqApiKeyApi *PqApiKeyApi) GetPqApiKeyPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	pqApiKeyService.GetPqApiKeyPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的pqApiKey表接口信息",
	}, "获取成功", c)
}
