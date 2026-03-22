package ai

import (
	"context"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"gorm.io/gorm"
)

type PqAiModelService struct{}

// CreatePqAiModel 创建pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService) CreatePqAiModel(ctx context.Context, pqAiModel *ai.PqAiModel) (err error) {
	db := global.GVA_DB.WithContext(ctx)
	prices := cloneModelPrices(pqAiModel.PqAiModelPrice)

	err = db.Transaction(func(tx *gorm.DB) error {
		pqAiModel.PqAiModelPrice = nil
		if err := tx.Omit("PqAiModelPrice").Create(pqAiModel).Error; err != nil {
			return err
		}
		if pqAiModel.Id == nil || len(prices) == 0 {
			return nil
		}
		for i := range prices {
			prices[i].Id = 0
			prices[i].PqAiModelId = *pqAiModel.Id
		}
		return tx.Create(&prices).Error
	})
	return err
}

// DeletePqAiModel 删除pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService) DeletePqAiModel(ctx context.Context, id string) (err error) {
	db := global.GVA_DB.WithContext(ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("pq_ai_model_id = ?", id).Delete(&ai.PqAiModelPrice{}).Error; err != nil {
			return err
		}
		return tx.Delete(&ai.PqAiModel{}, "id = ?", id).Error
	})
	return err
}

// DeletePqAiModelByIds 批量删除pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService) DeletePqAiModelByIds(ctx context.Context, ids []string) (err error) {
	db := global.GVA_DB.WithContext(ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		if len(ids) == 0 {
			return nil
		}
		if err := tx.Where("pq_ai_model_id in ?", ids).Delete(&ai.PqAiModelPrice{}).Error; err != nil {
			return err
		}
		return tx.Delete(&[]ai.PqAiModel{}, "id in ?", ids).Error
	})
	return err
}

// UpdatePqAiModel 更新pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService) UpdatePqAiModel(ctx context.Context, pqAiModel ai.PqAiModel) (err error) {
	if pqAiModel.Id == nil {
		return errors.New("id不能为空")
	}

	db := global.GVA_DB.WithContext(ctx)
	prices := cloneModelPrices(pqAiModel.PqAiModelPrice)
	pqAiModel.PqAiModelPrice = nil

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ai.PqAiModel{}).Where("id = ?", *pqAiModel.Id).Omit("PqAiModelPrice").Updates(&pqAiModel).Error; err != nil {
			return err
		}
		if err := tx.Where("pq_ai_model_id = ?", *pqAiModel.Id).Delete(&ai.PqAiModelPrice{}).Error; err != nil {
			return err
		}
		if len(prices) == 0 {
			return nil
		}
		for i := range prices {
			prices[i].Id = 0
			prices[i].PqAiModelId = *pqAiModel.Id
		}
		return tx.Create(&prices).Error
	})
	return err
}

// GetPqAiModel 根据id获取pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService) GetPqAiModel(ctx context.Context, id string) (pqAiModel ai.PqAiModel, err error) {
	err = global.GVA_DB.WithContext(ctx).Preload("PqAiModelPrice").Where("id = ?", id).First(&pqAiModel).Error
	return
}

// GetPqAiModelInfoList 分页获取pqAiModel表记录
// Author [yourname](https://github.com/yourname)
func (pqAiModelService *PqAiModelService) GetPqAiModelInfoList(ctx context.Context, info aiReq.PqAiModelSearch) (list []ai.PqAiModel, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.WithContext(ctx).Model(&ai.PqAiModel{})
	var pqAiModels []ai.PqAiModel

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Preload("PqAiModelPrice").Find(&pqAiModels).Error
	return pqAiModels, total, err
}

func (pqAiModelService *PqAiModelService) GetPqAiModelPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

func cloneModelPrices(prices []ai.PqAiModelPrice) []ai.PqAiModelPrice {
	if len(prices) == 0 {
		return nil
	}
	cloned := make([]ai.PqAiModelPrice, 0, len(prices))
	for _, price := range prices {
		cloned = append(cloned, price)
	}
	return cloned
}
