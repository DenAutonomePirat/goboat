package server

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Web struct {
	mux *Mux
}

func NewWeb() *Web {
	return &Web{mux: NewMux()}
}

func (w *Web) ListenAndServe(g *GameSetup) {

	r := gin.Default()
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("redboat", store))

	r.Static("/assets", "./server/assets")
	r.Static("/css", "./server/assets/css")
	r.Static("/javascripts", "./server/assets/javascripts")
	r.Static("/images", "./server/assets/images")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/assets/index.html")
	})

	r.GET("/logon", func(c *gin.Context) {
		c.Redirect(301, "/assets/logon.html")
	})

	r.GET("/api/gamesetup", func(c *gin.Context) {
		c.JSON(200, g)
	})

	r.GET("/ws", func(c *gin.Context) {
		w.mux.Handle(c.Writer, c.Request)
	})

	r.Run()
}
