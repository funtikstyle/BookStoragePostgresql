package clientHandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"../../domain"
	"github.com/gorilla/mux"
)

type clientHandler struct {
	service domain.ClientService
}

func NewClientHandler(clientService domain.ClientService) *clientHandler {
	return &clientHandler{service: clientService}
}

func (ch *clientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	newClient := domain.Client{}

	err := json.NewDecoder(r.Body).Decode(&newClient)
	log.Println(newClient)
	if err != nil {
		log.Panicln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok := ch.service.CreateClientService(newClient)
	if err != nil {
		log.Panicln(ok)
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(newClient)
	_, _ = w.Write(response)
}

func (ch *clientHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(ch.service.GetClientsService())
	w.Write(response)
}

func (ch *clientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)

	err := ch.service.DeleteClientService(id)
	if err != nil {
		if err.Error() == "internal server error" {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (ch *clientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)

	book, ok := ch.service.GetClientService(id)
	if !ok {
		http.Error(w, "Клиент отсутствует", http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(book)
	w.Write(response)
}

func (ch *clientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)

	newClient := domain.Client{}
	err := json.NewDecoder(r.Body).Decode(&newClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, error := ch.service.UpdateClientService(id, newClient)
	if !ok {
		http.Error(w, "Клиент отсутствует", http.StatusBadRequest)
		return
	}
	if error != nil {
		return  
	}
}
