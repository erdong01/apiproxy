package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(ai.PqAiTask{}, ai.PqAiModel{}, ai.PqApiKey{})
	if err != nil {
		return err
	}
	return nil
}
