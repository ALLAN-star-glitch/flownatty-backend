package email

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v3"
)

type EmailService struct {
	client *resend.Client
	from   string
}

func NewEmailService(apiKey, from string) *EmailService {
	return &EmailService{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

// OTPEmailData - for OTP verification emails
type OTPEmailData struct {
	To      string
	Name    string
	OTP     string
	Expires string
}

// WelcomeEmailData - for welcome emails
type WelcomeEmailData struct {
	To   string
	Name string
}

func (s *EmailService) SendOTP(data OTPEmailData) error {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OTP Verification - Flownatty</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
            background-color: #F8FAFC;
            padding: 20px;
            margin: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #FFFFFF;
            border-radius: 16px;
            padding: 32px 24px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            border: 1px solid #E2E8F0;
        }
        .logo-container {
            text-align: center;
            margin-bottom: 4px;
        }
        .logo {
            max-width: 150px;
            height: auto;
        }
        .tagline {
            color: #64748B;
            font-size: 13px;
            font-weight: 400;
            text-align: center;
            margin-top: 4px;
            margin-bottom: 20px;
            letter-spacing: 0.3px;
        }
        .subtitle {
            color: #64748B;
            font-size: 16px;
            text-align: center;
            margin-top: 0;
            margin-bottom: 24px;
        }
        .otp-box {
            background: #F8FAFC;
            padding: 16px;
            text-align: center;
            font-size: 32px;
            font-weight: 700;
            letter-spacing: 8px;
            border-radius: 12px;
            margin: 20px 0;
            color: #1E293B;
            border: 1px solid #E2E8F0;
        }
        .expiry-text {
            color: #64748B;
            font-size: 14px;
            text-align: center;
        }
        .divider {
            border: none;
            border-top: 1px solid #E2E8F0;
            margin: 24px 0 20px 0;
        }
        .footer {
            color: #94A3B8;
            font-size: 13px;
            text-align: center;
        }
        .highlight {
            color: #2DD4BF;
        }
        .brand-tagline {
            color: #2DD4BF;
            font-size: 12px;
            font-weight: 500;
            letter-spacing: 1px;
            text-transform: uppercase;
        }
        .text-primary {
            color: #1E293B;
        }
        .text-secondary {
            color: #64748B;
        }
        .text-center {
            text-align: center;
        }
        .mt-24 {
            margin-top: 24px;
        }
        @media (max-width: 480px) {
            .container {
                padding: 24px 16px;
                border-radius: 12px;
            }
            .otp-box {
                font-size: 28px;
                padding: 12px;
                letter-spacing: 6px;
            }
            .logo {
                max-width: 120px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo-container">
            <img src="https://cms.acop.co.ke/wp-content/uploads/2026/07/flownattylogo.png" alt="Flownatty" class="logo">
        </div>

        <p class="tagline">Flowing with You, Naturally</p>
        <p class="subtitle">Verify Your Account</p>

        <p class="text-primary" style="font-size: 16px;">Hi <strong>%s</strong>,</p>
        <p class="text-primary" style="font-size: 16px; line-height: 1.6;">
            Thank you for signing up for <span class="highlight">Flownatty</span>.
            Please use the OTP below to verify your account:
        </p>

        <div class="otp-box">%s</div>

        <p class="expiry-text">This OTP expires in <strong>%s</strong>.</p>

        <p class="text-secondary text-center mt-24">
            If you didn't request this, please ignore this email.
        </p>

        <hr class="divider">

        <div class="footer">
            <p class="brand-tagline">Flowing with You, Naturally</p>
            <p>Best regards,<br><strong style="color: #2DD4BF;">The Flownatty Team</strong></p>
            <p style="font-size: 11px; color: #94A3B8; margin-top: 12px;">
                Kenya's Social Commerce Super App
            </p>
        </div>
    </div>
</body>
</html>
	`, data.Name, data.OTP, data.Expires)

	text := fmt.Sprintf(`Flownatty OTP Verification
Flowing with You, Naturally

Hi %s,

Thank you for signing up for Flownatty.
Your OTP is: %s

This OTP expires in %s.

If you didn't request this, please ignore this email.

Best regards,
The Flownatty Team
--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Name, data.OTP, data.Expires)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{data.To},
		Subject: "Your OTP - Flownatty",
		Html:    html,
		Text:    text,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send OTP email to %s: %v", data.To, err)
		return err
	}

	log.Printf("OTP email sent to %s (ID: %s)", data.To, sent.Id)
	return nil
}

func (s *EmailService) SendWelcome(data WelcomeEmailData) error {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Flownatty</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
            background-color: #F8FAFC;
            padding: 20px;
            margin: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #FFFFFF;
            border-radius: 16px;
            padding: 32px 24px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            border: 1px solid #E2E8F0;
        }
        .logo-container {
            text-align: center;
            margin-bottom: 4px;
        }
        .logo {
            max-width: 150px;
            height: auto;
        }
        .tagline {
            color: #64748B;
            font-size: 13px;
            font-weight: 400;
            text-align: center;
            margin-top: 4px;
            margin-bottom: 20px;
            letter-spacing: 0.3px;
        }
        .subtitle {
            color: #64748B;
            font-size: 16px;
            text-align: center;
            margin-top: 0;
            margin-bottom: 8px;
        }
        .welcome-message {
            margin-top: 8px;
        }
        .welcome-message p {
            color: #1E293B;
            font-size: 16px;
            line-height: 1.7;
            margin: 0 0 16px 0;
        }
        .feature-grid {
            display: table;
            width: 100%%;
            margin: 16px 0 8px 0;
            border-collapse: collapse;
        }
        .feature-grid .feature-item {
            display: table-row;
        }
        .feature-grid .feature-icon {
            display: table-cell;
            vertical-align: middle;
            padding: 8px 12px 8px 0;
            width: 24px;
            color: #2DD4BF;
            font-size: 18px;
        }
        .feature-grid .feature-text {
            display: table-cell;
            vertical-align: middle;
            padding: 8px 0;
            color: #1E293B;
            font-size: 15px;
            border-bottom: 1px solid #F1F5F9;
        }
        .feature-grid .feature-item:last-child .feature-text {
            border-bottom: none;
        }
        .divider {
            border: none;
            border-top: 1px solid #E2E8F0;
            margin: 24px 0 20px 0;
        }
        .footer {
            color: #94A3B8;
            font-size: 13px;
            text-align: center;
        }
        .highlight {
            color: #2DD4BF;
        }
        .brand-tagline {
            color: #2DD4BF;
            font-size: 12px;
            font-weight: 500;
            letter-spacing: 1px;
            text-transform: uppercase;
        }
        .text-primary {
            color: #1E293B;
        }
        .text-secondary {
            color: #64748B;
        }
        .text-center {
            text-align: center;
        }
        .mt-24 {
            margin-top: 24px;
        }
        @media (max-width: 480px) {
            .container {
                padding: 24px 16px;
                border-radius: 12px;
            }
            .logo {
                max-width: 120px;
            }
            .feature-grid .feature-text {
                font-size: 14px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo-container">
            <img src="https://cms.acop.co.ke/wp-content/uploads/2026/07/flownattylogo.png" alt="Flownatty" class="logo">
        </div>

        <p class="tagline">Flowing with You, Naturally</p>
        <p class="subtitle">Welcome Aboard!</p>

        <div class="welcome-message">
            <p class="text-primary" style="font-size: 16px;">Hi <strong>%s</strong>,</p>
            <p class="text-primary" style="font-size: 16px; line-height: 1.7;">
                We're thrilled to have you join the <span class="highlight">Flownatty</span> community.
                Your account has been successfully created.
            </p>

            <p class="text-primary" style="font-size: 15px; font-weight: 600; margin-top: 20px;">
                Here's what you can do:
            </p>

            <div class="feature-grid">
                <div class="feature-item">
                    <span class="feature-icon">✦</span>
                    <span class="feature-text">Discover local businesses near you</span>
                </div>
                <div class="feature-item">
                    <span class="feature-icon">✦</span>
                    <span class="feature-text">Shop products and services effortlessly</span>
                </div>
                <div class="feature-item">
                    <span class="feature-icon">✦</span>
                    <span class="feature-text">Connect directly with businesses</span>
                </div>
                <div class="feature-item">
                    <span class="feature-icon">✦</span>
                    <span class="feature-text">Follow your favorite brands</span>
                </div>
                <div class="feature-item">
                    <span class="feature-icon">✦</span>
                    <span class="feature-text">Track orders and bookings</span>
                </div>
            </div>

            <p class="text-primary" style="font-size: 16px; line-height: 1.7; margin-top: 20px;">
                Start exploring the app today and experience commerce that flows naturally.
            </p>
        </div>

        <hr class="divider">

        <div class="footer">
            <p class="brand-tagline">Flowing with You, Naturally</p>
            <p>Best regards,<br><strong style="color: #2DD4BF;">The Flownatty Team</strong></p>
            <p style="font-size: 11px; color: #94A3B8; margin-top: 12px;">
                Kenya's Social Commerce Super App
            </p>
        </div>
    </div>
</body>
</html>
	`, data.Name)

	text := fmt.Sprintf(`Welcome to Flownatty!
Flowing with You, Naturally

Hi %s,

We're thrilled to have you join the Flownatty community.
Your account has been successfully created.

Here's what you can do:
- Discover local businesses near you
- Shop products and services effortlessly
- Connect directly with businesses
- Follow your favorite brands
- Track orders and bookings

Start exploring the app today and experience commerce that flows naturally.

Best regards,
The Flownatty Team
--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Name)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{data.To},
		Subject: "Welcome to Flownatty!",
		Html:    html,
		Text:    text,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send welcome email to %s: %v", data.To, err)
		return err
	}

	log.Printf("Welcome email sent to %s", data.To)
	return nil
}