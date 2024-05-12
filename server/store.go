package server

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type store struct {
	conn *sql.DB
}

func initDB() (db *store, err error) {
	var conn *sql.DB
	if conn, err = sql.Open("sqlite3", ":memory:"); err != nil {
		return nil, err
	}
	if _, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS sensors (
			name string PRIMARY KEY,
			latitude REAL,
			longitude REAL,
			unit string
		)
	`); err != nil {
		return nil, err
	}

	if err = storeInitialSensors(conn); err != nil {
		return nil, err
	}

	return &store{conn: conn}, nil
}

func storeInitialSensors(conn *sql.DB) (err error) {
	_, err = conn.Exec(`
		INSERT INTO sensors(name, latitude, longitude, unit) 
		VALUES('L1MAG', 0, 0, 'volts'),
		      ('L1ANG', 33.8, 117.9, 'deg'),
			  ('C1MAG', 37.8, 175.7, 'amps')
	`)
	return err
}

func (db *store) querySensor(name string) (sensor *Sensor, err error) {
	var rows *sql.Rows
	if rows, err = db.conn.Query("SELECT * FROM sensors WHERE name=(?)", name); err != nil {
		return nil, err
	}

	var unit string
	var lat, lon float64
	for rows.Next() {
		rows.Scan(&name, &lat, &lon, &unit)
		sensor = CreateSensor(name, unit, lat, lon)
	}

	if sensor == nil {
		return nil, errors.New("Sensor not found in store")
	}

	return sensor, nil
}

func (db *store) queryAllSensors() (sensors []*Sensor, err error) {
	var rows *sql.Rows
	if rows, err = db.conn.Query("SELECT * FROM sensors"); err != nil {
		return nil, err
	}

	var lat, lon float64
	var name, unit string
	for rows.Next() {
		rows.Scan(&name, &lat, &lon, &unit)
		newSensor := CreateSensor(name, unit, lat, lon)
		sensors = append(sensors, newSensor)
	}
	return sensors, nil
}

func (db *store) insertSensor(name, unit string, lat, lon float64) (err error) {
	insertStatement := "INSERT INTO sensors (name, latitude, longitude, unit) VALUES(?, ?, ?, ?)"
	if _, err = db.conn.Exec(insertStatement, name, lat, lon, unit); err != nil {
		return err
	}
	return nil
}

func (db *store) updateSensor(name, unit string, lat, lon float64) (err error) {
	updateStatement := "UPDATE sensors SET name=?, latitude=?, longitude=?, unit=? WHERE name = ?"
	if _, err = db.conn.Exec(updateStatement, name, lat, lon, unit, name); err != nil {
		return err
	}
	return nil
}

func CreateSensor(name, unit string, lat, lon float64) *Sensor {
	return &Sensor{
		Name: name,
		Location: Coordinates{
			Latitude:  lat,
			Longitude: lon,
		},
		Tags: SensorTags{
			Name: name,
			Unit: unit,
		},
	}
}
