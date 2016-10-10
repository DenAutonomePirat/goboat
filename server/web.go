package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Web struct {
	mux *Mux
}

func NewWeb() *Web {
	return &Web{mux: NewMux()}
}

func (w *Web) ListenAndServe() {

	r := gin.Default()

	r.Static("/assets", "./server/assets")
	r.Static("/css", "./server/assets/css")
	r.Static("/javascripts", "./server/assets/javascripts")
	r.Static("/images", "./server/assets/images")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/assets/index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		w.mux.Handle(c.Writer, c.Request)
	})

	r.Run()
}
