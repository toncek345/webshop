# webshop

small generic webshop that is far from ideal go and needs major refactor :)

## running project

create database `webshopGo`

default database login uses default postgres superuser with password postgres
it can be changed through input parameters

building with go 1.9.7

	# fetch deps
	dep ensure

	# backend
	go run main.go (see other args with --help)

	# frontend
	cd front
	`npm run start` or `yarn start`
	# or build frontend and serve it
	`npm run build` or `yarn build`

## deps that should be installed

go

dep (dependency management for go)

postgres (up & running)

node

npm
