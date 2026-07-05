package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/ALLAN-star-glitch/flownatty-backend/internal/config"
    "github.com/ALLAN-star-glitch/flownatty-backend/internal/modules/auth"
    "github.com/ALLAN-star-glitch/flownatty-backend/pkg/email"
    "github.com/hibiken/asynq"
)

func main() {
    cfg := config.Load()

    emailService := email.NewEmailService(cfg.Resend.ApiKey, cfg.Resend.From)
    emailWorker := auth.NewEmailWorker(emailService)


	//  Create Asynq server
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
    
    //  Register handlers
    mux.HandleFunc(auth.TypeEmailOTP, emailWorker.HandleOTPEmail)
    mux.HandleFunc(auth.TypeEmailWelcome, emailWorker.HandleWelcomeEmail)
    mux.HandleFunc(auth.TypeEmailPasswordResetOTP, emailWorker.HandlePasswordResetOTP)
    mux.HandleFunc(auth.TypeEmailLoginNotification, emailWorker.HandleLoginNotification)
    mux.HandleFunc(auth.TypeEmailPasswordResetConfirm, emailWorker.HandlePasswordResetConfirm) 

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