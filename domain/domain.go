package domain

type Book struct {
	ID       int32
	Author   string
	Title    string
	ClientID int64
	IsTaken  bool
}

type Client struct {
	ID        int32
	Name      string
	PhoneName string
}

type BookStorage interface {
	GetBookStorage(id int64) (Book, error)
	UpdateBookStorage(id int64, book Book) error
	CreateBookStorage(book Book) error
	GetBooksStorage() (map[int64]Book, error)
	DeleteBookStorage(id int64) error
	StatusClientByBooks(id int64) (bool, error)
	GetBooksByClientId(id int64) (map[int64]Book, error)
	GetNotTakenBookByIds(ids []int64) map[int64]Book
}

type BookService interface {
	GetBookService(id int64) (Book, error)
	TakeABookService(id int64, ids []int64) []int64
	CreateBookService(book Book) error
	GetBooksService() (map[int64]Book, error)
	ReturnABookService(id int64) error
	DeleteBookService(id int64) error
	UpdateBookService(id int64, book Book) bool
	GetBooksByClientId(id int64) (map[int64]Book, error)
}

type ClientStorage interface {
	CreateClientStorage(client Client)error
	GetClientsStorage() map[int64]Client
	DeleteClientStorage(id int64) error
	GetClientStorage(id int64) (Client, error)
	UpdateClientStorage(id int64, client Client)error
}

type ClientService interface {
	CreateClientService(client Client)error
	GetClientsService() map[int64]Client
	DeleteClientService(id int64) error
	GetClientService(id int64) (Client, bool)
	UpdateClientService(id int64, client Client) (bool,error)
}
