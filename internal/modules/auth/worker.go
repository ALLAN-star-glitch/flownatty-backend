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

// ================================================
// HANDLE OTP EMAIL
// ================================================

// HandleOTPEmail handles registration OTP email
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

// ================================================
// HANDLE WELCOME EMAIL
// ================================================

// HandleWelcomeEmail handles welcome email
func (w *EmailWorker) HandleWelcomeEmail(ctx context.Context, task *asynq.Task) error {
    var data WelcomeEmailTask
    if err := json.Unmarshal(task.Payload(), &data); err != nil {
        log.Printf("Failed to parse welcome email task: %v", err)
        return err
    }

    log.Printf("Processing welcome email for %s", data.To)

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

// ================================================
// HANDLE PASSWORD RESET OTP EMAIL
// ================================================

// HandlePasswordResetOTP handles password reset OTP email
func (w *EmailWorker) HandlePasswordResetOTP(ctx context.Context, task *asynq.Task) error {
    var data PasswordResetOTPTask
    if err := json.Unmarshal(task.Payload(), &data); err != nil {
        log.Printf("Failed to parse password reset OTP task: %v", err)
        return err
    }

    log.Printf("Processing password reset OTP for %s", data.To)

    err := w.emailService.SendPasswordResetOTP(email.PasswordResetOTPData{
        To:      data.To,
        Name:    data.Name,
        OTP:     data.OTP,
        Expires: data.Expires,
    })

    if err != nil {
        log.Printf("Failed to send password reset OTP to %s: %v", data.To, err)
        return err
    }

    log.Printf("Password reset OTP sent to %s", data.To)
    return nil
}

// ================================================
// HANDLE LOGIN NOTIFICATION EMAIL
// ================================================

// HandleLoginNotification handles login notification email
func (w *EmailWorker) HandleLoginNotification(ctx context.Context, task *asynq.Task) error {
    var data LoginNotificationTask
    if err := json.Unmarshal(task.Payload(), &data); err != nil {
        log.Printf("Failed to parse login notification task: %v", err)
        return err
    }

    log.Printf("Processing login notification for %s", data.To)

    err := w.emailService.SendLoginNotification(email.LoginNotificationData{
        To:        data.To,
        Name:      data.Name,
        Time:      data.Time,
        IPAddress: data.IPAddress,
        UserAgent: data.UserAgent,
    })

    if err != nil {
        log.Printf("Failed to send login notification to %s: %v", data.To, err)
        return err
    }

    log.Printf("Login notification sent to %s", data.To)
    return nil
}

// ================================================
// HANDLE PASSWORD RESET CONFIRMATION EMAIL
// ================================================

// HandlePasswordResetConfirm handles password reset confirmation email
func (w *EmailWorker) HandlePasswordResetConfirm(ctx context.Context, task *asynq.Task) error {
    var data PasswordResetConfirmTask
    if err := json.Unmarshal(task.Payload(), &data); err != nil {
        log.Printf("Failed to parse password reset confirm task: %v", err)
        return err
    }

    log.Printf("Processing password reset confirmation for %s", data.To)

    err := w.emailService.SendPasswordResetConfirm(email.PasswordResetConfirmData{
        To:   data.To,
        Name: data.Name,
    })

    if err != nil {
        log.Printf("Failed to send password reset confirmation to %s: %v", data.To, err)
        return err
    }

    log.Printf("Password reset confirmation sent to %s", data.To)
    return nil
}