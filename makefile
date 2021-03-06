c9:
	reflex -e -c reflex.conf

mac-build:
	go build -ldflags -s

dev:
	reflex -s go run main.go

dbUpdate:
	go run ./changescripts/script.go

generate-tls:
	sudo openssl req -x509 -nodes -days 365 -newkey rsa:4096 -keyout ./key.pem -out ./cert.pem

install:
	godep get

deploy: install dbUpdate 
	go build && ./wbu-server -p
