# webshop

small webshop

## running project

make database `webshopGo`

default database login uses default postgres superuser (potgres with empty password)

building with go 1.9.3

	# fetch deps
	govendor sync or make deps

	# make front & back
	make

	# run webshop
	./webshop (serves on default port 9000 run --help for more options)

## deps that should be installed

go
govendor (optional)
postgres (up & running)
node
npm