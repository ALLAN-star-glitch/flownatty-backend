package background

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
// HANDLE OTP EMAIL - Registration
// ================================================

// HandleOTPEmail handles registration OTP email
func (w *EmailWorker) HandleOTPEmail(ctx context.Context, task *asynq.Task) error {
    var data OTPEmailTask
    if err := json.Unmarshal(task.Payload(), &data); err != nil {
        log.Printf("Failed to parse OTP email task: %v", err)
        return err
    }

    log.Printf("Processing OTP email for %s", data.To)

    err := w.emailService.SendSignupOTP(data.To, data.Name, data.OTP, data.Expires)

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
// HANDLE BUSINESS OTP EMAIL (NEW)
// ================================================

// HandleBusinessOTPEmail handles business OTP email
func (w *EmailWorker) HandleBusinessOTPEmail(ctx context.Context, task *asynq.Task) error {
    var data BusinessOTPEmailTask
    if err := json.Unmarshal(task.Payload(), &data); err != nil {
        log.Printf("Failed to parse business OTP email task: %v", err)
        return err
    }

    log.Printf("Processing business OTP email for %s", data.To)

    err := w.emailService.SendBusinessOTP(data.To, data.BusinessName, data.OTP, data.Expires)

    if err != nil {
        log.Printf("Failed to send business OTP email to %s: %v", data.To, err)
        return err
    }

    log.Printf("Business OTP email sent to %s", data.To)
    return nil
}

// ================================================
// HANDLE BUSINESS WELCOME EMAIL (NEW)
// ================================================

// HandleBusinessWelcomeEmail handles business welcome email
func (w *EmailWorker) HandleBusinessWelcomeEmail(ctx context.Context, task *asynq.Task) error {
    var data BusinessWelcomeEmailTask
    if err := json.Unmarshal(task.Payload(), &data); err != nil {
        log.Printf("Failed to parse business welcome email task: %v", err)
        return err
    }

    log.Printf("Processing business welcome email for %s", data.To)

    err := w.emailService.SendBusinessWelcome(email.BusinessWelcomeData{
        To:           data.To,
        BusinessName: data.BusinessName,
        OwnerName:    data.OwnerName,
    })

    if err != nil {
        log.Printf("Failed to send business welcome email to %s: %v", data.To, err)
        return err
    }

    log.Printf("Business welcome email sent to %s", data.To)
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

    err := w.emailService.SendPasswordResetOTP(data.To, data.Name, data.OTP, data.Expires)

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

// ================================================
// HANDLE TWO-FACTOR OTP EMAIL
// ================================================

// HandleTwoFactorOTP handles 2FA OTP email
func (w *EmailWorker) HandleTwoFactorOTP(ctx context.Context, task *asynq.Task) error {
    var data TwoFactorOTPTask
    if err := json.Unmarshal(task.Payload(), &data); err != nil {
        log.Printf("Failed to parse 2FA OTP task: %v", err)
        return err
    }

    log.Printf("Processing 2FA OTP email for %s", data.To)

    err := w.emailService.SendTwoFactorOTP(data.To, data.Name, data.OTP, data.Expires)

    if err != nil {
        log.Printf("Failed to send 2FA OTP to %s: %v", data.To, err)
        return err
    }

    log.Printf("2FA OTP email sent to %s", data.To)
    return nil
}