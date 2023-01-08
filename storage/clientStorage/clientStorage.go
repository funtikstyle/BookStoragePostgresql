package clientStorage

import (
	"context"
	"fmt"

	"../../domain"
	"github.com/jackc/pgx/pgxpool"
)

type ClientStorage struct {
	list    map[int64]domain.Client
	connect *pgxpool.Pool
}

func NewClientStorage(dbpool *pgxpool.Pool) *ClientStorage {
	return &ClientStorage{
		list:    map[int64]domain.Client{},
		connect: dbpool,
	}
}

func (cs *ClientStorage) CreateClientStorage(client domain.Client) error {

	_, err := cs.connect.Query(
		context.Background(),
		fmt.Sprintf(
			"INSERT INTO client (\"Name\", \"Phone\")"+
				"VALUES ('%s', '%s')", client.Name, client.PhoneName),
	)
	if err != nil {
		return err
	}
	return nil
}

func (cs *ClientStorage) GetClientsStorage() map[int64]domain.Client {
	return cs.list
}

func (cs *ClientStorage) DeleteClientStorage(id int64) error {
	_, err := cs.connect.Query(
		context.Background(),
		fmt.Sprintf(
			"DELETE FROM client "+
				"WHERE \"id\" = %d", id),
	)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ClientStorage) GetClientStorage(id int64) (domain.Client, error) {
	client := domain.Client{}

	rows, err := cs.connect.Query(context.Background(),
		fmt.Sprintf("SELECT \"id\", \"Name\", \"Phone\" FROM clientWHERE\"id\" = %d", id))
	if err != nil {
		return client, err
	}

	if rows.Next() {
		val, err := rows.Values()
		if err != nil {
			return client, err
		}
		client = valuesToClient(val)
	}

	return client, nil
}

func (cs *ClientStorage) UpdateClientStorage(id int64, client domain.Client) error {

	_, err := cs.connect.Query(
		context.Background(),
		fmt.Sprintf(
			"UPDATE book "+
				"SET \"Name\" = '%s', "+
				"\"Phone\" = '%s', "+
				"WHERE \"id\" = %d", client.Name, client.PhoneName, id),
	)
	if err != nil {
		return err
	}
	return nil

	// cs.list[id] = client
}

func valuesToClient(val []any) domain.Client {
	id := val[0].(int32)
	Name := val[1].(string)
	Phone := val[2].(int32)

	client := domain.Client{
		ID:        id,
		Name:      Name,
		PhoneName: string(Phone),
	}
	return client
}
