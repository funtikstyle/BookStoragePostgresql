package bookHandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"../../domain"
	"github.com/gorilla/mux"
)

type bookHandler struct {
	service domain.BookService
}

func NewBookHandler(bookService domain.BookService) *bookHandler {
	return &bookHandler{service: bookService}
}

func (bh *bookHandler) TakeABook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["clientId"], 10, 0)

	var dat map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dat)

	if err != nil {
		panic(err)
	}
	strs := dat["ids"].([]interface{})

	ids := []int64{}

	for _, val := range strs {
		ids = append(ids, int64(val.(float64)))
	}

	response, _ := json.Marshal(bh.service.TakeABookService(id, ids))
	w.Write(response)
}

func (bh *bookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	newBook := domain.Book{}
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		log.Panicln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok:= bh.service.CreateBookService(newBook)
	if ok != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(newBook.Author)
	_, _ = w.Write(response)
}

func (bh *bookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := bh.service.GetBooksService()
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(books)
		if err !=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	w.Write(response)
}

func (bh *bookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)

	book, err := bh.service.GetBookService(id)
	if err != nil {
		http.Error(w, "Книга отсутствует", http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(book)
	w.Write(response)
}

func (bh *bookHandler) ReturnABook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)
	err := bh.service.ReturnABookService(id)
	if err != nil {
		http.Error(w, "Книга отсутствует", http.StatusBadRequest)
		return
	}
}

func (bh *bookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)

	err := bh.service.DeleteBookService(id)
	if err !=nil {
		http.Error(w, "Книга отсутствует", http.StatusBadRequest)
		return
	}
}

func (bh *bookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)

	newBook := domain.Book{}
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		log.Panicln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok := bh.service.UpdateBookService(id, newBook)
	if !ok {
		http.Error(w, "Книга отсутствует", http.StatusBadRequest)
		return
	}
}

func (bh *bookHandler) GetBooksByClientId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)

	booklist, err := bh.service.GetBooksByClientId(id)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	
	response, _ := json.Marshal(booklist)

	w.Write(response)
}
