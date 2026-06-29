package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	httpadapter "prepareGo/internal/adapter/in/http"
	"prepareGo/internal/adapter/out/memory"
	"prepareGo/internal/adapter/out/provider"
	"prepareGo/internal/application/service"

)

func main() {
	smsProvider := provider.NewSMSProvider()
	otpRepository := memory.NewOTPRepository()
	otpRateLimiter := memory.NewRateLimiter(3, 10*time.Minute)
	idempotencyStore := memory.NewIdempotencyStore()

	otpService := service.NewOTPService(
		smsProvider,
		otpRepository,
		otpRateLimiter,
		idempotencyStore,
		service.GenerateCode,
	)


	// mon http in handler a besoin d'un useCase interface --> service OTPService l'implement 
	otpHandler := httpadapter.NewOTPHandler(otpService)

	// Création du routeur HTTP
	mux := http.NewServeMux()
	// demandes au handler OTP d’ajouter ses routes dans le routeur.
	otpHandler.RegisterRoutes(mux)

	// Création du serveur HTTP
	// http://localhost:8080
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		// Temps maximum autorisé pour lire la requête du client.
		ReadTimeout:  5 * time.Second,
		// Temps maximum autorisé pour écrire la réponse au client.
		WriteTimeout: 5 * time.Second,
		// Temps maximum pendant lequel une connexion peut rester ouverte sans rien faire.
		IdleTimeout:  30 * time.Second,
	}

	fmt.Println("server started on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}