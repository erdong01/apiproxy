import service from '@/utils/request'
// @Tags PqApiKey
// @Summary 创建pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PqApiKey true "创建pqApiKey表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /pqApiKey/createPqApiKey [post]
export const createPqApiKey = (data) => {
  return service({
    url: '/pqApiKey/createPqApiKey',
    method: 'post',
    data
  })
}

// @Tags PqApiKey
// @Summary 删除pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PqApiKey true "删除pqApiKey表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pqApiKey/deletePqApiKey [delete]
export const deletePqApiKey = (params) => {
  return service({
    url: '/pqApiKey/deletePqApiKey',
    method: 'delete',
    params
  })
}

// @Tags PqApiKey
// @Summary 批量删除pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除pqApiKey表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pqApiKey/deletePqApiKey [delete]
export const deletePqApiKeyByIds = (params) => {
  return service({
    url: '/pqApiKey/deletePqApiKeyByIds',
    method: 'delete',
    params
  })
}

// @Tags PqApiKey
// @Summary 更新pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PqApiKey true "更新pqApiKey表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /pqApiKey/updatePqApiKey [put]
export const updatePqApiKey = (data) => {
  return service({
    url: '/pqApiKey/updatePqApiKey',
    method: 'put',
    data
  })
}

// @Tags PqApiKey
// @Summary 用id查询pqApiKey表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.PqApiKey true "用id查询pqApiKey表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /pqApiKey/findPqApiKey [get]
export const findPqApiKey = (params) => {
  return service({
    url: '/pqApiKey/findPqApiKey',
    method: 'get',
    params
  })
}

// @Tags PqApiKey
// @Summary 分页获取pqApiKey表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取pqApiKey表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pqApiKey/getPqApiKeyList [get]
export const getPqApiKeyList = (params) => {
  return service({
    url: '/pqApiKey/getPqApiKeyList',
    method: 'get',
    params
  })
}

// @Tags PqApiKey
// @Summary 不需要鉴权的pqApiKey表接口
// @Accept application/json
// @Produce application/json
// @Param data query aiReq.PqApiKeySearch true "分页获取pqApiKey表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /pqApiKey/getPqApiKeyPublic [get]
export const getPqApiKeyPublic = () => {
  return service({
    url: '/pqApiKey/getPqApiKeyPublic',
    method: 'get',
  })
}
