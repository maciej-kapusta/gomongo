package web

import (
	"github.com/gin-gonic/gin"
	"github.com/maciej-kapusta/gomongo/config"
)

type Server struct {
	port       string
	engine     *gin.Engine
	controller *DocController
}

func New(cfg *config.Config) (*Server, error) {
	controller, err := NewDocController(cfg)
	if err != nil {
		return nil, err
	}

	r := gin.Default()
	r.POST("/doc", controller.PostDoc)
	r.GET("/doc/:doc", controller.ReadDoc)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})
	return &Server{
		controller: controller,
		engine:     r,
		port:       cfg.Port,
	}, nil
}

func (s *Server) Serve() error {
	port := ":" + s.port
	return s.engine.Run(port)
}
