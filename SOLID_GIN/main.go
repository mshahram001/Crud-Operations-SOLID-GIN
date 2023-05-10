package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Name string
	Age  int
	City string
}

type base interface {
	AddUser(u *User) error
	GetUser(u *User) error
	UpdateUser(u *User) error
	delete_user(u *User) error
}

type HTTPHandler struct {
	db base
}
type database1 struct {
	db *sql.DB
}

func NewMySQLdbase(connectionString string) (*database1, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &database1{db}, nil
}

// GET

func (d *database1) GetUser(u *User) error {
	sql_query := fmt.Sprintf(`SELECT * FROM info WHERE ID='%d'`, u.ID)
	_, err := d.db.Exec(sql_query)
	return err
}

func (h *HTTPHandler) GetUser(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.db.GetUser(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, user)
}

// ADD
func (d *database1) AddUser(u *User) error {
	query_data := fmt.Sprintf(`INSERT INTO info VALUES(%d,"%s",%d,"%s")`, u.ID, u.Name, u.Age, u.City)
	_, err := d.db.Exec(query_data)
	return err
}

func (h *HTTPHandler) AddUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := h.db.AddUser(&user); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, user)

}

// UPDATE
func (d *database1) UpdateUser(u *User) error {
	query_data := fmt.Sprintf("UPDATE info SET Name='%s', Age=%d, City='%s' WHERE ID=%d", u.Name, u.Age, u.City, u.ID)
	_, err := d.db.Exec(query_data)
	return err
}

func (h *HTTPHandler) UpdateUser(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	err = h.db.UpdateUser(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, user)
}

// DELETE
func (d *database1) delete_user(u *User) error {
	query_data := fmt.Sprintf("DELETE FROM info WHERE ID = %d", u.ID)
	_, err := d.db.Exec(query_data)
	return err
}

func (h *HTTPHandler) delete_user(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.db.delete_user(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, user)
	fmt.Println("User deleted succesfully!!!")
}

func Err(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

func db_creation() {
	//connecting to mysql
	db, err := sql.Open("mysql", "root:india@123@tcp(localhost:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// database creation
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS solid_crud")
	if err != nil {
		panic(err)
	}
}

func db_connection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:india@123@tcp(localhost:3306)/solid_crud")

	if err != nil {
		return nil, err
	}
	return db, nil
}

func user_table_creation() {
	db, err := db_connection()
	if err != nil {
		panic(err)
	}
	// SQL table creation
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS info (ID INT NOT NULL, Name VARCHAR(20), Age INT, City VARCHAR(20), PRIMARY KEY (ID));")
	if err != nil {
		panic(err)
	}
	fmt.Println("info Table Created")
}

func main() {
	db_creation()
	user_table_creation()
	// DB connectivity
	db, err := NewMySQLdbase("root:india@123@tcp(localhost:3306)/solid_crud")
	if err != nil {
		log.Fatal(err)
	}
	defer db.db.Close()

	handler := &HTTPHandler{db}

	router := gin.Default()
	router.POST("/postuser", handler.AddUser)
	router.GET("/getuser", handler.GetUser)
	router.PUT("/updateuser", handler.UpdateUser)
	router.DELETE("/deleteuser", handler.delete_user)
	router.Run("localhost:8000")
}
