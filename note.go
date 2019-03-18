package main

import (
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (n Note) Create() error {
	// Realizamos la conexión a la base de datos.
	db := GetConnection()

	// Query para insertar los datos en la tabla notes
	q := `INSERT INTO notes (title, description, updated_at)
			VALUES(?, ?, ?)`

	// Preparamos la petición para insertar los datos de manera segura
	// y evitar código malicioso.
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	// Ejecutamos la petición pasando los datos correspondientes. El orden
	// es importante, corresponde con los "?" delstring q.
	r, err := stmt.Exec(n.Title, n.Description, time.Now())
	if err != nil {
		return err
	}
	// Confirmamos que una fila fuera afectada, debido a que insertamos un
	// registro en la tabla. En caso contrario devolvemos un error.
	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("ERROR: Se esperaba una fila afectada")
	}
	// Si llegamos a este punto consideramos que todo el proceso fue exitoso
	// y retornamos un nil para confirmar que no existe un error.
	return nil
}

func (n *Note) GetAll() ([]Note, error) {
	db := GetConnection()
	q := `SELECT
			id, title, description, created_at, updated_at
			FROM notes`
	// Ejecutamos la query
	rows, err := db.Query(q)
	if err != nil {
		return []Note{}, err
	}
	// Cerramos el recurso
	defer rows.Close()

	// Declaramos un slice de notas para que almacene las notas que retorne
	// la petición.
	notes := []Note{}
	// El método Next retorna un bool, mientras sea true indicará que existe
	// un valor siguiente para leer.
	for rows.Next() {
		// Escaneamos el valor actual de la fila e insertamos el retorno
		// en los correspondientes campos de la nota.
		rows.Scan(&n.ID, &n.Title, &n.Description, &n.CreatedAt, &n.UpdatedAt)
		// Añadimos cada nueva nota al slice de notas que declaramos antes.
		notes = append(notes, *n)
	}
	return notes, nil
}

func (n *Note) GetByID(id int) (Note, error) {
	db := GetConnection()
	q := `SELECT
		id, title, description, created_at, updated_at
		FROM notes WHERE id=?`

	err := db.QueryRow(q, id).Scan(
		&n.ID, &n.Title, &n.Description, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return Note{}, err
	}

	return *n, nil
}

func (n Note) Update() error {
	db := GetConnection()
	q := `UPDATE notes set title=?, description=?, updated_at=?
		WHERE id=?`
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(n.Title, n.Description, time.Now(), n.ID)
	if err != nil {
		return err
	}
	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("ERROR: Se esperaba una fila afectada")
	}
	return nil
}

func (n Note) Delete(id int) error {
	db := GetConnection()

	q := `DELETE FROM notes
		WHERE id=?`
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("ERROR: Se esperaba una fila afectada")
	}
	return nil
}
