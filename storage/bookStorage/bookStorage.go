package bookStorage

import (
	"context"
	_ "database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"../../domain"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgxpool"
)

type BookStorage struct {
	list    map[int64]domain.Book
	connect *pgxpool.Pool
}

func NewBookStorage(dbpool *pgxpool.Pool) *BookStorage {

	return &BookStorage{
		list:    map[int64]domain.Book{},
		connect: dbpool,
	}
}

func (pb BookStorage) GetBooksStorage() (map[int64]domain.Book, error) {
	rows, err := pb.connect.Query(context.Background(), "SELECT * FROM book")
	if err != nil {
		return nil, err
	}

	return rowsToBooksList(rows)
}

func (pb *BookStorage) GetBookStorage(id int64) (domain.Book, error) {
	// book, ok := pb.list[id]
	book := domain.Book{}
	// ok := false

	rows, err := pb.connect.Query(
		context.Background(),
		fmt.Sprintf(
			"SELECT id, Author, Title, ClientID, isTaken"+
				"FROM book AS b WHERE b.id = %d", id),
	)
	if err != nil {
		return book, err
	}

	if rows.Next() {
		val, err := rows.Values()
		if err != nil {
			return book, err
		}

		book = valuesToBook(val)
		// ok = true
	}

	return book, nil
}

func (pb *BookStorage) UpdateBookStorage(id int64, book domain.Book) error {
	// pb.list[id] = book

	_, err := pb.connect.Query(
		context.Background(),
		fmt.Sprintf(
			"UPDATE book "+
				"SET \"Author\" = '%s', "+
				"\"Title\" = '%s', "+
				"\"ClientID\" = %d "+
				"WHERE \"id\" = %d", book.Author, book.Title, book.ClientID, id),
	)
	if err != nil {
		return err
	}
	return nil
}

func (pb *BookStorage) CreateBookStorage(book domain.Book) error {
	_, err := pb.connect.Query(
		context.Background(),
		fmt.Sprintf(
			"INSERT INTO book (\"Author\", \"Title\")"+
				"VALUES ('%s', '%s')", book.Author, book.Title),
	)
	if err != nil {
		return err
	}
	return nil
}

func (pb *BookStorage) DeleteBookStorage(id int64) error {
	_, err := pb.connect.Query(
		context.Background(),
		fmt.Sprintf(
			"DELETE FROM book "+
				"WHERE \"id\" = %d", id),
	)
	if err != nil {
		return err

	}
	return nil
}

func (pb *BookStorage) StatusClientByBooks(id int64) (bool, error) {
	rows, err := pb.connect.Query(context.Background(),
		fmt.Sprintf("SELECT * FROM book WHERE \"ClientID\" = %d LIMIT 1", id))
	if err != nil {
		println(err.Error())
		return false, err
	}

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func (pb *BookStorage) GetBooksByClientId(id int64) (map[int64]domain.Book, error) {
	rows, err := pb.connect.Query(context.Background(),
		fmt.Sprintf("SELECT * FROM book WHERE \"ClientID\" = %d ", id))

	if err != nil {
		println(err.Error())
		return nil, err
	}

	books, err := rowsToBooksList(rows)

	if err != nil {
		println(err.Error())
		return nil, err
	}

	return books, nil
}

func (pb *BookStorage) GetNotTakenBookByIds(ids []int64) map[int64]domain.Book {
	idsString := make([]string, len(ids))
	for k, v := range ids {
		idsString[k] = strconv.FormatInt(v, 10)
	}

	rows, _ := pb.connect.Query(context.Background(),
		fmt.Sprintf("SELECT * FROM book WHERE \"id\" IN (%s) AND \"isTaken\" = false", strings.Join(idsString, ",")))
	books, _ := rowsToBooksList(rows)

	return books
}

func rowsToBooksList(rows pgx.Rows) (map[int64]domain.Book, error) {
	booklist := map[int64]domain.Book{}

	for rows.Next() {
		val, err := rows.Values()
		if err != nil {
			return nil, errors.New("error while iterating dataset")
		}

		book := valuesToBook(val)
		booklist[int64(book.ID)] = book

	}

	return booklist, nil
}

func valuesToBook(val []any) domain.Book {
	id := val[0].(int32)
	Author := val[1].(string)
	Title := val[2].(string)
	var ClientID int64

	if val[3] != nil {
		ClientID = int64(val[3].(int32))
	}

	IsTaken := val[4].(bool)

	book := domain.Book{
		ID:       id,
		Author:   Author,
		Title:    Title,
		ClientID: ClientID,
		IsTaken:  IsTaken,
	}
	return book
}
