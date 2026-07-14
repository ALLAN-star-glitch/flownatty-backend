package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
	"github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/authentication/background"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/email"
	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/redis"
	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.Load()

	// ================================================
	// Initialize Redis (Global)
	// ================================================
	if err := redis.Init(cfg.Redis.URL); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	log.Println("Redis initialized successfully")

	emailService := email.NewEmailService(cfg.Resend.ApiKey, cfg.Resend.From)
	emailWorker := background.NewEmailWorker(emailService)

	// Create Asynq server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.URL},
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	// Register handlers
	mux.HandleFunc(background.TypeEmailOTP, emailWorker.HandleOTPEmail)
	mux.HandleFunc(background.TypeEmailWelcome, emailWorker.HandleWelcomeEmail)
	mux.HandleFunc(background.TypeEmailBusinessOTP, emailWorker.HandleBusinessOTPEmail)
	mux.HandleFunc(background.TypeEmailBusinessWelcome, emailWorker.HandleBusinessWelcomeEmail)
	mux.HandleFunc(background.TypeEmailPasswordResetOTP, emailWorker.HandlePasswordResetOTP)
	mux.HandleFunc(background.TypeEmailLoginNotification, emailWorker.HandleLoginNotification)
	mux.HandleFunc(background.TypeEmailPasswordResetConfirm, emailWorker.HandlePasswordResetConfirm)
	mux.HandleFunc(background.TypeEmailTwoFactorOTP, emailWorker.HandleTwoFactorOTP)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()

	log.Println("Asynq worker started. Press Ctrl+C to stop.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	srv.Shutdown()
}