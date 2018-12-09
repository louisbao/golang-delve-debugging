package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Product represents a single apple product
type Product struct {
	Name     string `json:"name"`
	Year     string `json:"year"`
	Obsolete bool   `json:"obsolete"`
}

func (p *Product) isObsolete() {
	p.Obsolete = true
}

func main() {
	products := []Product{
		{
			Name: "iPhoneXS",
			Year: "2018",
		},
		{
			Name: "Macbook",
			Year: "2018",
		},
		{
			Name: "iBook",
			Year: "2006",
		},
	}

	products[2].isObsolete()

	jsonBytes, err := json.Marshal(products)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(jsonBytes))
}
