
MOCKERY_VERSION=2.43.0

IMAGE_NAME=user
CONTAINER_NAME=user-api

run:
	go run cmd/main.go

createuserreq:
	curl -X POST localhost:3000/users -d '{"firstName":"enes"}'

compose-up:
	docker compose up --build -d 

compose-down:
	docker compose down

mockery-install:
	go install github.com/vektra/mockery/v2/...@v$(MOCKERY_VERSION)

mocks: mockery-install
	rm -rf internal/mocks
	mockery  --dir internal --all --keeptree --output internal/mocks

.PHONY: docker-build
docker-build:
	docker build --tag $(IMAGE_NAME) .

docker-run: 
	docker run --env-file .env -p 8080:3000 --name $(CONTAINER_NAME) $(IMAGE_NAME):latest 
docker-stop:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

test:
	go test ./...