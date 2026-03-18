package ai

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/ai"
    aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type PqAiTaskApi struct {}



// CreatePqAiTask 创建pqAiTask表
// @Tags PqAiTask
// @Summary 创建pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ai.PqAiTask true "创建pqAiTask表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /pqAiTask/createPqAiTask [post]
func (pqAiTaskApi *PqAiTaskApi) CreatePqAiTask(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pqAiTask ai.PqAiTask
	err := c.ShouldBindJSON(&pqAiTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = pqAiTaskService.CreatePqAiTask(ctx,&pqAiTask)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeletePqAiTask 删除pqAiTask表
// @Tags PqAiTask
// @Summary 删除pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ai.PqAiTask true "删除pqAiTask表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /pqAiTask/deletePqAiTask [delete]
func (pqAiTaskApi *PqAiTaskApi) DeletePqAiTask(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	id := c.Query("id")
	err := pqAiTaskService.DeletePqAiTask(ctx,id)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeletePqAiTaskByIds 批量删除pqAiTask表
// @Tags PqAiTask
// @Summary 批量删除pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /pqAiTask/deletePqAiTaskByIds [delete]
func (pqAiTaskApi *PqAiTaskApi) DeletePqAiTaskByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ids := c.QueryArray("ids[]")
	err := pqAiTaskService.DeletePqAiTaskByIds(ctx,ids)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdatePqAiTask 更新pqAiTask表
// @Tags PqAiTask
// @Summary 更新pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ai.PqAiTask true "更新pqAiTask表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /pqAiTask/updatePqAiTask [put]
func (pqAiTaskApi *PqAiTaskApi) UpdatePqAiTask(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var pqAiTask ai.PqAiTask
	err := c.ShouldBindJSON(&pqAiTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = pqAiTaskService.UpdatePqAiTask(ctx,pqAiTask)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindPqAiTask 用id查询pqAiTask表
// @Tags PqAiTask
// @Summary 用id查询pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query int true "用id查询pqAiTask表"
// @Success 200 {object} response.Response{data=ai.PqAiTask,msg=string} "查询成功"
// @Router /pqAiTask/findPqAiTask [get]
func (pqAiTaskApi *PqAiTaskApi) FindPqAiTask(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	id := c.Query("id")
	repqAiTask, err := pqAiTaskService.GetPqAiTask(ctx,id)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(repqAiTask, c)
}
// GetPqAiTaskList 分页获取pqAiTask表列表
// @Tags PqAiTask
// @Summary 分页获取pqAiTask表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query aiReq.PqAiTaskSearch true "分页获取pqAiTask表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /pqAiTask/getPqAiTaskList [get]
func (pqAiTaskApi *PqAiTaskApi) GetPqAiTaskList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo aiReq.PqAiTaskSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := pqAiTaskService.GetPqAiTaskInfoList(ctx,pageInfo)
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

// GetPqAiTaskPublic 不需要鉴权的pqAiTask表接口
// @Tags PqAiTask
// @Summary 不需要鉴权的pqAiTask表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /pqAiTask/getPqAiTaskPublic [get]
func (pqAiTaskApi *PqAiTaskApi) GetPqAiTaskPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    pqAiTaskService.GetPqAiTaskPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的pqAiTask表接口信息",
    }, "获取成功", c)
}
