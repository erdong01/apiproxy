import service from '@/utils/request'
// @Tags PqAiModel
// @Summary 创建pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PqAiModel true "创建pqAiModel表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /pqAiModel/createPqAiModel [post]
export const createPqAiModel = (data) => {
  return service({
    url: '/pqAiModel/createPqAiModel',
    method: 'post',
    data
  })
}

// @Tags PqAiModel
// @Summary 删除pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PqAiModel true "删除pqAiModel表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pqAiModel/deletePqAiModel [delete]
export const deletePqAiModel = (params) => {
  return service({
    url: '/pqAiModel/deletePqAiModel',
    method: 'delete',
    params
  })
}

// @Tags PqAiModel
// @Summary 批量删除pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除pqAiModel表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pqAiModel/deletePqAiModel [delete]
export const deletePqAiModelByIds = (params) => {
  return service({
    url: '/pqAiModel/deletePqAiModelByIds',
    method: 'delete',
    params
  })
}

// @Tags PqAiModel
// @Summary 更新pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PqAiModel true "更新pqAiModel表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /pqAiModel/updatePqAiModel [put]
export const updatePqAiModel = (data) => {
  return service({
    url: '/pqAiModel/updatePqAiModel',
    method: 'put',
    data
  })
}

// @Tags PqAiModel
// @Summary 用id查询pqAiModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.PqAiModel true "用id查询pqAiModel表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /pqAiModel/findPqAiModel [get]
export const findPqAiModel = (params) => {
  return service({
    url: '/pqAiModel/findPqAiModel',
    method: 'get',
    params
  })
}

// @Tags PqAiModel
// @Summary 分页获取pqAiModel表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取pqAiModel表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pqAiModel/getPqAiModelList [get]
export const getPqAiModelList = (params) => {
  return service({
    url: '/pqAiModel/getPqAiModelList',
    method: 'get',
    params
  })
}

// @Tags PqAiModel
// @Summary 不需要鉴权的pqAiModel表接口
// @Accept application/json
// @Produce application/json
// @Param data query aiReq.PqAiModelSearch true "分页获取pqAiModel表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /pqAiModel/getPqAiModelPublic [get]
export const getPqAiModelPublic = () => {
  return service({
    url: '/pqAiModel/getPqAiModelPublic',
    method: 'get',
  })
}
