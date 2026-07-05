package email

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v3"
)

// ================================================
// EMAIL SERVICE
// ================================================

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

// ================================================
// EMAIL DATA STRUCTS
// ================================================

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

// PasswordResetOTPData - for password reset OTP emails
type PasswordResetOTPData struct {
	To      string
	Name    string
	OTP     string
	Expires string
}

// LoginNotificationData - for login notification emails
type LoginNotificationData struct {
	To        string
	Name      string
	Time      string
	IPAddress string
	UserAgent string
}

// PasswordResetConfirmData - for password reset confirmation emails
type PasswordResetConfirmData struct {
	To   string
	Name string
}

// ================================================
// BASE EMAIL TEMPLATE
// ================================================

func getBaseTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        /* ============================================
           RESET & BASE
           ============================================ */
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
            background-color: #F8FAFC;
            padding: 20px;
            margin: 0;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #FFFFFF;
            border-radius: 16px;
            padding: 40px 32px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
            border: 1px solid #E2E8F0;
        }

        /* ============================================
           HEADER
           ============================================ */
        .header {
            text-align: center;
            padding-bottom: 24px;
            border-bottom: 1px solid #E2E8F0;
            margin-bottom: 24px;
        }
        .logo {
            max-width: 140px;
            height: auto;
        }
        .tagline {
            color: #64748B;
            font-size: 12px;
            font-weight: 500;
            letter-spacing: 1px;
            text-transform: uppercase;
            margin-top: 8px;
        }
        .brand-name {
            color: #2DD4BF;
            font-weight: 700;
        }

        /* ============================================
           TYPOGRAPHY
           ============================================ */
        .title {
            font-size: 24px;
            font-weight: 700;
            color: #1E293B;
            text-align: center;
            margin-bottom: 8px;
        }
        .subtitle {
            font-size: 16px;
            color: #64748B;
            text-align: center;
            margin-bottom: 24px;
            font-weight: 400;
        }
        .greeting {
            font-size: 16px;
            color: #1E293B;
            margin-bottom: 16px;
            font-weight: 600;
        }
        .greeting strong {
            color: #2DD4BF;
        }
        .body-text {
            font-size: 15px;
            line-height: 1.7;
            color: #1E293B;
            margin-bottom: 16px;
        }
        .body-text .highlight {
            color: #2DD4BF;
            font-weight: 600;
        }
        .body-text .secondary {
            color: #64748B;
        }
        .text-center {
            text-align: center;
        }
        .text-muted {
            color: #64748B;
            font-size: 14px;
        }

        /* ============================================
           OTP BOX
           ============================================ */
        .otp-box {
            background: #F8FAFC;
            padding: 20px;
            text-align: center;
            font-size: 36px;
            font-weight: 700;
            letter-spacing: 12px;
            border-radius: 12px;
            margin: 24px 0;
            color: #1E293B;
            border: 2px solid #E2E8F0;
            font-family: 'Courier New', monospace;
        }
        .otp-box.reset {
            background: #FEF3C7;
            border-color: #F59E0B;
            color: #92400E;
        }
        .otp-expiry {
            color: #64748B;
            font-size: 14px;
            text-align: center;
            margin-top: -8px;
            margin-bottom: 16px;
        }

        /* ============================================
           DETAILS BOX
           ============================================ */
        .details-box {
            background: #F8FAFC;
            padding: 16px 20px;
            border-radius: 12px;
            margin: 16px 0;
            border: 1px solid #E2E8F0;
        }
        .details-box .row {
            display: flex;
            justify-content: space-between;
            padding: 6px 0;
            border-bottom: 1px solid #F1F5F9;
        }
        .details-box .row:last-child {
            border-bottom: none;
        }
        .details-box .label {
            color: #64748B;
            font-size: 13px;
            font-weight: 500;
        }
        .details-box .value {
            color: #1E293B;
            font-size: 14px;
            font-weight: 500;
        }

        /* ============================================
           ALERT BOXES
           ============================================ */
        .alert-warning {
            background: #FEF3C7;
            padding: 14px 16px;
            border-radius: 8px;
            margin: 16px 0;
            border-left: 4px solid #F59E0B;
        }
        .alert-warning p {
            margin: 0;
            color: #92400E;
            font-size: 14px;
        }
        .alert-success {
            background: #D1FAE5;
            padding: 16px;
            border-radius: 12px;
            text-align: center;
            border: 1px solid #34D399;
            margin: 16px 0;
        }
        .alert-success .message {
            color: #065F46;
            font-size: 16px;
            font-weight: 600;
        }
        .alert-info {
            background: #F8FAFC;
            padding: 14px 16px;
            border-radius: 8px;
            margin: 16px 0;
            border-left: 4px solid #2DD4BF;
        }
        .alert-info p {
            margin: 0;
            color: #1E293B;
            font-size: 14px;
        }

        /* ============================================
           FEATURE GRID
           ============================================ */
        .feature-grid {
            display: table;
            width: 100%;
            margin: 16px 0 8px 0;
            border-collapse: collapse;
        }
        .feature-grid .item {
            display: table-row;
        }
        .feature-grid .icon {
            display: table-cell;
            vertical-align: middle;
            padding: 8px 12px 8px 0;
            width: 24px;
            color: #2DD4BF;
            font-size: 16px;
        }
        .feature-grid .text {
            display: table-cell;
            vertical-align: middle;
            padding: 8px 0;
            color: #1E293B;
            font-size: 14px;
            border-bottom: 1px solid #F1F5F9;
        }
        .feature-grid .item:last-child .text {
            border-bottom: none;
        }

        /* ============================================
           DIVIDER & FOOTER
           ============================================ */
        .divider {
            border: none;
            border-top: 1px solid #E2E8F0;
            margin: 28px 0 24px 0;
        }
        .footer {
            color: #94A3B8;
            font-size: 13px;
            text-align: center;
            line-height: 1.8;
        }
        .footer .brand-tagline {
            color: #2DD4BF;
            font-size: 12px;
            font-weight: 600;
            letter-spacing: 1.5px;
            text-transform: uppercase;
        }
        .footer .team-name {
            color: #2DD4BF;
            font-weight: 600;
        }
        .footer .copyright {
            font-size: 11px;
            color: #94A3B8;
            margin-top: 8px;
        }

        /* ============================================
           RESPONSIVE
           ============================================ */
        @media (max-width: 480px) {
            body {
                padding: 12px;
            }
            .container {
                padding: 24px 16px;
                border-radius: 12px;
            }
            .otp-box {
                font-size: 28px;
                padding: 16px;
                letter-spacing: 8px;
            }
            .logo {
                max-width: 110px;
            }
            .title {
                font-size: 20px;
            }
            .details-box .row {
                flex-direction: column;
                padding: 8px 0;
            }
            .details-box .value {
                margin-top: 2px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- HEADER -->
        <div class="header">
            <img src="https://cms.acop.co.ke/wp-content/uploads/2026/07/flownattylogo.png" alt="Flownatty" class="logo">
            <p class="tagline">Flowing with You, Naturally</p>
        </div>

        <!-- CONTENT -->
        {{.Content}}

        <!-- DIVIDER -->
        <hr class="divider">

        <!-- FOOTER -->
        <div class="footer">
            <p class="brand-tagline">Flowing with You, Naturally</p>
            <p>Best regards,<br><span class="team-name">The Flownatty Team</span></p>
            <p class="copyright">Kenya's Social Commerce Super App</p>
        </div>
    </div>
</body>
</html>
`
}

// ================================================
// BUILD HTML TEMPLATE HELPER
// ================================================

func buildHTML(title, content string) string {
	template := getBaseTemplate()
	return fmt.Sprintf(template, title, content)
}

// ================================================
// SEND OTP EMAIL
// ================================================

func (s *EmailService) SendOTP(data OTPEmailData) error {
	content := fmt.Sprintf(`
        <h1 class="title">Verify Your Account</h1>
        <p class="subtitle">Complete your Flownatty registration</p>

        <p class="greeting">Hi <strong>%s</strong>,</p>
        <p class="body-text">
            Thank you for signing up for <span class="highlight">Flownatty</span>.
            Please use the verification code below to complete your account setup.
        </p>

        <div class="otp-box">%s</div>
        <p class="otp-expiry">This code expires in <strong>%s</strong>.</p>

        <p class="body-text text-muted text-center">
            If you did not request this, please ignore this email.
        </p>
    `, data.Name, data.OTP, data.Expires)

	html := buildHTML("OTP Verification - Flownatty", content)

	text := fmt.Sprintf(`Flownatty - Verify Your Account

Hi %s,

Thank you for signing up for Flownatty.
Your verification code is: %s

This code expires in %s.

If you did not request this, please ignore this email.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Name, data.OTP, data.Expires)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{data.To},
		Subject: "Verify Your Account - Flownatty",
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

// ================================================
// SEND WELCOME EMAIL
// ================================================

func (s *EmailService) SendWelcome(data WelcomeEmailData) error {
	content := fmt.Sprintf(`
        <h1 class="title">Welcome to Flownatty</h1>
        <p class="subtitle">Your journey starts here</p>

        <p class="greeting">Hi <strong>%s</strong>,</p>
        <p class="body-text">
            We are thrilled to welcome you to the <span class="highlight">Flownatty</span> community.
            Your account has been successfully created.
        </p>

        <p class="body-text" style="font-weight: 600; margin-top: 20px;">
            Here is what you can do:
        </p>

        <div class="feature-grid">
            <div class="item">
                <span class="icon">&#8226;</span>
                <span class="text">Discover local businesses near you</span>
            </div>
            <div class="item">
                <span class="icon">&#8226;</span>
                <span class="text">Shop products and services effortlessly</span>
            </div>
            <div class="item">
                <span class="icon">&#8226;</span>
                <span class="text">Connect directly with businesses</span>
            </div>
            <div class="item">
                <span class="icon">&#8226;</span>
                <span class="text">Follow your favorite brands</span>
            </div>
            <div class="item">
                <span class="icon">&#8226;</span>
                <span class="text">Track orders and bookings</span>
            </div>
        </div>

        <p class="body-text" style="margin-top: 20px;">
            Start exploring the app today and experience commerce that flows naturally.
        </p>
    `, data.Name)

	html := buildHTML("Welcome to Flownatty", content)

	text := fmt.Sprintf(`Welcome to Flownatty!

Hi %s,

We are thrilled to welcome you to the Flownatty community.
Your account has been successfully created.

Here is what you can do:
- Discover local businesses near you
- Shop products and services effortlessly
- Connect directly with businesses
- Follow your favorite brands
- Track orders and bookings

Start exploring the app today and experience commerce that flows naturally.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Name)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{data.To},
		Subject: "Welcome to Flownatty",
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

// ================================================
// SEND PASSWORD RESET OTP EMAIL
// ================================================

func (s *EmailService) SendPasswordResetOTP(data PasswordResetOTPData) error {
	content := fmt.Sprintf(`
        <h1 class="title">Reset Your Password</h1>
        <p class="subtitle">Secure access to your account</p>

        <p class="greeting">Hi <strong>%s</strong>,</p>
        <p class="body-text">
            We received a request to reset your password for your <span class="highlight">Flownatty</span> account.
            Use the verification code below to continue.
        </p>

        <div class="otp-box reset">%s</div>
        <p class="otp-expiry">This code expires in <strong>%s</strong>.</p>

        <div class="alert-warning">
            <p>If you did not request this, please ignore this email.</p>
        </div>
    `, data.Name, data.OTP, data.Expires)

	html := buildHTML("Password Reset - Flownatty", content)

	text := fmt.Sprintf(`Password Reset Request - Flownatty

Hi %s,

We received a request to reset your password.
Your verification code is: %s

This code expires in %s.

If you did not request this, please ignore this email.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Name, data.OTP, data.Expires)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{data.To},
		Subject: "Password Reset OTP - Flownatty",
		Html:    html,
		Text:    text,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send password reset OTP to %s: %v", data.To, err)
		return err
	}

	log.Printf("Password reset OTP sent to %s", data.To)
	return nil
}

// ================================================
// SEND LOGIN NOTIFICATION EMAIL
// ================================================

func (s *EmailService) SendLoginNotification(data LoginNotificationData) error {
	content := fmt.Sprintf(`
        <h1 class="title">New Login Detected</h1>
        <p class="subtitle">Security notification</p>

        <p class="greeting">Hi <strong>%s</strong>,</p>
        <p class="body-text">
            We detected a new login to your <span class="highlight">Flownatty</span> account.
        </p>

        <div class="details-box">
            <div class="row">
                <span class="label">Time</span>
                <span class="value">%s</span>
            </div>
            <div class="row">
                <span class="label">IP Address</span>
                <span class="value">%s</span>
            </div>
            <div class="row">
                <span class="label">Device</span>
                <span class="value">%s</span>
            </div>
        </div>

        <div class="alert-warning">
            <p>If this was you, you can safely ignore this notification.</p>
            <p style="margin-top: 4px;">If you did not log in, please reset your password immediately.</p>
        </div>
    `, data.Name, data.Time, data.IPAddress, data.UserAgent)

	html := buildHTML("New Login - Flownatty", content)

	text := fmt.Sprintf(`New Login Detected - Flownatty

Hi %s,

We detected a new login to your Flownatty account.

Time: %s
IP Address: %s
Device: %s

If this was you, you can safely ignore this notification.
If you did not log in, please reset your password immediately.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Name, data.Time, data.IPAddress, data.UserAgent)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{data.To},
		Subject: "New Login Notification - Flownatty",
		Html:    html,
		Text:    text,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send login notification to %s: %v", data.To, err)
		return err
	}

	log.Printf("Login notification sent to %s", data.To)
	return nil
}

// ================================================
// SEND PASSWORD RESET CONFIRMATION EMAIL
// ================================================

func (s *EmailService) SendPasswordResetConfirm(data PasswordResetConfirmData) error {
	content := fmt.Sprintf(`
        <h1 class="title">Password Reset Confirmation</h1>
        <p class="subtitle">Your account security is up to date</p>

        <div class="alert-success">
            <div class="message">Password Reset Successful</div>
        </div>

        <p class="greeting">Hi <strong>%s</strong>,</p>
        <p class="body-text">
            Your <span class="highlight">Flownatty</span> password has been successfully changed.
        </p>

        <div class="alert-info">
            <p>If you did not perform this action, please contact our support team immediately.</p>
        </div>
    `, data.Name)

	html := buildHTML("Password Reset Confirmation - Flownatty", content)

	text := fmt.Sprintf(`Password Reset Confirmation - Flownatty

Hi %s,

Your Flownatty password has been successfully changed.

If you did not perform this action, please contact our support team immediately.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Name)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{data.To},
		Subject: "Password Reset Confirmation - Flownatty",
		Html:    html,
		Text:    text,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send password reset confirmation to %s: %v", data.To, err)
		return err
	}

	log.Printf("Password reset confirmation sent to %s", data.To)
	return nil
}