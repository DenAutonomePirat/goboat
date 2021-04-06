package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type session struct {
	UserID         string
	Authenticated  bool
	ConnectedAt    time.Time     `json:"connectedAt" bson:"connectedAt"`
	OnlineDuration time.Duration `json:"onlineDuration"`
}

type Web struct {
	db         *Store
	mux        *Mux
	router     *mux.Router
	conf       *Configuration
	sessions   map[string]*session
	skipperReg chan string
}

func NewWeb(s *Store) *Web {
	return &Web{
		db:         s,
		mux:        NewMux(),
		router:     mux.NewRouter(),
		sessions:   make(map[string]*session),
		skipperReg: make(chan string),
	}
}

func (web *Web) ListenAndServe(g *Configuration) {
	web.conf = g
	web.router.HandleFunc("/", web.index)
	web.router.HandleFunc("/login", web.login).Methods("get")
	web.router.HandleFunc("/login", web.auth).Methods("post")
	web.router.HandleFunc("/logout", web.logout)
	web.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		web.serveWS(w, r)
	})
	web.router.HandleFunc("/api/gamesetup", web.api)
	web.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./server/assets")))
	http.ListenAndServe(":80", web.router)
}

func (web *Web) api(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, web.conf)
}

func (web *Web) index(w http.ResponseWriter, r *http.Request) {
	if web.isAuthenticated(w, r) {
		fmt.Println("Serving da shit")
		http.ServeFile(w, r, "./server/assets/index.html")
	} else {
		http.Redirect(w, r, "/login", 303)
	}
}

func (web *Web) logout(w http.ResponseWriter, r *http.Request) {
	ck, err := r.Cookie("redboat")
	if err != nil {
		http.Redirect(w, r, "/login", 303)
	}
	s, ok := web.sessions[ck.Value]
	if ok {
		fmt.Printf("deleting session for %s\n", s.UserID)
		web.skipperReg <- s.UserID
		delete(web.sessions, ck.Value)
	}

	http.Redirect(w, r, "/login", 303)
}

func (web *Web) auth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := NewUser()
	uname := r.Form.Get("u")
	pwd := r.Form.Get("p")
	u, err := web.db.getUser(uname)
	if err != nil {
		fmt.Printf("creating user %s", uname)
		u.UserName = uname
		u.SetPassword(pwd)
		web.db.AddUser(u)
		time.Sleep(time.Second)
	}

	if u.CheckPassword(pwd) {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s := &session{
			UserID:        u.UserName,
			Authenticated: true,
			ConnectedAt:   time.Now(),
		}
		random := RandomString(16)
		web.sessions[random] = s
		ck := http.Cookie{
			Name:     "redboat",
			Value:    random,
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

func (web *Web) isAuthenticated(w http.ResponseWriter, r *http.Request) bool {

	ck, err := r.Cookie("redboat")
	if err != nil {
		return false
	}
	s := web.sessions[ck.Value]
	if s == nil {
		return false
	}
	return s.Authenticated

}

func (web *Web) getSession(w http.ResponseWriter, r *http.Request) *session {

	ck, err := r.Cookie("redboat")
	if err != nil {
		return &session{
			UserID:         "",
			Authenticated:  false,
			ConnectedAt:    time.Time{},
			OnlineDuration: 0,
		}
	}
	s := web.sessions[ck.Value]
	return s
}

// serveWs handles websocket requests from the peer.
func (web *Web) serveWS(w http.ResponseWriter, r *http.Request) {
	if web.isAuthenticated(w, r) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Could not upgrade http request: %s", err.Error())
			return
		}
		conn := NewConn(web.mux, ws)
		conn.user = web.getSession(w, r).UserID
		web.mux.register <- conn
		return
	}
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
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
