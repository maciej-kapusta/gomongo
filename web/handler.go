package web

import (
	"github.com/gin-gonic/gin"
	"github.com/maciej-kapusta/gomongo/repo"
	"github.com/rs/zerolog/log"
)

type Doc struct {
	Name string `json:"name,omitempty"`
	Info string `json:"info,omitempty"`
}

type DocHandler struct {
	repo repo.Repo[Doc]
}

func NewDocController(docRepo repo.Repo[Doc]) *DocHandler {
	return &DocHandler{
		repo: docRepo,
	}
}

func (d *DocHandler) PostDoc(c *gin.Context) {
	var doc Doc
	err := c.BindJSON(&doc)
	if err != nil {
		log.Err(err).Msg("Could not decode json")
		_ = c.Error(err)
		return
	}

	id, err := d.repo.SaveObject(&doc)
	if err != nil {
		log.Err(err).Msg("Could not save error")
		_ = c.Error(err)
		return
	}
	c.String(200, id)
}

func (d *DocHandler) ReadDoc(c *gin.Context) {
	id := c.Param("doc")
	object, err := d.repo.ReadObject(id)
	if err != nil {
		log.Err(err).Msg("Could not read object")
		_ = c.Error(err)
		return
	}
	c.JSON(200, object)
}
