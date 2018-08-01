DEPS = github.com/gorilla/handlers github.com/gorilla/mux \
	github.com/lib/pq github.com/lib/pq/oid github.com/satori/go.uuid \
	github.com/senko/clog golang.org/x/crypto/bcrypt golang.org/x/crypto/blowfish

DIRECTORY = $(shell if [ ! -d "static" ]; then \
		echo "folder does not exists, creating static"; \
		mkdir "static"; \
		fi)

all: backend front

backend:
	@echo "checking for static folder"
	echo $(DIRECTORY)

	@echo "building backend"
	go build

frontend:
	@echo "getting dependencies"
	npm install

	@echo "building frontend"
	cd front; npm run build; cd ..;

	@echo "removing old frontend"
	rm -rf static/*.json static/*.ico static/*.html static/*.js

	@echo "copying front static"
	cp -r front/build/* static/

	@echo "fixing getting of js through go"
	mv static/static/js static/
	rm -rf static/static

deps:
	go get $(DEPS)

clean:
	rm -rf static/*.json static/*.ico static/*.html static/*.js
	rm -rf static/js/
	rm webshop
