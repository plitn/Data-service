package main

import (
	"fmt"
	"github.com/Data-service/internal/repository"
	"github.com/Data-service/internal/service/data_service"
	"github.com/Data-service/internal/service/handler"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"time"
)

func main() {
	//	cfg := config.LoadConfig()
	dbHost := os.Getenv("DB_HOST")
	minioHost := os.Getenv("MINIO_HOST")
	connString := fmt.Sprintf("host=%s port=5432 user=postgres password=postgres dbname=postgres sslmode=disable timezone=UTC", dbHost)
	dbInst, err := sqlx.Open("postgres", connString)
	if err != nil {
		fmt.Println("failed to connect database")
		return
	}
	dbInst.SetConnMaxLifetime(10 * time.Minute)
	err = dbInst.Ping()
	if err != nil {
		fmt.Printf("failed to ping database, err: %v\n", err)
		return
	}
	fmt.Println("connected to database")
	defer dbInst.Close()
	repo := repository.NewRepository(dbInst)
	minioConn := fmt.Sprintf("%s:9000", minioHost)
	dataService := data_service.NewDataService(repo, minioConn, "R6OupGDd8kdz8VCnBp0Z", "aMTnjuPtU99DXj7Y444tzq86pDYoDY8w8PFVddLz")
	handlerService := handler.New(dataService)
	fmt.Println("services created")

	mux := chi.NewRouter()
	// todo при каждом обращении к файлу получать новый линк и сохранять в бд его
	// todo ручка которая запускает питон скрипт с файлами, а на фронт отдает байты

	mux.Post("/save-data", handlerService.CreateFile)
	mux.Get("/get-item", handlerService.GetItem)
	mux.Put("/update-item", handlerService.UpdateItem) // тут поправить чтобы не затирались поля лишние
	mux.Delete("/delete-item", handlerService.DeleteItem)
	mux.Get("/get-user-items", handlerService.GetItemByUserId) // получение всех файлов по пользователю

	mux.Post("/cut-item", handlerService.CutItem)

	mux.Post("/create-capsule", handlerService.CreateCapsule)
	mux.Get("/get-capsule", handlerService.GetCapsule)
	mux.Get("/get-user-capsules", handlerService.GetCapsulesByUser)
	mux.Post("/update-capsule", handlerService.UpdateCapsule)
	mux.Delete("/delete-capsule", handlerService.DeleteCapsule)
	mux.Post("/add-to-capsule", handlerService.AddItemToCapsule)
	mux.Delete("/delete-from-capsule", handlerService.DeleteItemFromCapsule)
	mux.Get("/get-all-capsule", handlerService.GetItemsFromCapsule)

	mux.Get("/get-look-data", handlerService.GetLookData)
	mux.Post("/add-to-look", handlerService.AddToLook)
	mux.Delete("/delete-look", handlerService.DeleteLook)
	mux.Post("/create-look", handlerService.CreateLook)
	fmt.Println("starting server on port 8080")
	httpServer := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", 8080),
		Handler: mux,
	}

	fmt.Printf("listening to http://0.0.0.0:%d/ for debug http", 8080)
	if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("failed to listen on port 8080: %v", err)
	}
}
