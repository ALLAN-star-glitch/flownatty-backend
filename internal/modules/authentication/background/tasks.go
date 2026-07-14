package background

import (
    "encoding/json"
)

const (
    TypeEmailOTP                   = "email:otp"
    TypeEmailWelcome               = "email:welcome"
    TypeEmailPasswordResetOTP      = "email:password_reset_otp"
    TypeEmailLoginNotification     = "email:login_notification"
    TypeEmailPasswordResetConfirm  = "email:password_reset_confirm"
    TypeEmailTwoFactorOTP          = "email:two_factor_otp"
    
    // Business email tasks (NEW)
    TypeEmailBusinessOTP           = "email:business_otp"
    TypeEmailBusinessWelcome       = "email:business_welcome"
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
// BUSINESS OTP EMAIL TASK (NEW)
// ================================================

// BusinessOTPEmailTask - for business email verification
type BusinessOTPEmailTask struct {
    To           string `json:"to"`
    BusinessName string `json:"business_name"`
    OTP          string `json:"otp"`
    Expires      string `json:"expires"`
}

func (t BusinessOTPEmailTask) Payload() ([]byte, error) {
    return json.Marshal(t)
}

// ================================================
// BUSINESS WELCOME EMAIL TASK (NEW)
// ================================================

// BusinessWelcomeEmailTask - for business welcome email
type BusinessWelcomeEmailTask struct {
    To           string `json:"to"`
    BusinessName string `json:"business_name"`
    OwnerName    string `json:"owner_name"`
}

func (t BusinessWelcomeEmailTask) Payload() ([]byte, error) {
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

// ================================================
// TWO-FACTOR OTP EMAIL TASK
// ================================================

// TwoFactorOTPTask - matches email.OTPEmailData
type TwoFactorOTPTask struct {
    To      string `json:"to"`
    Name    string `json:"name"`
    OTP     string `json:"otp"`
    Expires string `json:"expires"`
}

func (t TwoFactorOTPTask) Payload() ([]byte, error) {
    return json.Marshal(t)
}