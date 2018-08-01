# webshop

small generic webshop

## running project

create database `webshopGo`

default database login uses default postgres superuser with password postgres
it can be changed through input parameters

building with go 1.9.7

	# fetch deps
	dep ensure

	# make front & back
	make

	# run webshop
	./webshop (serves on default port 9000 run with --help for more options)

## deps that should be installed

go

dep (dependency management for go)

postgres (up & running)

node

npm
