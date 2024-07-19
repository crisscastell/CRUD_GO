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
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
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
	http.HandleFunc("/borrar", borrar)
	http.HandleFunc("/editar", editar)
	http.HandleFunc("/actualizar", actualizar)
	fmt.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}

func borrar(w http.ResponseWriter, r *http.Request) {
	idTrabajador := r.URL.Query().Get("id")
	fmt.Println(idTrabajador)

	conexionEstablecida := conexionBD()
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM trabajadores WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idTrabajador)
	http.Redirect(w, r, "/", 301)

}

type Trabajador struct {
	Id       int
	Nombre   string
	Apellido string
	Correo   string
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

		trabajador.Id = id
		trabajador.Nombre = nombre
		trabajador.Apellido = apellido
		trabajador.Correo = correo

		arregloTrabajador = append(arregloTrabajador, trabajador)
	}
	//fmt.Println(arregloTrabajador)

	plantillas.ExecuteTemplate(w, "inicio", arregloTrabajador)

}

func editar(w http.ResponseWriter, r *http.Request) {
	idTrabajador := r.URL.Query().Get("id")
	fmt.Println(idTrabajador)

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM trabajadores WHERE id=?", idTrabajador)

	trabajador := Trabajador{}

	if err != nil {
		panic(err.Error())
	}

	for registro.Next() {
		var id int
		var nombre, apellido, correo string
		err = registro.Scan(&id, &nombre, &apellido, &correo)

		if err != nil {
			panic(err.Error())
		}

		trabajador.Id = id
		trabajador.Nombre = nombre
		trabajador.Apellido = apellido
		trabajador.Correo = correo
	}

	fmt.Println(trabajador)
	plantillas.ExecuteTemplate(w, "editar", trabajador)
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

func actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		apellido := r.FormValue("apellido")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()

		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE trabajadores SET nombre=?, apellido=?, correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(nombre, apellido, correo, id)
		http.Redirect(w, r, "/", 301)

	}
}
