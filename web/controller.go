package web

import (
	"github.com/gin-gonic/gin"
	"github.com/maciej-kapusta/gomongo/config"
	"github.com/maciej-kapusta/gomongo/repo"
	"github.com/rs/zerolog/log"
)

type Doc struct {
	Name string `json:"name,omitempty"`
	Info string `json:"info,omitempty"`
}
type DocController struct {
	repo repo.Repo[Doc]
}

func NewDocController(config *config.Config) (*DocController, error) {
	mongoRepo, err := repo.Connect[Doc](config.MongoUri, config.MongoDb, "docs")
	if err != nil {
		return nil, err
	}

	return &DocController{
		repo: mongoRepo,
	}, nil
}

func (d *DocController) PostDoc(c *gin.Context) {
	var doc Doc
	err := c.BindJSON(&doc)
	if err != nil {
		log.Err(err).Msg("Could not decode json")
		_ = c.Error(err)
		return
	}

	id, err := d.repo.SaveObject(&doc)
	c.String(200, id)
}

func (d *DocController) ReadDoc(c *gin.Context) {
	id := c.Param("doc")
	object, err := d.repo.ReadObject(id)
	if err != nil {
		log.Err(err).Msg("Could not read object")
		_ = c.Error(err)
		return
	}
	c.JSON(200, object)
}
