package server

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type store struct {
	conn *sql.DB
}

// Create a new store struct, initializing
// the database with three default sensors.
func newStore() (db *store, err error) {
	var conn *sql.DB
	if conn, err = sql.Open("sqlite3", ":memory:"); err != nil {
		return nil, err
	}

	createStatement := `
		CREATE TABLE IF NOT EXISTS sensors (
			name string PRIMARY KEY,
			latitude REAL,
			longitude REAL,
			unit string,
			ingress string,
			distiller string
		)`
	if _, err = conn.Exec(createStatement); err != nil {
		return nil, err
	}

	insertStatement := `
		INSERT INTO sensors(name, latitude, longitude, unit, ingress, distiller) 
		VALUES('L1MAG', 0, 0, 'volts', 'Middle of the Ocean', 'foo'),
			('L1ANG', 33.8, 117.9, 'deg', 'Anaheim', 'bar'),
			('C1MAG', 37.8, 175.7, 'amps', 'New Zealand', 'baz')`
	if _, err = conn.Exec(insertStatement); err != nil {
		return nil, err
	}

	return &store{conn: conn}, nil
}

// Query a particular sensor by name, returning it as a sensor struct.
func (db *store) querySensor(name string) (sensor *Sensor, err error) {
	var rows *sql.Rows
	selectStatement := `
		SELECT latitude, longitude, unit, ingress, distiller 
		FROM sensors 
		WHERE name=(?)`
	if rows, err = db.conn.Query(selectStatement, name); err != nil {
		return nil, err
	}

	var lat, lon float64
	var unit, ingress, distiller string
	for rows.Next() {
		rows.Scan(&lat, &lon, &unit, &ingress, &distiller)
		sensor = CreateSensor(name, unit, ingress, distiller, lat, lon)
	}

	if sensor == nil {
		return nil, errors.New("Sensor not found in store")
	}

	return sensor, nil
}

// Queries all sensors from the database and returns them as a slice.
func (db *store) queryAllSensors() (sensors []*Sensor, err error) {
	var rows *sql.Rows
	selectStatement := `
		SELECT name, latitude, longitude, unit, ingress, distiller 
		FROM sensors`
	if rows, err = db.conn.Query(selectStatement); err != nil {
		return nil, err
	}

	var lat, lon float64
	var unit, name, ingress, distiller string
	for rows.Next() {
		rows.Scan(&name, &lat, &lon, &unit, &ingress, &distiller)
		newSensor := CreateSensor(name, unit, ingress, distiller, lat, lon)
		sensors = append(sensors, newSensor)
	}

	return sensors, nil
}

// Insert sensor into the database.
func (db *store) insertSensor(name, ingress, distiller, unit string, lat, lon float64) (err error) {
	insertStatement := `
		INSERT INTO sensors (name, latitude, longitude, unit, ingress, distiller) 
		VALUES(?, ?, ?, ?, ?, ?)`
	if _, err = db.conn.Exec(insertStatement, name, lat, lon, unit, ingress, distiller); err != nil {
		return err
	}

	return nil
}

// Update a sensor already in the database.
func (db *store) updateSensor(name, ingress, distiller, unit string, lat, lon float64) (err error) {
	updateStatement := `
		UPDATE sensors 
		SET name=?, latitude=?, longitude=?, unit=?, ingress=?, distiller=? 
		WHERE name = ?`
	if _, err = db.conn.Exec(updateStatement, name, lat, lon, unit, ingress, distiller, name); err != nil {
		return err
	}
	return nil
}

// Helper function to create a sensor struct.
func CreateSensor(name, unit, ingress, distiller string, lat, lon float64) *Sensor {
	return &Sensor{
		Name: name,
		Location: Coordinates{
			Latitude:  lat,
			Longitude: lon,
		},
		Tags: SensorTags{
			Name:      name,
			Unit:      unit,
			Ingress:   ingress,
			Distiller: distiller,
		},
	}
}
