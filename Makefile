.PHONY: cert
cert:
	mkdir -p cert
	openssl genrsa -out cert/id_rsa 4096
	openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub

.PHONY: migrate
migrate:
	cd db/migration; ./sqitch deploy

.PHONY: seed
seed:
	go run db/seed/*.go

.PHONY: run
run:
	go run cmd/api/*.go
