.PHONY: build

build:
	go build -o bin/main ./cmd

run:
	@go run ./cmd/

.PHONY: docker

docker:
	docker-compose up -d

	@echo "\n=========================\n"
	docker images

	@echo "\n=========================\n"
	@echo "docker containers"
	docker ps -a

	@echo "Running server:"
	@echo "\n***************************"
	@echo "*                         *"
	@echo "* http://localhost:27969/ *"
	@echo "*                         *"
	@echo "***************************\n"

docker-stop:
	docker-compose down

docker-catalog:
	@echo "\n=========================\n"
	@echo "browse the catalog"
	docker exec -it web ls -la
	@echo "\n========================="

docker-delete:
	docker-compose down
	docker rmi forum_app:latest
	docker rmi golang:1.17

docker-delete-volume:
	docker volume rm forum_web

.DEFAULT_GOAL := build