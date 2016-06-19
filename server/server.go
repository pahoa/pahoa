package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pahoa/pahoa/core"
)

type Server struct {
	mux      *mux.Router
	board    *core.Board
	model    core.Model
	executor *core.Executor
}

type NewServerOptions struct {
	Board    *core.Board
	Model    core.Model
	Executor *core.Executor
}

func NewServer(opts *NewServerOptions) *Server {
	serveMux := mux.NewRouter()

	server := &Server{
		mux:      serveMux,
		board:    opts.Board,
		model:    opts.Model,
		executor: opts.Executor,
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

	card, err := core.AddCard(s.board, s.model, s.executor, &options)
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
	step := r.URL.Query().Get("step")
	cards, err := s.model.ListCards(step)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cards); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) updateCardStepHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
