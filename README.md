# small webshop

## running project

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

* frontend

* database is currently mocked -> it will be upgraded to postgres db

* always refactor :)