package routers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bitCaskKV/global"
	"github.com/bitCaskKV/pkg/app"
	"github.com/bitCaskKV/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Recovery())
		r.Use(gin.Logger())
	} else {
		// Add Server Deploy Middleware
		r.Use(gin.Recovery())
		// r.Use(gin.Logger())
	}

	// r.POST("/v1/db", func(c *gin.Context) {
	// 	response := app.NewResponse(c)
	// 	op := c.Query("Operate")
	// 	c.FullPath()
	// 	if op == "" {
	// 		response.ToErrorResponse(errcode.InvalidParams.WithDetails(fmt.Sprintf("No Valid Operate Param")))
	// 		return
	// 	}
	// 	switch op {
	// 	case "Put":
	// 		DBPut(c, response)
	// 	case "Get":
	// 		DBGet(c, response)
	// 	case "Del":
	// 		DBDel(c, response)
	// 	default:
	// 		response.ToErrorResponse(errcode.InvalidParams.WithDetails(fmt.Sprintf("No Valid Operate Param")))
	// 	}
	// })

	r.POST("/v1/dbs", func(c *gin.Context) {
		db, ok := c.GetQuery("dbname")
		if !ok {
			c.JSON(errcode.InvalidParams.StatusCode(), gin.H{"msg": fmt.Sprintf("")})
			return
		}
		fmt.Println(global.DefaultBitCaskEngine)
		block, err := global.DefaultBitCaskEngine.Block(db)

		if err != nil && block == nil {
			global.DefaultBitCaskEngine.NewBlock(db)
		}

		c.JSON(http.StatusOK, gin.H{})

	})

	r.POST("/v1/db/:dbname", func(c *gin.Context) {
		name, ok := c.Params.Get("dbname")
		if !ok {
			c.JSON(errcode.InvalidParams.StatusCode(), gin.H{"msg": fmt.Sprintf("no Path Variable for dbname")})
			return
		}
		log.Println("QQQQQ")
		fmt.Println(global.DefaultBitCaskEngine)
		block, err := global.DefaultBitCaskEngine.Block(name)
		if err != nil {
			c.JSON(errcode.NotFound.StatusCode(), gin.H{"msg": fmt.Sprintf("Not Found DB Name %s", name)})
			return
		}
		key, ok1 := c.GetQuery("Key")
		value, ok2 := c.GetQuery("Value")
		if !ok1 || !ok2 {
			c.JSON(errcode.InvalidParams.StatusCode(), gin.H{"msg": fmt.Sprintf("no Key and Value Param in Request")})
			return
		}
		err = block.Put([]byte(key), []byte(value))
		if err != nil {
			c.JSON(errcode.ServerError.StatusCode(), gin.H{"msg": errcode.ServerError.Msg()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/v1/db/:dbname", func(c *gin.Context) {
		name, ok := c.Params.Get("dbname")
		if !ok {
			c.JSON(errcode.InvalidParams.StatusCode(), gin.H{"msg": fmt.Sprintf("no Path Variable for dbname")})
			return
		}
		block, err := global.DefaultBitCaskEngine.Block(name)
		if err != nil {
			c.JSON(errcode.NotFound.StatusCode(), gin.H{"msg": fmt.Sprintf("Not Found DB Name %s", name)})
			return
		}

		key, ok := c.GetQuery("Key")
		if !ok {
			c.JSON(errcode.InvalidParams.StatusCode(), gin.H{"msg": fmt.Sprintf("no Key and Value Param in Request")})
			return
		}

		v, err := block.Get([]byte(key))
		if err != nil {
			c.JSON(errcode.ServerError.StatusCode(), gin.H{"msg": errcode.ServerError.Msg()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Value": string(v)})

	})

	r.DELETE("/v1/db/:dbname", func(c *gin.Context) {
		name, ok := c.Params.Get("dbname")
		if !ok {
			c.JSON(errcode.InvalidParams.StatusCode(), gin.H{"msg": fmt.Sprintf("no Path Variable for dbname")})
			return
		}
		block, err := global.DefaultBitCaskEngine.Block(name)
		if err != nil {
			c.JSON(errcode.NotFound.StatusCode(), gin.H{"msg": fmt.Sprintf("Not Found DB Name %s", name)})
			return
		}

		key, ok := c.GetQuery("Key")
		if !ok {
			c.JSON(errcode.InvalidParams.StatusCode(), gin.H{"msg": fmt.Sprintf("no Key and Value Param in Request")})
			return
		}

		err = block.Del([]byte(key))
		if err != nil {
			c.JSON(errcode.ServerError.StatusCode(), gin.H{"msg": errcode.ServerError.Msg()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	return r
}

func DBPut(c *gin.Context, response *app.Response) {
	key, ok1 := c.GetQuery("Key")
	value, ok2 := c.GetQuery("Value")
	if !ok1 || !ok2 {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(fmt.Sprintf("No Key and Value Param")))
		return
	}

	response.ToResponse(gin.H{"Hello": "Put", key: value})
}

func DBGet(c *gin.Context, response *app.Response) {
	key, ok := c.GetQuery("Key")
	if ok {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(fmt.Sprintf("No Key Param")))
		return
	}
	response.ToResponse(gin.H{key: "Get"})
}

func DBDel(c *gin.Context, response *app.Response) {
	key, ok := c.GetQuery("Key")
	if ok {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(fmt.Sprintf("No Key Param")))
		return
	}
	response.ToResponse(gin.H{key: "Del"})

}
