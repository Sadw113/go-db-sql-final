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

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	// верните идентификатор последней добавленной записи

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		parsel := Parcel{}

		err := rows.Scan(&parsel.Number, &parsel.Client, &parsel.Status, &parsel.Address, &parsel.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, parsel)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("update parcel set status = :status where number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	// var status string

	// row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))

	// err := row.Scan(&status)

	// if err != nil {
	// 	return err
	// }

	// if status != ParcelStatusRegistered {
	// 	return fmt.Errorf("Статус не зарегистрирован")
	// }

	_, err := s.db.Exec("update parcel set address = :address where number = :number and status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number and status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		return err
	}

	return nil
}
