help:
	@echo "Доступные команды:"
	@echo "  make swagger     - генерация swagger документации"
	@echo "  make up          - запуск docker-контейнера"
	@echo "  make all         - генерация swagger документации и запуск docker-контейнера"
	@echo "  make down        - остановка docker-контейнера"
	@echo "  make logs        - просмотр логов контейнера"

swagger:
	@echo "Генерация swagger документации"
	@swag init -g ./cmd/main.go
	@echo "Swagger документация сгенерирована"

up:
	@docker-compose up -d --build

all: swagger up

down:
	@docker-compose down

logs: 
	@docker-compose logs -f
