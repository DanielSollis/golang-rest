package server

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() (conn *sql.DB, err error) {
	if conn, err = sql.Open("sqlite3", ":memory:"); err != nil {
		return nil, err
	}
	if _, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS sensors (
			name string PRIMARY KEY,
			latitude REAL,
			longitude REAL,
			unit string,
			Ingress string,
			distiller string
		)
	`); err != nil {
		return nil, err
	}

	if err = storeInitialSensors(conn); err != nil {
		return nil, err
	}

	return conn, nil
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

func (s *Server) querySensor(name string) (_ []Sensor, err error) {
	return nil, nil
}

func (s *Server) queryAllSensors() (sensors []*Sensor, err error) {
	var rows *sql.Rows
	if rows, err = s.db.Query("SELECT * FROM sensors"); err != nil {
		return nil, err
	}

	var lat, lon float64
	var name, unit string
	for rows.Next() {
		rows.Scan(&name, &lat, &lon, &unit)
		sensors = append(sensors, CreateSensor(name, unit, lat, lon))
	}
	return sensors, nil
}

const insertString = "INSERT INTO sensors(name, latitude, longitude, unit, Ingress, distiller) VALUES(?, ?, ?, ?, ?, ?)"

func (s *Server) insertSensor(name, unit string, lat, lon float64) (err error) {
	return nil
}
