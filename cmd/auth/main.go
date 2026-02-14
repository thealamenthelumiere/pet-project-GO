package main

import (
    "fmt"       
    "log"       
    "net/http" 

    "github.com/thealamenthelumiere/pet-project-GO/configs"
    "github.com/thealamenthelumiere/pet-project-GO/internal/handler"
    "github.com/thealamenthelumiere/pet-project-GO/internal/service"
    "github.com/thealamenthelumiere/pet-project-GO/internal/store"
)

func main() {
    // Загружаем конфигурацию и обрабатываем возможную ошибку
    cfg, err := configs.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Инициализируем зависимости
    userStore := store.NewInMemoryStore()
    userService := service.NewUserService(userStore, cfg.Secret)
    loginHandler := handler.NewLoginHandler(userService)
    verifyHandler := handler.NewVerifyHandler(userService)

    // Настраиваем маршруты
    http.HandleFunc("/login", loginHandler.Handle)
    http.HandleFunc("/verify", verifyHandler.Handle)

    // Запускаем сервер
    address := fmt.Sprintf(":%d", cfg.Port)
    log.Printf("Starting server on %s", address)
    if err := http.ListenAndServe(address, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

