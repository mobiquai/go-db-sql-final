package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int64, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client), sql.Named("status", p.Status), sql.Named("address", p.Address), sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId() // получаем последний добавленный идентификатор
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return id, nil
}

func (s ParcelStore) Get(number int64) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	onerow := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = :id", sql.Named("id", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := onerow.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}

		// Читаем данные из текущей строки, куда перемещён курсор
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}

		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int64, status string) error {
	// реализуйте обновление статуса в таблице parcel

	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :id", sql.Named("status", status), sql.Named("id", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int64, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :id AND status = 'registered'", sql.Named("address", address), sql.Named("id", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int64) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :id AND status = 'registered'", sql.Named("id", number))
	if err != nil {
		return err
	}

	return nil
}
