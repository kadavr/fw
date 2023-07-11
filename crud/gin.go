package crud

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LogicFunc func(ctx *gin.Context, db *gorm.DB) *gorm.DB

type CRud[T Model] struct {
	stub         *T
	router       gin.IRouter
	db           *gorm.DB
	model        *gorm.DB
	lcmodelname  string
	idparamname  string
	responser    Responser[T]
	errResponser ErrorResponser

	ReadOneLogic   LogicFunc
	CreateOneLogic LogicFunc
	ReadManyLogic  LogicFunc
	UpdateOneLogic LogicFunc
	DeleteOneLogic LogicFunc

	readManyPrloader []string
}

type CrudOpt struct {
}

type ICrud interface {
}

func Crud[T Model](router gin.IRouter, childCrud ICrud, crudOpt ...CrudOpt) CRud[T] {
	modelStub := new(T)
	modelStubForName := new(T)
	modelName := getModleNameNoPtr(*modelStubForName)
	routeName := strings.ToLower(modelName)

	db := LoadGormDB()
	crd := CRud[T]{
		db:           db,
		lcmodelname:  routeName,
		idparamname:  routeName + "_id",
		router:       router,
		model:        db.Model(modelStub),
		stub:         modelStub,
		responser:    NewResponseRegistry[T](),
		errResponser: NewErrResponseRegistry(),
	}
	var preloaderNames []string
	{ //detect preloads

		v := reflect.ValueOf(*modelStub)
		for i := 0; i < v.NumField(); i++ {
			t := v.Type()
			f := t.Field(i)
			if strings.Contains(f.Tag.Get("gorm"), "foreignKey") {
				preloaderNames = append(preloaderNames, f.Name)
			}
		}
		crd.readManyPrloader = preloaderNames
	}
	crd.CrudRouter(router)

	return crd
}

func (c *CRud[T]) CrudRouter(router gin.IRouter) {
	rg := router.Group(fmt.Sprintf("/%s", c.lcmodelname))
	rg.GET("/", c.readMany)
	rg.POST("/", c.createOne)
	rg.GET(fmt.Sprintf("/:%s", c.idparamname), c.readOne)
	rg.PATCH(fmt.Sprintf("/:%s", c.idparamname), c.updateOne)
	rg.DELETE(fmt.Sprintf("/:%s", c.idparamname), c.deleteOne)
	c.router = rg
}

func (c *CRud[T]) createOne(ctx *gin.Context) {
	var body T
	if err := ctx.ShouldBind(&body); err != nil {
		errBody := c.errResponser.Match(err)
		ctx.JSON(400, errBody)
		return
	}
	ctxc, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel()
	result := c.model.WithContext(ctxc).Clauses(clause.Returning{}).Create(&body)
	if result.Error != nil {
		errBody := c.errResponser.Match(result.Error)
		ctx.JSON(int(errBody.Code), errBody)
		return
	}
	if result.RowsAffected > 0 {
		ctx.JSON(200, c.responser.ResponseOne(200, body))
		return
	}
}

func (c *CRud[T]) readMany(ctx *gin.Context) {
	var objects []readManywrapper[T]

	ctxc, cancel := context.WithCancel(ctx.Request.Context())

	defer cancel()
	result := c.model.WithContext(ctxc)
	for _, preload := range c.readManyPrloader {
		result = result.Preload(preload)
	}
	result.Select("*", "count(id) OVER() as count")
	result, rangeCfg := rangectx(ctx, 10, 0, result)
	result = result.Find(&objects)

	if result.Error != nil {
		errBody := c.errResponser.Match(result.Error)
		ctx.JSON(int(errBody.Status), errBody)
		return
	}

	var count int64
	if len(objects) != 0 {
		count = objects[0].Count
	}
	ctx.JSON(200, NewListResponse[readManywrapper[T]](200, rangeCfg, count, objects))
	return

}

func (c *CRud[T]) readOne(ctx *gin.Context) {
	id, err := GetIntPathParam[int64](ctx, c.idparamname)
	if err != nil {
		errBody := c.errResponser.Match(err)
		ctx.JSON(400, errBody)
		return
	}
	var object T
	ctxc, cancel := context.WithCancel(ctx.Request.Context())
	defer cancel()

	result := c.model.WithContext(ctxc).First(&object, id)
	if result.Error != nil {
		errBody := c.errResponser.Match(result.Error)
		ctx.JSON(int(errBody.Status), errBody)
		return
	}
	ctx.JSON(200, c.responser.ResponseOne(200, object))
	return
}

func (c *CRud[T]) updateOne(ctx *gin.Context) {
	id, err := GetIntPathParam[int64](ctx, c.idparamname)
	if err != nil {
		errBody := c.errResponser.Match(err)
		ctx.JSON(400, errBody)
		return
	}
	var body T
	if err := ctx.ShouldBind(&body); err != nil {
		errBody := c.errResponser.Match(err)
		ctx.JSON(400, errBody)
		return
	}

	ctxc, cancel := context.WithCancel(ctx)
	defer cancel()
	fmt.Printf("WTF2: %+v, %d", body, id)
	result := c.db.Session(&gorm.Session{FullSaveAssociations: true}).WithContext(ctxc).Save(&body)
	if result.Error != nil {
		errBody := c.errResponser.Match(result.Error)
		ctx.JSON(int(errBody.Status), errBody)
		return
	}
	ctx.JSON(200, c.responser.ResponseOne(200, body))
	return

}
func (c *CRud[T]) deleteOne(ctx *gin.Context) {}

// func (this *crud[T]) PrintRoutes() {
// 	for _, item := range this.router.Routes() {
// 		println("method:", item.Method, "path:", item.Path)
// 	}
// }

// func CrudNest[Tparent Model, TChild Model]() crudNest[Tparent, TChild] {
// 	modelStub := new(Tparent)
// 	_ = getName(*modelStub)

// 	(*modelStub).Identity()
// 	return crudNest[Tparent, TChild]{
// 		db: LoadGormDB(),
// 	}
// }

func (c *CRud[T]) Migrate(addit ...any) error {
	migrator := c.db.Migrator()
	args := []any{} //any(c.stub)
	args = append(args, addit...)
	return migrator.AutoMigrate(args...)

}

// API(http,grpc,...) -> Model -> Logic -> View
//       |<----------------------------------|
//	filter input   BaseQModel   JoinsEtc RenderView
