package auth

import (
    "encoding/json"
)

const (
    TypeEmailOTP                   = "email:otp"
    TypeEmailWelcome               = "email:welcome"
    TypeEmailPasswordResetOTP      = "email:password_reset_otp"
    TypeEmailLoginNotification     = "email:login_notification"
    TypeEmailPasswordResetConfirm  = "email:password_reset_confirm"
)

// ================================================
// OTP EMAIL TASK
// ================================================

// OTPEmailTask - matches email.OTPEmailData
type OTPEmailTask struct {
    To      string `json:"to"`
    Name    string `json:"name"`
    OTP     string `json:"otp"`
    Expires string `json:"expires"`
}

func (t OTPEmailTask) Payload() ([]byte, error) {
    return json.Marshal(t)
}

// ================================================
// WELCOME EMAIL TASK
// ================================================

// WelcomeEmailTask - matches email.WelcomeEmailData
type WelcomeEmailTask struct {
    To   string `json:"to"`
    Name string `json:"name"`
}

func (t WelcomeEmailTask) Payload() ([]byte, error) {
    return json.Marshal(t)
}

// ================================================
// PASSWORD RESET OTP EMAIL TASK
// ================================================

// PasswordResetOTPTask - matches email.PasswordResetOTPData
type PasswordResetOTPTask struct {
    To      string `json:"to"`
    Name    string `json:"name"`
    OTP     string `json:"otp"`
    Expires string `json:"expires"`
}

func (t PasswordResetOTPTask) Payload() ([]byte, error) {
    return json.Marshal(t)
}

// ================================================
// LOGIN NOTIFICATION EMAIL TASK
// ================================================

// LoginNotificationTask - matches email.LoginNotificationData
type LoginNotificationTask struct {
    To        string `json:"to"`
    Name      string `json:"name"`
    Time      string `json:"time"`
    IPAddress string `json:"ip_address"`
    UserAgent string `json:"user_agent"`
}

func (t LoginNotificationTask) Payload() ([]byte, error) {
    return json.Marshal(t)
}

// ================================================
// PASSWORD RESET CONFIRMATION EMAIL TASK
// ================================================

// PasswordResetConfirmTask - matches email.PasswordResetConfirmData
type PasswordResetConfirmTask struct {
    To   string `json:"to"`
    Name string `json:"name"`
}

func (t PasswordResetConfirmTask) Payload() ([]byte, error) {
    return json.Marshal(t)
}