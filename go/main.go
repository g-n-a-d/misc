package main

import (
	"os"
	"encoding/json"
	"database/sql"
	"io/ioutil"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type Products struct {
	P_list []p `json:"products"`
	Total int `json:"total"`
	Skip int `json:"skip"`
	Limit int `json:"limit"`
}

type p struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price float32 `json:"price"`
	DiscountPercentage float32 `json:"discountPercentage"`
	Rating float32 `json:"rating"`
	Stock float32 `json:"stock"`
	Brand string `json:"brand"`
	Category string `json:"category"`
	Thumbnail string `json:"thumbnail"`
	Images []string `json:"images"`
}

func main() {
	// open database
	db, err := sql.Open("mysql", "user:user@tcp(localhost:3306)/poo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the database connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL statement to create the table
	_, err = db.Exec(
		`create table if not exists product (
			id int not null primary key,
			title varchar(255),
			description varchar(255),
			price float(17),
			discountPercentage float(17),
			rating float(17),
			stock float(17),
			brand varchar(255),
			category varchar(255),
			thumbnail varchar(255),
			images varchar(255)
		);`,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Open the JSON file
	file, err := os.Open("./products.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(file)

	var products Products

	json.Unmarshal(byteValue, &products)

	// fmt.Println(products)

	for i := 0; i < len(products.P_list); i++ {
		_, err = db.Exec(fmt.Sprintf(
			`insert into product values (%d, "%s", "%s", %f, %f, %f, %f, "%s", "%s", "%s", "%s");`,
			products.P_list[i].Id,
			products.P_list[i].Title,
			products.P_list[i].Description,
			products.P_list[i].Price,
			products.P_list[i].DiscountPercentage,
			products.P_list[i].Rating,
			products.P_list[i].Stock,
			products.P_list[i].Brand,
			products.P_list[i].Category,
			products.P_list[i].Thumbnail,
			products.P_list[i].Images[0],
		))
		if err != nil {
			log.Fatal(err)
		}
	}
}