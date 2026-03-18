import service from '@/utils/request'
// @Tags PqAiTask
// @Summary 创建pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PqAiTask true "创建pqAiTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /pqAiTask/createPqAiTask [post]
export const createPqAiTask = (data) => {
  return service({
    url: '/pqAiTask/createPqAiTask',
    method: 'post',
    data
  })
}

// @Tags PqAiTask
// @Summary 删除pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PqAiTask true "删除pqAiTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pqAiTask/deletePqAiTask [delete]
export const deletePqAiTask = (params) => {
  return service({
    url: '/pqAiTask/deletePqAiTask',
    method: 'delete',
    params
  })
}

// @Tags PqAiTask
// @Summary 批量删除pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除pqAiTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pqAiTask/deletePqAiTask [delete]
export const deletePqAiTaskByIds = (params) => {
  return service({
    url: '/pqAiTask/deletePqAiTaskByIds',
    method: 'delete',
    params
  })
}

// @Tags PqAiTask
// @Summary 更新pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PqAiTask true "更新pqAiTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /pqAiTask/updatePqAiTask [put]
export const updatePqAiTask = (data) => {
  return service({
    url: '/pqAiTask/updatePqAiTask',
    method: 'put',
    data
  })
}

// @Tags PqAiTask
// @Summary 用id查询pqAiTask表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.PqAiTask true "用id查询pqAiTask表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /pqAiTask/findPqAiTask [get]
export const findPqAiTask = (params) => {
  return service({
    url: '/pqAiTask/findPqAiTask',
    method: 'get',
    params
  })
}

// @Tags PqAiTask
// @Summary 分页获取pqAiTask表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取pqAiTask表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pqAiTask/getPqAiTaskList [get]
export const getPqAiTaskList = (params) => {
  return service({
    url: '/pqAiTask/getPqAiTaskList',
    method: 'get',
    params
  })
}

// @Tags PqAiTask
// @Summary 不需要鉴权的pqAiTask表接口
// @Accept application/json
// @Produce application/json
// @Param data query aiReq.PqAiTaskSearch true "分页获取pqAiTask表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /pqAiTask/getPqAiTaskPublic [get]
export const getPqAiTaskPublic = () => {
  return service({
    url: '/pqAiTask/getPqAiTaskPublic',
    method: 'get',
  })
}
