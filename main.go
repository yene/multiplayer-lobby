package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/speps/go-hashids"
)

var g Games
var h *hashids.HashID

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 7
	h = hashids.NewWithData(hd)

	g.lobbies = make(map[string]*Lobby)

	router := httprouter.New()
	router.GET("/create", createLobby)
	router.GET("/lobbies", listLobbies)
	router.GET("/lobby/:hashid", showLobby)
	router.GET("/join/:hashid", joinLobby)
	static := httprouter.New()
	static.ServeFiles("/*filepath", http.Dir("frontend/dist"))
	router.NotFound = static
	log.Println("http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func listLobbies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g.RLock()
	defer g.RUnlock()
	js, err := json.Marshal(g.lobbies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func showLobby(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g.Lock()
	defer g.Unlock()
	hashid := ps.ByName("hashid")
	lobby, ok := g.lobbies[hashid]
	if !ok {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	js, err := json.Marshal(lobby)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func joinLobby(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g.Lock()
	defer g.Unlock()
	hashid := ps.ByName("hashid")
	lobby, ok := g.lobbies[hashid]
	if !ok {
		http.Error(w, "Lobby not found", http.StatusNotFound)
		return
	}

	newPlayerID := len(lobby.Players) + 1
	lobby.Players = append(lobby.Players, newPlayerID)

	answer := JoinAnswer{lobby, newPlayerID}

	js, err := json.Marshal(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createLobby(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	g.Lock()
	defer g.Unlock()
	newLobbyID := len(g.lobbies) + 1
	e, _ := h.Encode([]int{newLobbyID})
	newLobby := &Lobby{newLobbyID, e, time.Now(), []int{}}
	g.lobbies[e] = newLobby

	js, err := json.Marshal(newLobby)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
