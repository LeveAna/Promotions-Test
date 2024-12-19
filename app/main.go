package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	generator "promotions-test/dataset_generation"
	"promotions-test/services"
	"syscall"
	"time"
)

func main() {
	dsn := "root:password@tcp(db:3306)/promotions_db"
	application, err := services.SetUpApp(dsn)
	if err != nil {
		log.Fatal("Failed to set up application:", err)
	}
	defer application.DB.Close()

	// Wait for the database to be ready
	time.Sleep(5 * time.Second)

	// Generate 30,000 products with prices between 300,00 and 1500,00
	generator.GenerateRows(application.DB, 30000, 30000, 150000)

	bye := make(chan os.Signal, 1)
	signal.Notify(bye, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server running on port 8080...")
		err := application.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %q", err.Error())
		}
	}()

	// wait for the SIGINT
	sig := <-bye
	log.Printf("Detected os signal %s.", sig)

	ctx1, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	err = application.Server.Shutdown(ctx1)
	cancel()
	if err != nil {
		log.Printf("Error %s.", err.Error())
	}
}
