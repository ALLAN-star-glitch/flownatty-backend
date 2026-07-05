package auth

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ALLAN-star-glitch/flownatty-backend/pkg/email"
	"github.com/hibiken/asynq"
)

type EmailWorker struct {
	emailService *email.EmailService
}

func NewEmailWorker(emailService *email.EmailService) *EmailWorker {
	return &EmailWorker{
		emailService: emailService,
	}
}

// Handle OTP Email
func (w *EmailWorker) HandleOTPEmail(ctx context.Context, task *asynq.Task) error {
	var data OTPEmailTask
	if err := json.Unmarshal(task.Payload(), &data); err != nil {
		log.Printf("Failed to parse OTP email task: %v", err)
		return err
	}

	log.Printf("Processing OTP email for %s", data.To)

	err := w.emailService.SendOTP(email.OTPEmailData{
		To:      data.To,
		Name:    data.Name,
		OTP:     data.OTP,
		Expires: data.Expires,
	})

	if err != nil {
		log.Printf("Failed to send OTP email to %s: %v", data.To, err)
		return err
	}

	log.Printf("OTP email sent to %s", data.To)
	return nil
}

// Handle Welcome Email
func (w *EmailWorker) HandleWelcomeEmail(ctx context.Context, task *asynq.Task) error {
	var data WelcomeEmailTask
	if err := json.Unmarshal(task.Payload(), &data); err != nil {
		log.Printf("Failed to parse welcome email task: %v", err)
		return err
	}

	log.Printf("Processing welcome email for %s", data.To)

	// ✅ Use WelcomeEmailData (not OTPEmailData)
	err := w.emailService.SendWelcome(email.WelcomeEmailData{
		To:   data.To,
		Name: data.Name,
	})

	if err != nil {
		log.Printf("Failed to send welcome email to %s: %v", data.To, err)
		return err
	}

	log.Printf("Welcome email sent to %s", data.To)
	return nil
}