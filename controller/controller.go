package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/abserari/jx-co2-101-sensor/model/mysql"
	"github.com/gin-gonic/gin"
)

type dioxideDensity struct {
	db *sql.DB
}

func New(db *sql.DB) *dioxideDensity {
	return &dioxideDensity{
		db: db,
	}
}

func (d dioxideDensity) RegistRouter(r gin.IRouter) {
	r.POST("/dioxide", d.Add)
}

func (d dioxideDensity) Add(c *gin.Context) {
	var req struct {
		DioxideDensity int    `json: "dioxide" binding:"required"`
		DeviceId       string `json: "deviceId" binding:"required"`
		// Status         int `json: "status" binding:"required"`
		ZoneName string `json: "zoneName" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err := mysql.InsertDioxide(d.db, req.DioxideDensity, 0, req.ZoneName, req.DeviceId)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
