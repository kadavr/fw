package crud

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRange(ctx *gin.Context, limitDef, offsetDef int64) rangeConf {
	limit, _ := GetIntQueryParam[int64](ctx, "limit")
	if limit <= 0 {
		limit = limitDef
	}
	offset, _ := GetIntQueryParam[int64](ctx, "offset")
	if offset <= 0 {
		offset = offsetDef
	}
	return rangeConf{
		Limit:  limit,
		Offset: offset,
	}
}

func rangectx(ctx *gin.Context, limitDef, offsetDef int64, db *gorm.DB) (*gorm.DB, rangeConf) {
	rg := GetRange(ctx, limitDef, offsetDef)
	return db.Limit(int(rg.Limit)).Offset(int(rg.Offset)), rg
}
