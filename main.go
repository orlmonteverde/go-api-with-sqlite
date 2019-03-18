package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var note Note

// Punto de ejecución del ejecutable.
func main() {
	// Instancia de http.DefaultServeMux
	mux := http.NewServeMux()

	// flag para realizar la creación de las tablas en la base de datos.
	migrate := flag.Bool("migrate", false, "Crea las tablas en la base de datos")
	flag.Parse()

	if *migrate {
		if err := MakeMigrations(); err != nil {
			log.Fatal(err)
		}
	}

	// Rutas a manejar
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/notes", NotesHandler)

	// Log informativo
	log.Println("Corriendo en http://localhost:8080")

	// Servidor escuchando en el puerto 8080
	http.ListenAndServe(":8080", mux)
}

// IndexHandler nos permite manejar la petición a la ruta '/' y retornar "hola mundo"
// como respuesta al cliente.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hola mundo")
}

// GetNotesHandler nos permite manejar las peticiones a la ruta '/notes' con el método GET.
func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Puntero a una estructura de tipo Note
	n := new(Note)
	// Solicitando todas las notas en la base de datos
	notes, err := n.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// Conviertiendo el slice de Note a formato JSON, retorna un []byte
	j, err := json.Marshal(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// Escribiendo el código de respuesta.
	w.WriteHeader(http.StatusOK)
	// Estableciendo el tipo de contenido del cuerpo de la respuesta.
	w.Header().Set("Content-Type", "application/json")
	// Escribiendo la respuesta, es decir nuestro slice de notas en formato JSON.
	w.Write(j)
}

// CreateNotesHandler nos permite manejar las peticiones a la ruta '/notes' con el método POST.
func CreateNotesHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	// Tomando el cuerpo de la petición, en formato JSON, y decodificandola en
	// la variable note que acabamos de declarar.
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// Creamos la nueva nota gracias al método Create.
	err = note.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateNotesHandler nos permite manejar las peticiones a la ruta '/notes' con el método UPDATE.
func UpdateNotesHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// Actualizamos la nota correspondiente.
	err = note.Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteNotesHandler nos permite manejar las peticiones a la ruta '/notes' con el método DELETE.
func DeleteNotesHandler(w http.ResponseWriter, r *http.Request) {
	// obtendemosel valor pasado en la url como query correspondiente a id, del tipo ?id=3.
	idStr := r.URL.Query().Get("id")
	// Verificamos que no esté vacío.
	if idStr == "" {
		http.Error(w, "Query id es requerido", http.StatusBadRequest)
	}
	// Convertimos el valor obtenido del query a un int, de ser posible.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Query id debe ser un número", http.StatusBadRequest)
	}

	var note Note
	// Borramos la nota con el id correspondiente.
	err = note.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// NotesHandler nos permite manejar la petición a la ruta '/notes' y pasa el control al
// la función correspondiente según el método de la petición.
func NotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetNotesHandler(w, r)
	case http.MethodPost:
		CreateNotesHandler(w, r)
	case http.MethodPut:
		UpdateNotesHandler(w, r)
	case http.MethodDelete:
		DeleteNotesHandler(w, r)
	default:
		// Caso por defecto en caso de que se realice una petición con un
		// método deferente a los esperados.
		http.Error(w, "Metodo no permitido", http.StatusBadRequest)
	}
}
