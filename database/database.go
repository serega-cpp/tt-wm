package database

import (
	"database/sql"
	"fmt"

	"finleap/api"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func Create(host string, port int, user, pass, dbname string) (*DB, error) {
	cs := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		user, pass, dbname, host, port)

	db, err := sql.Open("postgres", cs)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (d *DB) Destroy() error {
	return d.db.Close()
}

///////////////////////////////////////////////////////////
// TABLES Creators

func (d *DB) CreateCitiesTable() error {
	stmt := `CREATE TABLE IF NOT EXISTS cities (
		id        SERIAL,
		name      VARCHAR (128) NOT NULL,
		latitude  NUMERIC(9, 6),
		longitude NUMERIC(9, 6))`
	if _, err := d.db.Exec(stmt); err != nil {
		return err
	}
	return nil
}

func (d *DB) CreateTemperaturesTable() error {
	stmt := `CREATE TABLE IF NOT EXISTS temperatures (
		id        SERIAL,
		city_id   INTEGER,
		max_cels  REAL,
		min_cels  REAL,
		time_unix BIGINT)`
	if _, err := d.db.Exec(stmt); err != nil {
		return err
	}
	return nil
}

func (d *DB) CreateWebhooksTable() error {
	stmt := `CREATE TABLE IF NOT EXISTS webhooks (
		id           SERIAL,
		city_id      INTEGER,
		callback_url VARCHAR (1024) NOT NULL)`
	if _, err := d.db.Exec(stmt); err != nil {
		return err
	}
	return nil
}

///////////////////////////////////////////////////////////
// TABLES Cleaners

func (d *DB) ClearCitiesTable() error {
	if _, err := d.db.Exec("DELETE FROM cities"); err != nil {
		return err
	}
	if _, err := d.db.Exec("ALTER SEQUENCE cities_id_seq RESTART"); err != nil {
		return err
	}
	return nil
}

func (d *DB) ClearTemperaturesTable() error {
	if _, err := d.db.Exec("DELETE FROM temperatures"); err != nil {
		return err
	}
	if _, err := d.db.Exec("ALTER SEQUENCE temperatures_id_seq RESTART"); err != nil {
		return err
	}
	return nil
}

func (d *DB) ClearWebhooksTable() error {
	if _, err := d.db.Exec("DELETE FROM webhooks"); err != nil {
		return err
	}
	if _, err := d.db.Exec("ALTER SEQUENCE webhooks_id_seq RESTART"); err != nil {
		return err
	}
	return nil
}

///////////////////////////////////////////////////////////
// Queries

func (d *DB) InsertCity(c *api.City) error {
	stmt := "INSERT INTO cities (name, latitude, longitude) VALUES ($1, $2, $3) RETURNING id"
	return d.db.QueryRow(stmt, c.Name, c.Latitude, c.Longitude).Scan(&c.ID)
}

func (d *DB) UpdateCity(c *api.City) error {
	stmt := "UPDATE cities SET name=$1, latitude=$2, longitude=$3 WHERE id=$4"
	res, err := d.db.Exec(stmt, c.Name, c.Latitude, c.Longitude, c.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("City with id=%d not exists", c.ID)
	}
	return nil
}

func (d *DB) DeleteCity(c *api.City) error {
	if err := d.SelectOnlyCity(c); err != nil {
		return err
	}

	stmt := "DELETE FROM cities WHERE id=$1"
	_, err := d.db.Exec(stmt, c.ID)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) SelectOnlyCity(c *api.City) error {
	stmt := "SELECT Name, Latitude, Longitude FROM cities WHERE id=$1"
	row := d.db.QueryRow(stmt, c.ID)
	return row.Scan(&c.Name, &c.Latitude, &c.Longitude)
}
