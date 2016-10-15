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

	/*r.Use(func(c *gin.Context) {
		s := sessions.Default(c)
		switch state := s.Get("state"); state {

		case nil:
			v := ""
			s.Set("state", v)
			s.Save()
			c.Redirect(301, "/assets/logon.html")
			return

		case "logged in":
			//chech random against map of active cookies
			c.Next()

		default:
			fmt.Println("Default fired")
		}
	})*/

	r.Static("/assets", "./server/assets")
	r.Static("/css", "./server/assets/css")
	r.Static("/javascripts", "./server/assets/javascripts")
	r.Static("/images", "./server/assets/images")

	r.GET("/logon", func(c *gin.Context) {
		c.Redirect(301, "/assets/logon.html")
	})
	r.GET("/api/gamesetup", func(c *gin.Context) {
		c.JSON(200, g)
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/assets/index.html")
	})
	r.GET("/ws", func(c *gin.Context) {
		w.mux.Handle(c.Writer, c.Request)
	})

	r.Run()
}
