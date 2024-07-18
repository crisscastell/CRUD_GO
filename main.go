package main

import (
	//"log"
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	driver := "mysql"
	usuario := "root"
	contrasenia := ""
	nombre := "sistema"

	conexion, err := sql.Open(driver, usuario+":"+contrasenia+"@tcp(127.0.0.1)/"+nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", inicio)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insertar)

	fmt.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}

type Trabajador struct {
	id       int
	nombre   string
	apellido string
	correo   string
}

func inicio(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM trabajadores")

	if err != nil {
		panic(err.Error())
	}

	trabajador := Trabajador{}
	arregloTrabajador := []Trabajador{}

	for registros.Next() {
		var id int
		var nombre, apellido, correo string
		err = registros.Scan(&id, &nombre, &apellido, &correo)

		if err != nil {
			panic(err.Error())
		}

		trabajador.id = id
		trabajador.nombre = nombre
		trabajador.apellido = apellido
		trabajador.correo = correo

		arregloTrabajador = append(arregloTrabajador, trabajador)
	}
	fmt.Println(arregloTrabajador[0].nombre)
	fmt.Print(plantillas)

	plantillas.ExecuteTemplate(w, "inicio", arregloTrabajador)
}
func crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}
func insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		apellido := r.FormValue("apellido")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO trabajadores (nombre,apellido,correo) VALUES (?,?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, apellido, correo)
		http.Redirect(w, r, "/", 301)

	}
}
