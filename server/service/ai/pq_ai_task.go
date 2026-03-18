
package ai

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
    aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
)

type PqAiTaskService struct {}
// CreatePqAiTask 创建pqAiTask表记录
// Author [yourname](https://github.com/yourname)
func (pqAiTaskService *PqAiTaskService) CreatePqAiTask(ctx context.Context, pqAiTask *ai.PqAiTask) (err error) {
	err = global.GVA_DB.Create(pqAiTask).Error
	return err
}

// DeletePqAiTask 删除pqAiTask表记录
// Author [yourname](https://github.com/yourname)
func (pqAiTaskService *PqAiTaskService)DeletePqAiTask(ctx context.Context, id string) (err error) {
	err = global.GVA_DB.Delete(&ai.PqAiTask{},"id = ?",id).Error
	return err
}

// DeletePqAiTaskByIds 批量删除pqAiTask表记录
// Author [yourname](https://github.com/yourname)
func (pqAiTaskService *PqAiTaskService)DeletePqAiTaskByIds(ctx context.Context, ids []string) (err error) {
	err = global.GVA_DB.Delete(&[]ai.PqAiTask{},"id in ?",ids).Error
	return err
}

// UpdatePqAiTask 更新pqAiTask表记录
// Author [yourname](https://github.com/yourname)
func (pqAiTaskService *PqAiTaskService)UpdatePqAiTask(ctx context.Context, pqAiTask ai.PqAiTask) (err error) {
	err = global.GVA_DB.Model(&ai.PqAiTask{}).Where("id = ?",pqAiTask.Id).Updates(&pqAiTask).Error
	return err
}

// GetPqAiTask 根据id获取pqAiTask表记录
// Author [yourname](https://github.com/yourname)
func (pqAiTaskService *PqAiTaskService)GetPqAiTask(ctx context.Context, id string) (pqAiTask ai.PqAiTask, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&pqAiTask).Error
	return
}
// GetPqAiTaskInfoList 分页获取pqAiTask表记录
// Author [yourname](https://github.com/yourname)
func (pqAiTaskService *PqAiTaskService)GetPqAiTaskInfoList(ctx context.Context, info aiReq.PqAiTaskSearch) (list []ai.PqAiTask, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&ai.PqAiTask{})
    var pqAiTasks []ai.PqAiTask
    // 如果有条件搜索 下方会自动创建搜索语句
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&pqAiTasks).Error
	return  pqAiTasks, total, err
}
func (pqAiTaskService *PqAiTaskService)GetPqAiTaskPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
