package crud

import (
	"encoding/json"
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"modules/crud/item"
)

type Products struct {
	Items_list []item.Item `json:"products"`
	Total int `json:"total"`
	Skip int `json:"skip"`
	Limit int `json:"limit"`
}

// Open database
var db, err = sql.Open("mysql", "user:password@tcp(mariadb:3306)/poo")

func init() {
	if err != nil {
		log.Fatal(err)
	}

	// Check if the database connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	Sync()
}

func CreateTable() {
	// Execute the SQL statement to create the table
	_, err = db.Exec(
		`create table if not exists products (
			id int not null primary key,
			title varchar(255),
			description varchar(255),
			price float(17),
			discountPercentage float(17),
			rating float(17),
			stock float(17),
			brand varchar(255),
			category varchar(255),
			thumbnail varchar(255)
		);`,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteTable() {
	_, err = db.Exec(`drop table if exists products;`)
	if err != nil {
		log.Fatal(err)
	}
}

func Insert(it item.Item) (error) {
	_, err = db.Exec(fmt.Sprintf(
		`insert into products values (%d, "%s", "%s", %f, %f, %f, %f, "%s", "%s", "%s");`,
		it.Id,
		it.Title,
		it.Description,
		it.Price,
		it.DiscountPercentage,
		it.Rating,
		it.Stock,
		it.Brand,
		it.Category,
		it.Thumbnail,
	))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func Update(it item.Item) (error) {
	_, err = db.Exec(fmt.Sprintf(
		`update products set
		title = "%s",
		description = "%s",
		price = %f,
		discountPercentage = %f,
		rating = %f,
		stock = %f,
		brand = "%s",
		category = "%s",
		thumbnail = "%s"
		where id = %d;`,
		it.Title,
		it.Description,
		it.Price,
		it.DiscountPercentage,
		it.Rating,
		it.Stock,
		it.Brand,
		it.Category,
		it.Thumbnail,
		it.Id,
	))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func Delete(id int) (error) {
	_, err = db.Exec(fmt.Sprintf(`delete from products where id = %d`, id))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func Query(id int) (item.Item, error) {
	var it item.Item
	
	rows, err := db.Query(fmt.Sprintf(`select * from products where id = %d`, id))
	if err != nil {
		log.Fatal(err)
		return it, err
	}
	
	rows.Next()
	err = rows.Scan(
		&it.Id,
		&it.Title,
		&it.Description,
		&it.Price,
		&it.DiscountPercentage,
		&it.Rating,
		&it.Stock,
		&it.Brand,
		&it.Category,
		&it.Thumbnail,
	)
	if err != nil {
		log.Fatal(err)
		return it, err
	}

	return it, nil
}

func Sync() {
	DeleteTable()
	CreateTable()

	// Get the JSON file
	response, err := http.Get("https://dummyjson.com/products")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a variable to hold the decoded JSON data
	var products Products

	// Decode JSON from the response body
	err = json.NewDecoder(response.Body).Decode(&products)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(products.Items_list); i++ {
		Insert(products.Items_list[i])
	}
}