package crud

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	fiber "github.com/gofiber/fiber/v2"
	"golang.org/x/exp/constraints"
)

func GetIntPathParam[T constraints.Integer](ctx *gin.Context, key string) (ret T, err error) {
	paramPath := ctx.Param(key)
	switch any(ret).(type) {
	case int, int16, int32, int64:
		val, err := strconv.ParseInt(paramPath, 10, 64)
		if err != nil {
			return T(0), err
		}
		ret = T(val)
	case uint, uint16, uint32, uint64:
		val, err := strconv.ParseUint(paramPath, 10, 64)
		if err != nil {
			return T(0), err
		}
		ret = T(val)
	}
	return
}
func GetIntQueryParam[T constraints.Integer](ctx *gin.Context, key string) (ret T, err error) {
	paramPath := ctx.Param(key)
	switch any(ret).(type) {
	case int, int16, int32, int64:
		val, err := strconv.ParseInt(paramPath, 10, 64)
		if err != nil {
			return T(0), err
		}
		ret = T(val)
	case uint, uint16, uint32, uint64:
		val, err := strconv.ParseUint(paramPath, 10, 64)
		if err != nil {
			return T(0), err
		}
		ret = T(val)
	}
	return
}
func GetIntPathParamFiber[T constraints.Integer](ctx *fiber.Ctx, key string) (ret T, err error) {
	paramPath := ctx.Params(key)
	switch any(ret).(type) {
	case int, int16, int32, int64:
		val, err := strconv.ParseInt(paramPath, 10, 64)
		if err != nil {
			return T(0), err
		}
		ret = T(val)
	case uint, uint16, uint32, uint64:
		val, err := strconv.ParseUint(paramPath, 10, 64)
		if err != nil {
			return T(0), err
		}
		ret = T(val)
	}
	return
}

func getModleNameNoPtr(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func getName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
