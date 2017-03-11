c9:
	reflex -e -c reflex.conf

dev:
	reflex -s go run main.go

deploy:
	go build ./main.go && ./main -p

generate-tls:
	sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./tls.key -out ./tls.crt
