package auth

import (
    "encoding/json"
)

const (
    TypeEmailOTP     = "email:otp"
    TypeEmailWelcome = "email:welcome"
)

// OTP Email Task - matches email.OTPEmailData
type OTPEmailTask struct {
    To      string `json:"to"`
    Name    string `json:"name"`
    OTP     string `json:"otp"`
    Expires string `json:"expires"`
}

func (t OTPEmailTask) Payload() ([]byte, error) {
    return json.Marshal(t)
}

// Welcome Email Task - matches email.WelcomeEmailData
type WelcomeEmailTask struct {
    To   string `json:"to"`
    Name string `json:"name"`
}

func (t WelcomeEmailTask) Payload() ([]byte, error) {
    return json.Marshal(t)
}