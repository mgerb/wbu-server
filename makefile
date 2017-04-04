c9:
	reflex -e -c reflex.conf

mac-build:
	go build -ldflags -s

dev:
	reflex -s go run main.go

deploy:
	go build ./main.go && ./main -p

generate-tls:
	sudo openssl req -x509 -nodes -days 365 -newkey rsa:4096 -keyout ./key.pem -out ./cert.pem
