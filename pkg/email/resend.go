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

type OTPEmailData struct {
	To      string
	Name    string
	OTP     string
	Expires string
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
            padding: 40px;
            margin: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #FFFFFF;
            border-radius: 16px;
            padding: 48px 40px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            border: 1px solid #E2E8F0;
        }
        .logo-container {
            text-align: center;
            margin-bottom: 4px;
        }
        .logo {
            max-width: 180px;
            height: auto;
        }
        .tagline {
            color: #64748B;
            font-size: 14px;
            font-weight: 400;
            text-align: center;
            margin-top: 4px;
            margin-bottom: 24px;
            letter-spacing: 0.3px;
        }
        .subtitle {
            color: #64748B;
            font-size: 16px;
            text-align: center;
            margin-top: 0;
            margin-bottom: 32px;
        }
        .otp-box {
            background: #F8FAFC;
            padding: 24px;
            text-align: center;
            font-size: 42px;
            font-weight: 700;
            letter-spacing: 12px;
            border-radius: 12px;
            margin: 24px 0;
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
            margin: 32px 0 24px 0;
        }
        .footer {
            color: #94A3B8;
            font-size: 14px;
            text-align: center;
        }
        .highlight {
            color: #2DD4BF;
        }
        .brand-tagline {
            color: #2DD4BF;
            font-size: 13px;
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
            <p style="font-size: 12px; color: #94A3B8; margin-top: 16px;">
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

func (s *EmailService) SendWelcome(data OTPEmailData) error {
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
            padding: 40px;
            margin: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #FFFFFF;
            border-radius: 16px;
            padding: 48px 40px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            border: 1px solid #E2E8F0;
        }
        .logo-container {
            text-align: center;
            margin-bottom: 4px;
        }
        .logo {
            max-width: 180px;
            height: auto;
        }
        .tagline {
            color: #64748B;
            font-size: 14px;
            font-weight: 400;
            text-align: center;
            margin-top: 4px;
            margin-bottom: 24px;
            letter-spacing: 0.3px;
        }
        .subtitle {
            color: #64748B;
            font-size: 16px;
            text-align: center;
            margin-top: 0;
            margin-bottom: 32px;
        }
        .card {
            background: #F8FAFC;
            border-radius: 12px;
            padding: 32px;
            margin: 24px 0;
            border: 1px solid #E2E8F0;
        }
        .card-title {
            color: #1E293B;
            font-size: 20px;
            font-weight: 600;
            margin-top: 0;
            margin-bottom: 16px;
        }
        .card-text {
            color: #1E293B;
            font-size: 16px;
            line-height: 1.8;
            margin: 0;
        }
        .feature-list {
            list-style: none;
            padding: 0;
            margin: 16px 0 0 0;
        }
        .feature-list li {
            color: #1E293B;
            font-size: 16px;
            padding: 8px 0;
            padding-left: 28px;
            background: url("data:image/svg+xml,%%3Csvg xmlns='http://www.w3.org/2000/svg' width='20' height='20' viewBox='0 0 24 24' fill='none' stroke='%%232DD4BF' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%%3E%%3Cpolyline points='20 6 9 17 4 12'%%3E%%3C/polyline%%3E%%3C/svg%%3E") left center no-repeat;
            background-size: 20px;
        }
        .divider {
            border: none;
            border-top: 1px solid #E2E8F0;
            margin: 32px 0 24px 0;
        }
        .footer {
            color: #94A3B8;
            font-size: 14px;
            text-align: center;
        }
        .highlight {
            color: #2DD4BF;
        }
        .brand-tagline {
            color: #2DD4BF;
            font-size: 13px;
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
        .mt-16 {
            margin-top: 16px;
        }
        .mb-8 {
            margin-bottom: 8px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo-container">
            <img src="https://cms.acop.co.ke/wp-content/uploads/2026/07/flownattylogo.png" alt="Flownatty" class="logo">
        </div>

        <p class="tagline">Flowing with You, Naturally</p>
        <p class="subtitle">Welcome to Flownatty!</p>

        <div class="card">
            <h2 class="card-title">Welcome aboard, %s!</h2>
            <p class="card-text">Your account has been successfully created.</p>
            <p class="card-text" style="margin-top: 16px;">You can now:</p>
            <ul class="feature-list">
                <li>Discover local businesses</li>
                <li>Shop products and services</li>
                <li>Connect with businesses</li>
                <li>Follow your favorite brands</li>
            </ul>
            <p class="card-text" style="margin-top: 16px;">
                Start exploring the app today!
            </p>
        </div>

        <hr class="divider">

        <div class="footer">
            <p class="brand-tagline">Flowing with You, Naturally</p>
            <p>Best regards,<br><strong style="color: #2DD4BF;">The Flownatty Team</strong></p>
            <p style="font-size: 12px; color: #94A3B8; margin-top: 16px;">
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

Your account has been successfully created.

You can now:
- Discover local businesses
- Shop products and services
- Connect with businesses
- Follow your favorite brands

Start exploring the app today!

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