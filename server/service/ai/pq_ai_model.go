
package ai

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
    aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
)

type PqAiModelService struct {}
// CreatePqAiModel 创建pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService) CreatePqAiModel(ctx context.Context, pqAiModel *ai.PqAiModel) (err error) {
	err = global.GVA_DB.Create(pqAiModel).Error
	return err
}

// DeletePqAiModel 删除pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService)DeletePqAiModel(ctx context.Context, id string) (err error) {
	err = global.GVA_DB.Delete(&ai.PqAiModel{},"id = ?",id).Error
	return err
}

// DeletePqAiModelByIds 批量删除pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService)DeletePqAiModelByIds(ctx context.Context, ids []string) (err error) {
	err = global.GVA_DB.Delete(&[]ai.PqAiModel{},"id in ?",ids).Error
	return err
}

// UpdatePqAiModel 更新pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService)UpdatePqAiModel(ctx context.Context, pqAiModel ai.PqAiModel) (err error) {
	err = global.GVA_DB.Model(&ai.PqAiModel{}).Where("id = ?",pqAiModel.Id).Updates(&pqAiModel).Error
	return err
}

// GetPqAiModel 根据id获取pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService)GetPqAiModel(ctx context.Context, id string) (pqAiModel ai.PqAiModel, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&pqAiModel).Error
	return
}
// GetPqAiModelInfoList 分页获取pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService)GetPqAiModelInfoList(ctx context.Context, info aiReq.PqAiModelSearch) (list []ai.PqAiModel, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&ai.PqAiModel{})
    var pqAiModels []ai.PqAiModel
    // 如果有条件搜索 下方会自动创建搜索语句
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&pqAiModels).Error
	return  pqAiModels, total, err
}
func (pqAiModelService *PqAiModelService)GetPqAiModelPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
