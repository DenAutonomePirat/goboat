package server

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type Web struct {
	db     *Store
	mux    *Mux
	tokens map[string]time.Time
}

func NewWeb(s *Store) *Web {
	return &Web{
		mux:    NewMux(),
		db:     s,
		tokens: make(map[string]time.Time),
	}
}

func (w *Web) makeCookie(c *gin.Context) {
	session := sessions.Default(c)
	temp := make([]byte, 16)
	rand.Read(temp)
	token := string(temp)
	w.tokens[token] = time.Now().Add(time.Hour * 72)
	session.Set("token", token)
	session.Save()
}

func (w *Web) checkCookieToken(c *gin.Context) bool {
	session := sessions.Default(c)
	token, ok := session.Get("token").(string)
	if !ok {
		return false
	}
	if time.Now().Before(w.tokens[token]) {
		w.tokens[token] = time.Now().Add(time.Hour * 72)
		return true
	}
	return false
}

func (w *Web) ListenAndServe(g *Configuration) {

	r := gin.Default()
	store := sessions.NewCookieStore([]byte("secret")) //todo get "secret" from conf
	r.Use(sessions.Sessions("redboat", store))

	r.Static("/assets", "./server/assets")
	r.Static("/css", "./server/assets/css")
	r.Static("/javascripts", "./server/assets/javascripts")
	r.Static("/images", "./server/assets/images")

	r.GET("/", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			c.Redirect(303, "/assets")
			return
		}
		c.Redirect(303, "/assets/login.html")
	})

	r.GET("/assets", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			c.Redirect(303, "/assets/index.html")
			return
		}
		c.Redirect(303, "/assets/login.html")
	})

	r.GET("/login", func(c *gin.Context) {
		u := NewUser()
		uname, _ := c.GetQuery("u")
		pwd, _ := c.GetQuery("p")
		u, _ = w.db.getUser(uname)

		if u.CheckPassword(pwd) {
			fmt.Println("match")
			w.makeCookie(c)
			c.Redirect(303, "/assets")
			return
		}
		c.Redirect(303, "/assets/login.html")
	})

	r.GET("/api/gamesetup", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			c.JSON(200, g)
		}
	})

	r.GET("/ws", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			w.mux.Handle(c.Writer, c.Request)
		}
	})
	r.Run()
}
