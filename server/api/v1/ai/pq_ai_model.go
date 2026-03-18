package ai

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/ai"
    aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type PqAiModelApi struct {}



// CreatePqAiModel 创建pqAiModel表
// @Tags PqAiModel
// @Summary 创建pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ai.PqAiModel true "创建pqAiModel表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /pqAiModel/createPqAiModel [post]
func (pqAiModelApi *PqAiModelApi) CreatePqAiModel(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pqAiModel ai.PqAiModel
	err := c.ShouldBindJSON(&pqAiModel)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = pqAiModelService.CreatePqAiModel(ctx,&pqAiModel)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeletePqAiModel 删除pqAiModel表
// @Tags PqAiModel
// @Summary 删除pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ai.PqAiModel true "删除pqAiModel表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /pqAiModel/deletePqAiModel [delete]
func (pqAiModelApi *PqAiModelApi) DeletePqAiModel(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	id := c.Query("id")
	err := pqAiModelService.DeletePqAiModel(ctx,id)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeletePqAiModelByIds 批量删除pqAiModel表
// @Tags PqAiModel
// @Summary 批量删除pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /pqAiModel/deletePqAiModelByIds [delete]
func (pqAiModelApi *PqAiModelApi) DeletePqAiModelByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ids := c.QueryArray("ids[]")
	err := pqAiModelService.DeletePqAiModelByIds(ctx,ids)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdatePqAiModel 更新pqAiModel表
// @Tags PqAiModel
// @Summary 更新pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ai.PqAiModel true "更新pqAiModel表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /pqAiModel/updatePqAiModel [put]
func (pqAiModelApi *PqAiModelApi) UpdatePqAiModel(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var pqAiModel ai.PqAiModel
	err := c.ShouldBindJSON(&pqAiModel)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = pqAiModelService.UpdatePqAiModel(ctx,pqAiModel)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindPqAiModel 用id查询pqAiModel表
// @Tags PqAiModel
// @Summary 用id查询pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query int true "用id查询pqAiModel表"
// @Success 200 {object} response.Response{data=ai.PqAiModel,msg=string} "查询成功"
// @Router /pqAiModel/findPqAiModel [get]
func (pqAiModelApi *PqAiModelApi) FindPqAiModel(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	id := c.Query("id")
	repqAiModel, err := pqAiModelService.GetPqAiModel(ctx,id)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(repqAiModel, c)
}
// GetPqAiModelList 分页获取pqAiModel表列表
// @Tags PqAiModel
// @Summary 分页获取pqAiModel表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query aiReq.PqAiModelSearch true "分页获取pqAiModel表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /pqAiModel/getPqAiModelList [get]
func (pqAiModelApi *PqAiModelApi) GetPqAiModelList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo aiReq.PqAiModelSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pqAiModelService.GetPqAiModelInfoList(ctx,pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}

// GetPqAiModelPublic 不需要鉴权的pqAiModel表接口
// @Tags PqAiModel
// @Summary 不需要鉴权的pqAiModel表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /pqAiModel/getPqAiModelPublic [get]
func (pqAiModelApi *PqAiModelApi) GetPqAiModelPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    pqAiModelService.GetPqAiModelPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的pqAiModel表接口信息",
    }, "获取成功", c)
}
