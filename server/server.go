package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pahoa/pahoa/core"
)

var actions map[string]core.Action = map[string]core.Action{
	"MergeRequestToDevelop": func() {},
}

type Server struct {
	mux   *mux.Router
	board *core.Board
	model *core.Model
}

func NewServer() *Server {
	serveMux := mux.NewRouter()

	board := core.NewBoard()
	board.AddTransition("", "todo")
	board.AddTransition("todo", "in-development")
	board.AddTransition("in-development", "waiting-for-code-review",
		actions["MergeRequestToDevelop"])

	server := &Server{
		mux:   serveMux,
		board: board,
		model: &core.Model{},
	}

	// add card
	serveMux.HandleFunc("/cards", server.addCardHandler).
		Methods("POST")
	// card list
	serveMux.HandleFunc("/cards", server.listCardsHandler).
		Methods("GET")
	// move card
	serveMux.HandleFunc("/cards/{id}/step", server.updateCardStepHandler).
		Methods("POST")

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) addCardHandler(w http.ResponseWriter, r *http.Request) {
	var options core.AddCardOptions
	if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
		http.Error(w, "Unable to decode body", http.StatusBadRequest)
		return
	}

	card, err := core.AddCard(s.board, s.model, &options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := json.Marshal(card)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}

func (s *Server) listCardsHandler(w http.ResponseWriter, r *http.Request) {
	cards := s.model.ListCards()

	if err := json.NewEncoder(w).Encode(cards); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) updateCardStepHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
