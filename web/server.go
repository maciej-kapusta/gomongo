package web

import (
	"github.com/gin-gonic/gin"
	"github.com/maciej-kapusta/gomongo/config"
	"github.com/maciej-kapusta/gomongo/repo"
)

func SetupAll(cfg *config.Config) (*gin.Engine, error) {
	mongoRepo, err := repo.Connect[Doc](cfg.MongoUri, cfg.MongoDb, "docs")
	if err != nil {
		return nil, err
	}
	handler := NewDocController(mongoRepo)
	if err != nil {
		return nil, err
	}

	server := setupServer(handler)
	return server, nil
}

func setupServer(handler *DocHandler) *gin.Engine {

	r := gin.Default()
	r.Use(errorHandler)

	r.POST("/doc", handler.PostDoc)
	r.GET("/doc/:doc", handler.ReadDoc)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

func errorHandler(c *gin.Context) {
	c.Next()
	last := c.Errors.Last()
	if last != nil {
		c.Status(500)
		_, _ = c.Writer.Write([]byte(last.Error()))
	}
}
