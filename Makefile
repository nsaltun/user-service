
MOCKERY_VERSION=2.43.0

run:
	go run cmd/main.go

createuserreq:
	curl -X POST localhost:3000/users -d '{"firstName":"enes"}'

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

mockery-install:
	go install github.com/vektra/mockery/v2/...@v$(MOCKERY_VERSION)

mocks: mockery-install
	rm -rf internal/mocks
	mockery  --dir internal --all --keeptree --output internal/mocks