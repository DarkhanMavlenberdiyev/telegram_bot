build:
	go build main.go
run:
	./main start
buildAPI:
	go build curl.go
runAPI:
	./curl start

depends:
	go get gopkg.in/tucnak/telebot.v2
	go get github.com/gorilla/mux
	go get github.com/urfave/cli
	go get github.com/go-pg/pg
	go get github.com/streadway/amqp
