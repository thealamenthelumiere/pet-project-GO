.PHONY: build test docker-build docker-run docker-down clean

# Сборка проекта
build:
	go build -o bin/auth-service ./cmd/auth-service

# Запуск тестов
test:
	go test -v ./...

# Сборка Docker образа
docker-build:
	docker build -t auth-service:latest .

# Запуск в Docker
docker-run:
	docker-compose up -d

# Остановка Docker
docker-down:
	docker-compose down

# Полная очистка
clean:
	rm -rf bin/
	docker-compose down -v
	docker rmi auth-service:latest || true

# Запуск с горячей перезагрузкой (для разработки)
dev:
	air

# Генерация моков
generate-mocks:
	mockgen -source=internal/service/user.go -destination=internal/service/mocks/user_service_mock.go -package=mocks