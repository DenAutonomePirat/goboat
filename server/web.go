package server

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Web struct {
	db  *Store
	mux *Mux
}

func NewWeb(s *Store) *Web {
	return &Web{
		mux: NewMux(),
		db:  s,
	}
}

func (w *Web) ListenAndServe(g *Configuration) {

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

	r.GET("/login", func(c *gin.Context) {
		u := NewUser()
		uname, _ := c.GetQuery("u")
		pwd, _ := c.GetQuery("p")
		_, u = w.db.getUser(uname)
		pwdcheck := u.CheckPassword(pwd)
		if pwdcheck {
			fmt.Println("match")
			c.Redirect(301, "/assets/index.html")
			return
		}
		c.Redirect(301, "/assets/login2.html")
	})

	r.GET("/api/gamesetup", func(c *gin.Context) {
		c.JSON(200, g)
	})

	r.GET("/ws", func(c *gin.Context) {
		w.mux.Handle(c.Writer, c.Request)
	})

	r.Run()
}
