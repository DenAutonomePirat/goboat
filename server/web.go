package server

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

type Web struct {
	db       *Store
	mux      *Mux
	router   *mux.Router
	conf     *Configuration
	skippers map[string]*Skipper
}

func NewWeb(s *Store) *Web {
	return &Web{
		db:       s,
		mux:      NewMux(),
		router:   mux.NewRouter(),
		skippers: make(map[string]*Skipper),
	}
}

func (web *Web) ListenAndServe(g *Configuration) {
	web.conf = g
	web.router.HandleFunc("/", index)
	web.router.HandleFunc("/login", web.login).Methods("get")
	web.router.HandleFunc("/login", web.auth).Methods("post")
	web.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		web.serveWS(w, r)
	})
	web.router.HandleFunc("/api/gamesetup", web.api)
	web.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./server/assets")))
	http.ListenAndServe(":8080", web.router)
}

func (web *Web) api(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, web.conf)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./server/assets/index.html")
}

func (web *Web) auth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := NewUser()
	uname := r.Form.Get("u")
	pwd := r.Form.Get("p")
	u, err := web.db.getUser(uname)
	if err != nil {
		fmt.Println("user not found")
	}

	if u.CheckPassword(pwd) {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s := &Skipper{
			UserID:        u.UserName,
			Authenticated: true,
			ConnectedAt:   time.Now(),
		}
		temp := make([]byte, 16)
		rand.Read(temp)
		web.skippers[string(temp)] = s
		ck := http.Cookie{
			Name:     "redboat",
			Value:    "gnarls",
			Path:     "/",
			Domain:   "",
			Expires:  time.Now().Add(time.Hour * 3),
			Secure:   false,
			HttpOnly: true,
			SameSite: 0,
		}

		http.SetCookie(w, &ck)

		fmt.Printf("User %s Signed in\n", u.UserName)
		http.Redirect(w, r, "/", http.StatusFound)

	}
}
func (web *Web) login(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./server/assets/login.html")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func (web *Web) serveWS(w http.ResponseWriter, r *http.Request) {
	ck, _ := r.Cookie("redboat")
	spew.Dump(ck.Value)

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Could not upgrade http request: %s", err.Error())
		return
	}

	conn := NewConn(web.mux, ws)
	web.mux.register <- conn
}

// getUser returns a user from session s
// on error returns an empty user
func getSkipper(s *sessions.Session) *Skipper {
	val := s.Values["skipper"]
	skipper, ok := val.(*Skipper)
	if !ok {
		fmt.Println("not found in cookiejar")
		return &Skipper{Authenticated: false}
	}
	fmt.Println("found in cookiejar")

	return skipper
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

/*
	router.HandleFunc("/", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			c.Redirect(303, "/assets")
			return
		}
		c.Redirect(303, "/assets/login.html")
	})

	router.GET("/assets", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			c.Redirect(303, "/assets/index.html")
			return
		}
		c.Redirect(303, "/assets/login.html")
	})

	router.GET("/login", func(c *gin.Context) {
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

	router.GET("/api/gamesetup", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			c.JSON(200, g)
		}
	})

	router.GET("/ws", func(c *gin.Context) {
		if w.checkCookieToken(c) {
			session := sessions.Default(c)
			token, ok := session.Get("token").(string)
			if !ok {
			}
			spew.Dump(session)
			spew.Dump(token)
			w.mux.Handle(c.Writer, c.Request)
		}
	})
	router.Run()
}
*/
