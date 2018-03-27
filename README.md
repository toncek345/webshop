# webshop

small webshop

## running project

make database `webshopGo`

default database login uses default postgres superures (potgres with empty password)

building with go 1.9.3

	go get github.com/gorilla/mux
	go get github.com/senko/clog
	go get github.com/lib/pq
	go get github.com/satori/go.uuid

	go run main.go
	
	or
	
	go build
	./webshop

## todo

* each product -> multiple images
* frontend
* always refactor :)
