package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(s *Server, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Upgrader error:%s \n", err)
	}
	c, err := s.NewClient(conn)
	if err != nil {
		log.Fatalf("Client init error:%s \n", err)
	}
	c.Start()
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Print(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./resource/index.html")
}

//return token then login through ws with token
//in production env, the token is generated by page
// this function must be depricated in production
func Login(s *Server, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := ps.ByName("account")
	session := s.GenSession(account)
	_, err := fmt.Fprint(w, session)
	//log.Printf("Login Session:%s, bytes:%d", session, n)
	if err != nil {
		log.Fatalf("Login error:%s \n", err)
	}
}
