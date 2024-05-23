package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	//"time"
	//"github.com/dgrijalva/jwt-go"

	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
)


type Tarea struct {
	Id int `json:"id"`
	Titulo string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Estado int `json:"estado"`
}


type Parametro struct {
	Id int `json:"id"`
	Key string `json:"key"`
	Value interface{} `json:"value"`
}


func main(){
	router := gin.Default()
	router.GET("/", obtenerTareas)
	router.GET("/:id", obtenerTarea)
	router.POST("/", crearTarea)
	router.PATCH("/", actualizarTarea)
	router.DELETE("/:id", eliminarTarea)

	router.Run()
}


// Funciones aplicadas a la base de datos.
func openDB() *sql.DB {
	db, error := sql.Open("sqlite3", "objetivos.db")
	if error != nil {
		fmt.Println("Error al abrir la base de datos:", error)
		panic(error)
	}
	return db // OJO: defer db.Close() donde se implemente esta funcion.
}


func getQuery(db *sql.DB, query string) []Tarea {
	defer db.Close()
	rows, error := db.Query(query)
	if error != nil {
		fmt.Println("Error al ejecutar la query:", error)
		panic(error)
	}
	defer rows.Close()

	var tareas []Tarea
	for rows.Next() {
		var tarea Tarea
		error := rows.Scan(&tarea.Id, &tarea.Titulo, &tarea.Descripcion, &tarea.Estado)
		if error != nil {
			fmt.Println("Error al escanear las rows obtenidas:", error)
			panic(error)
		}
		tareas = append(tareas, tarea)
	}
	return tareas
}


//var query = "INSERT INTO tareas (titulo, descripcion, estado) VALUES (?,?,?)"
func execQuery(db *sql.DB, query string, tarea Tarea){
	defer db.Close()
	_, error := db.Exec(query, tarea.Titulo, tarea.Descripcion, tarea.Estado)
	if error != nil {
		fmt.Println("Error al ejecutar query:", error)
		panic(error)
	}
}


func execQueryDelete(db *sql.DB, id string){
	defer db.Close()
	query := fmt.Sprintf("DELETE FROM tareas WHERE id=%s", id)
	_, error := db.Exec(query)
	if error != nil {
		fmt.Println("Error al ejecutar query:", error)
		panic(error)
	}
}


func execQueryUpdate(db *sql.DB, parametro Parametro){
	defer db.Close()
	query := fmt.Sprintf("UPDATE tareas SET %s=\"%v\" WHERE id=%v", parametro.Key, parametro.Value, parametro.Id)
	fmt.Println(query)
	_, error := db.Exec(query)
	if error != nil {
		fmt.Println("Error en la query al intentar actualizar parametro de tarea:", error)
		panic(error)
	}
}


// Controllers

func obtenerTareas(c *gin.Context){
	db := openDB()
	tareas := getQuery(db, "SELECT * FROM tareas")
	c.JSON(http.StatusOK, gin.H{"mensaje": "Bienvenido a nuestra API.", "tareas": tareas})
}


func obtenerTarea(c *gin.Context){
	db := openDB()
	id := c.Param("id")
	query := fmt.Sprintf("SELECT * FROM tareas WHERE id=%s", id)
	tareas := getQuery(db, query)
	c.JSON(http.StatusOK, gin.H{"mensaje": "Bienvenido a nuestra API.", "tareas": tareas})
}


func crearTarea(c *gin.Context){
	var tarea Tarea
	c.BindJSON(&tarea)
	db := openDB()
	execQuery(db, "INSERT INTO tareas (titulo, descripcion, estado) VALUES (?,?,?)", tarea)
	c.JSON(http.StatusOK, gin.H{"mensaje": "Tarea creada satisfactoriamente!"})
}


func actualizarTarea(c *gin.Context){
	db := openDB()
	var parametro Parametro
	c.BindJSON(&parametro)
	execQueryUpdate(db, parametro)
	c.JSON(http.StatusOK, gin.H{"mensaje": "Tarea actualizada satisfactoriamente!"})
}

func eliminarTarea(c *gin.Context){
	db := openDB()
	id := c.Param("id")
	execQueryDelete(db, id)
	c.JSON(http.StatusOK, gin.H{"mensaje": "Tarea eliminada satisfactoriamente!"})
}