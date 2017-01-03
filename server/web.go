package server

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type Web struct {
	db   *Store
	mux  *Mux
	Crew map[string]*SkipperWatch
}

func NewWeb(s *Store) *Web {
	return &Web{
		mux:  NewMux(),
		db:   s,
		Crew: make(map[string]*SkipperWatch),
	}
}

func (w *Web) newCookie(c *gin.Context, u *User) {
	session := sessions.Default(c)
	watch := newSkipperWatch()
	watch.UserName = u.UserName
	watch.Expiry = time.Now().Add(time.Hour * 72)

	temp := make([]byte, 16)
	rand.Read(temp)
	token := string(temp)
	fmt.Println(token)

	w.Crew[token] = watch
	session.Set("token", token)
	session.Save()
}

func (w *Web) checkCookieToken(c *gin.Context) bool {
	session := sessions.Default(c)
	token, ok := session.Get("token").(string)
	if !ok {
		return false
	}

	if crew, ok := w.Crew[token]; ok {
		if time.Now().Before(crew.Expiry) {
			w.Crew[token].Expiry = time.Now().Add(time.Hour * 72)
			return true
		}
	}
	return false
}

func (w *Web) getUserFromCookieToken(c *gin.Context) string {
	session := sessions.Default(c)
	token, ok := session.Get("token").(string)
	if !ok {
		return ""
	}
	return w.Crew[token].UserName
}

func (w *Web) removeCookie(c *gin.Context) {
	session := sessions.Default(c)
	token, _ := session.Get("token").(string)
	fmt.Println(token)
	delete(w.Crew, token)
	fmt.Println(token)
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

	r.POST("/login", func(c *gin.Context) {

		if c.PostForm("action") == "logout" {
			w.removeCookie(c)
			c.Redirect(303, "/assets/login.html")
			return
		}

		u := NewUser()

		uname := c.PostForm("uname")
		pwd := c.PostForm("psw")
		fmt.Println(uname)

		u, _ = w.db.getUser(uname)

		if u.CheckPassword(pwd) {
			fmt.Println("match")
			w.newCookie(c, u)
			c.Redirect(303, "/assets")
			return
		}
		c.Redirect(303, "/assets/login.html")
	})

	r.GET("/login", func(c *gin.Context) {
		c.Redirect(303, "/assets/login.html")
	})

	r.GET("/spectate", func(c *gin.Context) {
		c.Redirect(303, "/assets/spectate.html")
	})

	r.GET("/api/gamesetup", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			c.JSON(200, g)
		}
	})

	r.GET("/api/whoami", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			c.JSON(200, w.getUserFromCookieToken(c))
		}
	})

	r.GET("/ws", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			userName := w.getUserFromCookieToken(c)
			w.mux.Handle(c.Writer, c.Request, userName)
		}
	})
	r.Run()
}
