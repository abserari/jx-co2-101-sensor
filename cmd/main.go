package main

import (
	"database/sql"

	"github.com/abserari/jx-co2-101-sensor/controller"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3307)/metabase")
	if err != nil {
		panic(err)
	}

	d := controller.New(db)
	d.RegistRouter(r.Group("/api/v1"))

	r.Run(":9000")
}
