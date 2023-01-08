package service

import (
	"errors"
	"log"

	"../domain"
)

type service struct {
	bookStorage   domain.BookStorage
	clientStorage domain.ClientStorage
}

func NewServiceBook(bookstorage domain.BookStorage, clientstorage domain.ClientStorage) *service {
	return &service{
		bookStorage:   bookstorage,
		clientStorage: clientstorage,
	}
}

func (s *service) TakeABookService(clientId int64, ids []int64) []int64 {
	takeBooks := s.bookStorage.GetNotTakenBookByIds(ids)
	var listId []int64

	for id, book := range takeBooks {
		log.Println(clientId)
		book.ClientID = clientId
		book.IsTaken = true
		s.bookStorage.UpdateBookStorage(id, book)

		listId = append(listId, id)
	}

	return listId

}

func (s *service) CreateBookService(book domain.Book) error {

	err := s.bookStorage.CreateBookStorage(book)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetBooksService() (map[int64]domain.Book, error) {
	return s.bookStorage.GetBooksStorage()
}

func (s *service) GetBookService(id int64) (domain.Book, error) {
	book, err := s.bookStorage.GetBookStorage(id)
	return book, err
}

func (s *service) ReturnABookService(id int64) error {
	book, err := s.bookStorage.GetBookStorage(id)
	if err != nil {
		return err
	}

	if !book.IsTaken {
		return err
	}

	book.ClientID = 0
	book.IsTaken = false
	s.bookStorage.UpdateBookStorage(id, book)

	return err
}

func (s *service) DeleteBookService(id int64) error {
	err := s.bookStorage.DeleteBookStorage(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateBookService(id int64, book domain.Book) bool {
	s.bookStorage.UpdateBookStorage(id, book)

	return true
}

func (s *service) CreateClientService(client domain.Client) error {
	err := s.clientStorage.CreateClientStorage(client)

	return err

}

func (s *service) GetClientsService() map[int64]domain.Client {

	return s.clientStorage.GetClientsStorage()
}

func (s *service) DeleteClientService(id int64) error {
	client, err := s.clientStorage.GetClientStorage(id)

	if err != nil {
		println(err.Error())
		return errors.New("internal server error")
	}

	if client.ID == 0 {
		return errors.New("клиент не найден")
	}

	haveBooks, err := s.bookStorage.StatusClientByBooks(id)
	if err != nil {
		println(err.Error())
		return errors.New("internal server error")
	}
	if haveBooks {
		return errors.New("клиент не вернул книги")
	}

	err = s.clientStorage.DeleteClientStorage(id)

	if err != nil {
		println(err.Error())
		return errors.New("internal server error")
	}

	return nil
}

func (s *service) GetClientService(id int64) (domain.Client, bool) {
	client, _ := s.clientStorage.GetClientStorage(id)
	
	if client.ID == 0 {
		return client, false
	}

	return client, true
}

func (s *service) UpdateClientService(id int64, client domain.Client) (bool, error) {
	clientOld, err := s.clientStorage.GetClientStorage(id)
	
	if clientOld.ID == 0 {
		return false, err
	}

	error := s.clientStorage.UpdateClientStorage(id, client)

	return true, error
}

func (s *service) GetBooksByClientId(id int64) (map[int64]domain.Book, error) {
	booklist, err := s.bookStorage.GetBooksByClientId(id)

	return booklist, err
}
