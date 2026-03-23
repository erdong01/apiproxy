package ai

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
)

type PqApiKeyService struct{}

// CreatePqApiKey 创建pqApiKey表记录
// Author [yourname](https://github.com/yourname)
func (pqApiKeyService *PqApiKeyService) CreatePqApiKey(ctx context.Context, pqApiKey *ai.PqApiKey) (err error) {
	err = global.GVA_DB.Create(pqApiKey).Error
	return err
}

// DeletePqApiKey 删除pqApiKey表记录
// Author [yourname](https://github.com/yourname)
func (pqApiKeyService *PqApiKeyService) DeletePqApiKey(ctx context.Context, id string) (err error) {
	err = global.GVA_DB.Delete(&ai.PqApiKey{}, "id = ?", id).Error
	return err
}

// DeletePqApiKeyByIds 批量删除pqApiKey表记录
// Author [yourname](https://github.com/yourname)
func (pqApiKeyService *PqApiKeyService) DeletePqApiKeyByIds(ctx context.Context, ids []string) (err error) {
	err = global.GVA_DB.Delete(&[]ai.PqApiKey{}, "id in ?", ids).Error
	return err
}

// UpdatePqApiKey 更新pqApiKey表记录
// Author [yourname](https://github.com/yourname)
func (pqApiKeyService *PqApiKeyService) UpdatePqApiKey(ctx context.Context, pqApiKey ai.PqApiKey) (err error) {
	err = global.GVA_DB.Model(&ai.PqApiKey{}).Where("id = ?", pqApiKey.Id).Updates(&pqApiKey).Error
	return err
}

// GetPqApiKey 根据id获取pqApiKey表记录
// Author [yourname](https://github.com/yourname)
func (pqApiKeyService *PqApiKeyService) GetPqApiKey(ctx context.Context, id string) (pqApiKey ai.PqApiKey, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&pqApiKey).Error
	return
}

// GetPqApiKeyInfoList 分页获取pqApiKey表记录
// Author [yourname](https://github.com/yourname)
func (pqApiKeyService *PqApiKeyService) GetPqApiKeyInfoList(ctx context.Context, info aiReq.PqApiKeySearch) (list []ai.PqApiKey, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&ai.PqApiKey{})
	var pqApiKeys []ai.PqApiKey
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Order("id DESC").Find(&pqApiKeys).Error
	return pqApiKeys, total, err
}

func (pqApiKeyService *PqApiKeyService) GetPqApiKeyPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
