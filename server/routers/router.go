package routers

import (
	"fmt"

	_ "github.com/bitCaskKV/docs"
	"github.com/bitCaskKV/global"
	"github.com/bitCaskKV/pkg/app"
	"github.com/bitCaskKV/pkg/errcode"
	"github.com/bitCaskKV/server/middleware"
	"github.com/bitCaskKV/server/service"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func NewRouter() *gin.Engine {

	if global.ServerSetting.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Recovery())
		r.Use(gin.Logger())
	} else {
		// Add Server Deploy Middleware
		r.Use(gin.Recovery())
		r.Use(middleware.AccessLog())
		// r.Use(gin.Logger())
	}

	r.POST("/v1/dbs", CreateDB)
	r.POST("/v1/db/:dbname", Put)
	r.GET("/v1/db/:dbname", Get)
	r.DELETE("/v1/db/:dbname", Del)

	//swagger

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

// ---------------- Router Method -------------------//

// @Summary Create DB Block
// @Description Create DB Block
// @Produce json
// @Param dbname query string true "block name"
// @Success 200 {string} string "ok" "成功"
// @Failure 400 {string} string "err_code：10000001 参数错误"；
// @Router /v1/dbs [post]
func CreateDB(c *gin.Context) {
	response := app.NewResponse(c)
	param := service.CreateRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		e := errcode.InvalidParams.WithDetails(err.Error(), fmt.Sprintf("No dbname in request parameter"))
		response.ToErrorResponse(e)
		return
	}

	block, err := global.DefaultBitCaskEngine.Block(param.DBName)

	if err != nil && block == nil {
		global.DefaultBitCaskEngine.NewBlock(param.DBName)
	}

	response.ToResponse(gin.H{})
}

// @Summary Put Key - Value in DB
// @Description Put Key - Value in DB
// @Produce json
// @Param dbname path string true "block name"
// @Param key query string true "put key"
// @Param dbname query string true "put value"
// @Success 200 {string} string "ok" "成功"
// @Failure 400 {string} string "err_code：10000001 参数错误"；
// @Failure 500 {string} string "err_code：10000000 服务错误 err_code: 10000002 无数据"；
// @Router /v1/db/:dbname [post]
func Put(c *gin.Context) {
	response := app.NewResponse(c)
	param := service.PutRequest{}
	name, ok := c.Params.Get("dbname")
	if !ok {
		e := errcode.InvalidParams.WithDetails("no Path Variable for dbname")
		response.ToErrorResponse(e)
		return
	}

	block, err := global.DefaultBitCaskEngine.Block(name)
	if err != nil {
		e := errcode.NotFound.WithDetails(fmt.Sprintf("Not Found DB Name %s", name))
		response.ToErrorResponse(e)
		return
	}
	err = c.ShouldBind(&param)
	if err != nil {
		e := errcode.InvalidParams.WithDetails(fmt.Sprintf("no Key and Value Param in Request"))
		response.ToErrorResponse(e)
		return
	}

	err = block.Put([]byte(param.Key), []byte(param.Value))
	if err != nil {
		e := errcode.ServerError.WithDetails("Write into DB Failed")
		response.ToErrorResponse(e)
		return
	}
	response.ToResponse(gin.H{})
}

// @Summary Get Key - Value in DB
// @Description Get Key - Value in DB
// @Produce json
// @Param dbname path string true "block name"
// @Param key query string true "Get key"
// @Success 200 {string} string "ok" "成功"
// @Failure 400 {string} string "err_code：10000001 参数错误"；
// @Failure 500 {string} string "err_code：10000000 服务错误 err_code: 10000002 无数据"；
// @Router /v1/db/:dbname [get]
func Get(c *gin.Context) {
	response := app.NewResponse(c)
	param := service.GetRequest{}
	name, ok := c.Params.Get("dbname")
	if !ok {
		e := errcode.InvalidParams.WithDetails("no Path Variable for dbname")
		response.ToErrorResponse(e)
		return
	}

	block, err := global.DefaultBitCaskEngine.Block(name)
	if err != nil {
		e := errcode.NotFound.WithDetails(fmt.Sprintf("Not Found DB Name %s", name))
		response.ToErrorResponse(e)
		return
	}
	err = c.ShouldBind(&param)
	if err != nil {
		e := errcode.InvalidParams.WithDetails(fmt.Sprintf("no Key Parameter in Request"))
		response.ToErrorResponse(e)
		return
	}

	v, err := block.Get([]byte(param.Key))
	if err != nil {
		e := errcode.ServerError.WithDetails(fmt.Sprintf("Get Key %s from DB %s Failed", param.Key, name))
		response.ToErrorResponse(e)
		return
	}
	response.ToResponse(gin.H{"value": string(v)})
}

// @Summary Delete Key - Value in DB
// @Description Delete Key - Value in DB
// @Produce json
// @Param dbname path string true "block name"
// @Param key query string true "Get key"
// @Success 200 {string} string "ok" "成功"
// @Failure 400 {string} string "err_code：10000001 参数错误"；
// @Failure 500 {string} string "err_code：10000000 服务错误 err_code: 10000002 无数据"；
// @Router /v1/db/:dbname [delete]
func Del(c *gin.Context) {
	response := app.NewResponse(c)
	param := service.DelRequest{}
	name, ok := c.Params.Get("dbname")
	if !ok {
		e := errcode.InvalidParams.WithDetails("no Path Variable for dbname")
		response.ToErrorResponse(e)
		return
	}

	block, err := global.DefaultBitCaskEngine.Block(name)
	if err != nil {
		e := errcode.NotFound.WithDetails(fmt.Sprintf("Not Found DB Name %s", name))
		response.ToErrorResponse(e)
		return
	}
	err = c.ShouldBind(&param)
	if err != nil {
		e := errcode.InvalidParams.WithDetails(fmt.Sprintf("no Key and Value Param in Request"))
		response.ToErrorResponse(e)
		return
	}

	err = block.Del([]byte(param.Key))
	if err != nil {
		e := errcode.ServerError.WithDetails(fmt.Sprintf("Del Key %s from DB %s Failed", string(param.Key), name))
		response.ToErrorResponse(e)
		return
	}
	response.ToResponse(gin.H{})
}
