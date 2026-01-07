package main

import (
    "fmt"       // ← ДОБАВЬТЕ этот импорт
    "log"       // ← Этот уже есть
    "net/http"  // ← Этот уже есть

    "github.com/thealamenthelumiere/pet-project-GO/configs"
    "github.com/thealamenthelumiere/pet-project-GO/internal/handler"
    "github.com/thealamenthelumiere/pet-project-GO/internal/service"
    "github.com/thealamenthelumiere/pet-project-GO/internal/store"
)

func main() {
    // Загружаем конфигурацию
    cfg := configs.LoadConfig()  // ← L (заглавная), а не l (строчная)
    
    // Инициализируем зависимости
    userStore := store.NewInMemoryStore()  // ← InMemory, а не IntMemory
    userService := service.NewUserService(userStore, cfg.Secret)  // ← UserService, а не UsersService
    loginHandler := handler.NewLoginHandler(userService)
    verifyHandler := handler.NewVerifyHandler(userService)

    // Настраиваем маршруты
    http.HandleFunc("/login", loginHandler.Handle)
    http.HandleFunc("/verify", verifyHandler.Handle)

    // Запускаем сервер
    address := fmt.Sprintf(":%d", cfg.Port)  // ← правильный формат
    
    log.Printf("Starting server on %s", address)  // ← правильный синтаксис
    if err := http.ListenAndServe(address, nil); err != nil {  // ← правильный синтаксис
        log.Fatalf("Failed to start server: %v", err)
    }
}