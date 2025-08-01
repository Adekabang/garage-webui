package main

import (
	"Adekabang/garage-webui/router"
	"Adekabang/garage-webui/ui"
	"Adekabang/garage-webui/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize app
	godotenv.Load()
	utils.InitCacheManager()
	sessionMgr := utils.InitSessionManager()

	if err := utils.Garage.LoadConfig(); err != nil {
		log.Println("Cannot load garage config!", err)
	}

	basePath := os.Getenv("BASE_PATH")
	mux := http.NewServeMux()

	// Serve API
	apiPrefix := basePath + "/api"
	mux.Handle(apiPrefix+"/", http.StripPrefix(apiPrefix, router.HandleApiRouter()))

	// Static files
	ui.ServeUI(mux)

	// Redirect to UI if BASE_PATH is set
	if basePath != "" {
		mux.Handle("/", http.RedirectHandler(basePath, http.StatusMovedPermanently))
	}

	host := utils.GetEnv("HOST", "0.0.0.0")
	port := utils.GetEnv("PORT", "3909")

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting server on http://%s", addr)

	if err := http.ListenAndServe(addr, sessionMgr.LoadAndSave(mux)); err != nil {
		log.Fatal(err)
	}
}
