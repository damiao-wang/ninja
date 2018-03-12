package main

import (
	"fmt"

	"ninja/conf"
)

func main() {
	cf := conf.Get()
	fmt.Printf(`title: %v,
		database: %v,
		host: %v,
		books: %v,
		seervers: %v
		owner: %v
		clients: %v
		products: %v`,
		cf.Title, cf.Database, cf.Hosts, cf.Books, cf.Servers, cf.Owner, cf.Clients, cf.Products)
}
