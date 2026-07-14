package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/resend/resend-go/v3"
)

// ================================================
// EMAIL SERVICE
// ================================================

type EmailService struct {
	client *resend.Client
	from   string
	tmpl   *template.Template
}

func NewEmailService(apiKey, from string) *EmailService {
	// Parse all email templates
	tmpl := template.Must(template.New("").Parse(`
{{define "base"}}
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;background-color:#F8FAFC;padding:20px;margin:0}
        .container{max-width:600px;margin:0 auto;background:#FFFFFF;border-radius:16px;padding:40px 32px;box-shadow:0 4px 6px rgba(0,0,0,0.05);border:1px solid #E2E8F0}
        .header{text-align:center;padding-bottom:24px;border-bottom:2px solid #F1F5F9;margin-bottom:24px}
        .logo{max-width:160px;height:auto;display:inline-block}
        .tagline{color:#64748B;font-size:11px;letter-spacing:2px;text-transform:uppercase;margin-top:6px}
        .title{font-size:24px;font-weight:700;color:#1E293B;text-align:center;margin-bottom:4px}
        .title span{color:#2DD4BF}
        .subtitle{font-size:15px;color:#64748B;text-align:center;margin-bottom:24px}
        .greeting{font-size:16px;color:#1E293B;margin-bottom:12px;font-weight:600}
        .greeting strong{color:#2DD4BF}
        .body-text{font-size:15px;line-height:1.7;color:#334155;margin-bottom:14px}
        .body-text .highlight{color:#2DD4BF;font-weight:600}
        .text-muted{color:#64748B;font-size:13px}
        .text-center{text-align:center}
        .divider{border:none;border-top:2px solid #F1F5F9;margin:28px 0 20px 0}
        .footer{color:#94A3B8;font-size:12px;text-align:center;line-height:1.8}
        .footer .brand{color:#2DD4BF;font-weight:600}
        .footer .copyright{font-size:10px;color:#94A3B8;margin-top:4px}

        .otp-box{padding:20px;text-align:center;font-size:38px;font-weight:700;letter-spacing:12px;border-radius:12px;margin:20px 0;font-family:'Courier New',monospace}
        .otp-box.verify{background:#D1FAE5;color:#065F46;border:2px solid #34D399}
        .otp-box.reset{background:#FEF3C7;color:#92400E;border:2px solid #F59E0B}
        .otp-expiry{color:#64748B;font-size:13px;text-align:center;margin-top:-8px;margin-bottom:16px}
        .otp-expiry strong{color:#1E293B}

        .alert-success{background:#D1FAE5;padding:16px;border-radius:12px;text-align:center;border:2px solid #34D399;margin:16px 0}
        .alert-success .message{color:#065F46;font-size:17px;font-weight:700}
        .alert-success .sub{color:#065F46;font-size:13px;opacity:0.8}
        .alert-warning{background:#FEF3C7;padding:14px 18px;border-radius:10px;margin:14px 0;border-left:4px solid #F59E0B}
        .alert-warning p{margin:0;color:#92400E;font-size:13px}
        .alert-info{background:#EEF2FF;padding:14px 18px;border-radius:10px;margin:14px 0;border-left:4px solid #6366F1}
        .alert-info p{margin:0;color:#3730A3;font-size:13px}
        .alert-danger{background:#FEE2E2;padding:14px 18px;border-radius:10px;margin:14px 0;border-left:4px solid #EF4444}
        .alert-danger p{margin:0;color:#991B1B;font-size:13px}

        .feature-item{padding:8px 0;border-bottom:1px solid #F1F5F9;display:flex;align-items:center}
        .feature-item:last-child{border-bottom:none}
        .feature-icon{color:#2DD4BF;font-weight:700;margin-right:12px;width:20px;text-align:center}
        .feature-text{color:#1E293B;font-size:14px}

        .details-box{background:#F8FAFC;padding:16px 20px;border-radius:12px;margin:16px 0;border:1px solid #E2E8F0}
        .details-box .row{display:flex;justify-content:space-between;padding:6px 0;border-bottom:1px solid #F1F5F9}
        .details-box .row:last-child{border-bottom:none}
        .details-box .label{color:#64748B;font-size:12px;font-weight:500;text-transform:uppercase}
        .details-box .value{color:#1E293B;font-size:13px;font-weight:500}

        @media (max-width:480px){.container{padding:24px 16px}.otp-box{font-size:30px;letter-spacing:8px;padding:16px}.logo{max-width:120px}.title{font-size:20px}.details-box .row{flex-direction:column;padding:8px 0}}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <img src="https://cms.acop.co.ke/wp-content/uploads/2026/07/flownattylogo.png" alt="Flownatty" class="logo">
            <p class="tagline">Flowing with You, Naturally</p>
        </div>

        {{.Content}}

        <hr class="divider">
        <div class="footer">
            <p>Best regards,<br><span class="brand">The Flownatty Team</span></p>
            <p class="copyright">2026 Flownatty · Kenya's Social Commerce Super App</p>
        </div>
    </div>
</body>
</html>
{{end}}

{{define "otp"}}
<h1 class="title">{{.Title}} <span>{{.Subtitle}}</span></h1>
<p class="subtitle">{{.Description}}</p>

<p class="greeting">Hi <strong>{{.Name}}</strong>,</p>
<p class="body-text">
    {{.Message}}
</p>

<div class="otp-box {{.Type}}">{{.OTP}}</div>
<p class="otp-expiry">This code expires in <strong>{{.Expires}}</strong></p>

{{if .Warning}}
<div class="alert-warning">
    <p>{{.Warning}}</p>
</div>
{{end}}

<p class="body-text text-muted text-center">If you did not request this, please ignore this email.</p>
{{end}}

{{define "welcome"}}
<h1 class="title">Welcome to <span>Flownatty</span></h1>
<p class="subtitle">Your journey starts here</p>

<p class="greeting">Hi <strong>{{.Name}}</strong>,</p>
<p class="body-text">
    We are thrilled to welcome you to the <span class="highlight">Flownatty</span> community.
    Your account has been successfully created.
</p>

<div class="alert-success">
    <div class="message">Account Created Successfully</div>
    <div class="sub">Your Flownatty journey begins now</div>
</div>

<p class="body-text" style="font-weight:600;margin-top:16px">Here is what you can do:</p>

<div class="feature-item">
    <span class="feature-icon">◆</span>
    <span class="feature-text">Discover local businesses near you</span>
</div>
<div class="feature-item">
    <span class="feature-icon">◆</span>
    <span class="feature-text">Shop products and services effortlessly</span>
</div>
<div class="feature-item">
    <span class="feature-icon">◆</span>
    <span class="feature-text">Connect directly with businesses</span>
</div>
<div class="feature-item">
    <span class="feature-icon">◆</span>
    <span class="feature-text">Follow your favorite brands</span>
</div>
<div class="feature-item">
    <span class="feature-icon">◆</span>
    <span class="feature-text">Track orders and bookings</span>
</div>
{{end}}

{{define "business-otp"}}
<h1 class="title">Verify Your <span>Business Email</span></h1>
<p class="subtitle">Complete your business verification</p>

<p class="greeting">Hi <strong>{{.Name}}</strong>,</p>
<p class="body-text">
    Thank you for registering <span class="highlight">{{.BusinessName}}</span> on Flownatty.
    Please use the verification code below to confirm your business email address.
</p>

<div class="otp-box verify">{{.OTP}}</div>
<p class="otp-expiry">This code expires in <strong>{{.Expires}}</strong></p>

<div class="alert-info">
    <p>Once verified, you can start listing products and accepting orders.</p>
</div>

<p class="body-text text-muted text-center">If you did not register this business, please ignore this email.</p>
{{end}}

{{define "business-welcome"}}
<h1 class="title">Welcome <span>{{.BusinessName}}</span></h1>
<p class="subtitle">Your business is now on Flownatty</p>

<p class="greeting">Hi <strong>{{.OwnerName}}</strong>,</p>
<p class="body-text">
    Congratulations! Your business <span class="highlight">{{.BusinessName}}</span> has been successfully registered on Flownatty.
</p>

<div class="alert-success">
    <div class="message">Business Registration Complete</div>
    <div class="sub">Start selling and connecting with customers</div>
</div>

<p class="body-text" style="font-weight:600;margin-top:16px">Here is what you can do next:</p>

<div class="feature-item">
    <span class="feature-icon">◆</span>
    <span class="feature-text">Add your products and services</span>
</div>
<div class="feature-item">
    <span class="feature-icon">◆</span>
    <span class="feature-text">Set up your business profile</span>
</div>
<div class="feature-item">
    <span class="feature-icon">◆</span>
    <span class="feature-text">Connect with customers</span>
</div>
<div class="feature-item">
    <span class="feature-icon">◆</span>
    <span class="feature-text">Start receiving orders</span>
</div>

<div class="alert-info">
    <p>Complete your business profile to attract more customers.</p>
</div>
{{end}}

{{define "login-notification"}}
<h1 class="title">New <span>Login</span> Detected</h1>
<p class="subtitle">Security notification</p>

<p class="greeting">Hi <strong>{{.Name}}</strong>,</p>
<p class="body-text">
    We detected a new login to your <span class="highlight">Flownatty</span> account.
</p>

<div class="details-box">
    <div class="row">
        <span class="label">Time</span>
        <span class="value">{{.Time}}</span>
    </div>
    <div class="row">
        <span class="label">IP Address</span>
        <span class="value">{{.IPAddress}}</span>
    </div>
    <div class="row">
        <span class="label">Device</span>
        <span class="value">{{.UserAgent}}</span>
    </div>
</div>

<div class="alert-info">
    <p>If this was you, you can safely ignore this notification.</p>
</div>
<div class="alert-danger">
    <p>If you did not log in, please reset your password immediately.</p>
</div>
{{end}}

{{define "password-reset-confirm"}}
<h1 class="title">Password <span>Reset</span> Confirmation</h1>
<p class="subtitle">Your account security is up to date</p>

<div class="alert-success">
    <div class="message">Password Reset Successful</div>
    <div class="sub">Your account is now secure</div>
</div>

<p class="greeting">Hi <strong>{{.Name}}</strong>,</p>
<p class="body-text">
    Your <span class="highlight">Flownatty</span> password has been successfully changed.
</p>

<div class="alert-info">
    <p>If you did not perform this action, please contact our support team immediately.</p>
</div>
{{end}}
`))

	return &EmailService{
		client: resend.NewClient(apiKey),
		from:   from,
		tmpl:   tmpl,
	}
}

// ================================================
// EMAIL DATA STRUCTS
// ================================================

type OTPEmailData struct {
	To          string
	Name        string
	OTP         string
	Expires     string
	Title       string
	Subtitle    string
	Description string
	Message     string
	Type        string
	Warning     string
}

type WelcomeEmailData struct {
	To   string
	Name string
}

type BusinessOTPEmailData struct {
	To           string
	Name         string
	BusinessName string
	OTP          string
	Expires      string
}

type BusinessWelcomeData struct {
	To           string
	BusinessName string
	OwnerName    string
}

type LoginNotificationData struct {
	To        string
	Name      string
	Time      string
	IPAddress string
	UserAgent string
}

type PasswordResetConfirmData struct {
	To   string
	Name string
}

// ================================================
// HELPER FUNCTIONS
// ================================================

func (s *EmailService) renderEmail(templateName, title string, data interface{}) (string, error) {
	var htmlBuf bytes.Buffer

	var contentBuf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&contentBuf, templateName, data); err != nil {
		return "", err
	}

	baseData := struct {
		Title   string
		Content template.HTML
	}{
		Title:   title,
		Content: template.HTML(contentBuf.String()),
	}

	if err := s.tmpl.ExecuteTemplate(&htmlBuf, "base", baseData); err != nil {
		return "", err
	}

	return htmlBuf.String(), nil
}

func (s *EmailService) sendEmail(to, subject, html, text string) error {
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		Html:    html,
		Text:    text,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent to %s (ID: %s)", to, sent.Id)
	return nil
}

// ================================================
// GENERIC OTP SENDER
// ================================================

func (s *EmailService) SendOTP(data OTPEmailData) error {
	html, err := s.renderEmail("otp", data.Title, data)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`Flownatty - %s %s

Hi %s,

%s

Your verification code is: %s

This code expires in %s.

If you did not request this, please ignore this email.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Title, data.Subtitle, data.Name, data.Message, data.OTP, data.Expires)

	subject := fmt.Sprintf("%s %s - Flownatty", data.Title, data.Subtitle)
	return s.sendEmail(data.To, subject, html, text)
}

// ================================================
// BUSINESS EMAIL METHODS (NEW)
// ================================================

// SendBusinessOTP sends business email verification OTP
func (s *EmailService) SendBusinessOTP(to, businessName, otp, expires string) error {
	data := BusinessOTPEmailData{
		To:           to,
		Name:         to, // Use email as name if not provided
		BusinessName: businessName,
		OTP:          otp,
		Expires:      expires,
	}

	html, err := s.renderEmail("business-otp", "Verify Your Business Email", data)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`Verify Your Business Email - Flownatty

Hi,

Thank you for registering %s on Flownatty.
Please use the verification code below to confirm your business email address.

Your verification code is: %s

This code expires in %s.

Once verified, you can start listing products and accepting orders.

If you did not register this business, please ignore this email.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, businessName, otp, expires)

	return s.sendEmail(to, "Verify Your Business Email - Flownatty", html, text)
}

// SendBusinessWelcome sends business welcome email
func (s *EmailService) SendBusinessWelcome(data BusinessWelcomeData) error {
	html, err := s.renderEmail("business-welcome", "Welcome to Flownatty", data)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`Welcome %s to Flownatty!

Hi %s,

Congratulations! Your business %s has been successfully registered on Flownatty.

Here is what you can do next:
◆ Add your products and services
◆ Set up your business profile
◆ Connect with customers
◆ Start receiving orders

Complete your business profile to attract more customers.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.BusinessName, data.OwnerName, data.BusinessName)

	return s.sendEmail(data.To, "Welcome to Flownatty - Business Registration Complete", html, text)
}

// ================================================
// CONVENIENCE FUNCTIONS
// ================================================

func (s *EmailService) SendSignupOTP(to, name, otp, expires string) error {
	return s.SendOTP(OTPEmailData{
		To:          to,
		Name:        name,
		OTP:         otp,
		Expires:     expires,
		Title:       "Verify Your",
		Subtitle:    "Account",
		Description: "Complete your Flownatty registration",
		Message:     "Thank you for signing up for Flownatty. Please use the verification code below to complete your account setup.",
		Type:        "verify",
	})
}

func (s *EmailService) SendTwoFactorOTP(to, name, otp, expires string) error {
	return s.SendOTP(OTPEmailData{
		To:          to,
		Name:        name,
		OTP:         otp,
		Expires:     expires,
		Title:       "Two-Factor",
		Subtitle:    "Authentication",
		Description: "Secure your account access",
		Message:     "You requested a two-factor authentication code for your Flownatty account. Please use the verification code below to complete your login.",
		Type:        "verify",
	})
}

func (s *EmailService) SendPasswordResetOTP(to, name, otp, expires string) error {
	return s.SendOTP(OTPEmailData{
		To:          to,
		Name:        name,
		OTP:         otp,
		Expires:     expires,
		Title:       "Reset Your",
		Subtitle:    "Password",
		Description: "Secure access to your account",
		Message:     "We received a request to reset your password for your Flownatty account. Use the verification code below to continue.",
		Type:        "reset",
		Warning:     "If you did not request this, please ignore this email.",
	})
}

func (s *EmailService) SendLoginOTP(to, name, otp, expires string) error {
	return s.SendOTP(OTPEmailData{
		To:          to,
		Name:        name,
		OTP:         otp,
		Expires:     expires,
		Title:       "Login",
		Subtitle:    "Verification",
		Description: "Secure access to your account",
		Message:     "You requested to log in to your Flownatty account. Please use the verification code below to complete your login.",
		Type:        "verify",
	})
}

func (s *EmailService) SendEmailChangeOTP(to, name, otp, expires, newEmail string) error {
	return s.SendOTP(OTPEmailData{
		To:          to,
		Name:        name,
		OTP:         otp,
		Expires:     expires,
		Title:       "Verify Your",
		Subtitle:    "Email Change",
		Description: "Confirm your new email address",
		Message:     fmt.Sprintf("You requested to change your email address to %s. Please use the verification code below to confirm this change.", newEmail),
		Type:        "verify",
		Warning:     "If you did not request this, please contact support immediately.",
	})
}

func (s *EmailService) SendPhoneChangeOTP(to, name, otp, expires, newPhone string) error {
	return s.SendOTP(OTPEmailData{
		To:          to,
		Name:        name,
		OTP:         otp,
		Expires:     expires,
		Title:       "Verify Your",
		Subtitle:    "Phone Change",
		Description: "Confirm your new phone number",
		Message:     fmt.Sprintf("You requested to change your phone number to %s. Please use the verification code below to confirm this change.", newPhone),
		Type:        "verify",
		Warning:     "If you did not request this, please contact support immediately.",
	})
}

func (s *EmailService) SendWelcome(data WelcomeEmailData) error {
	html, err := s.renderEmail("welcome", "Welcome to Flownatty", data)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`Welcome to Flownatty!

Hi %s,

We are thrilled to welcome you to the Flownatty community.
Your account has been successfully created.

Here is what you can do:
◆ Discover local businesses near you
◆ Shop products and services effortlessly
◆ Connect directly with businesses
◆ Follow your favorite brands
◆ Track orders and bookings

Start exploring the app today and experience commerce that flows naturally.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Name)

	return s.sendEmail(data.To, "Welcome to Flownatty", html, text)
}

func (s *EmailService) SendLoginNotification(data LoginNotificationData) error {
	html, err := s.renderEmail("login-notification", "New Login Detected", data)
	if err != nil {
		return err
	}

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

	return s.sendEmail(data.To, "New Login Notification - Flownatty", html, text)
}

func (s *EmailService) SendPasswordResetConfirm(data PasswordResetConfirmData) error {
	html, err := s.renderEmail("password-reset-confirm", "Password Reset Confirmation", data)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`Password Reset Confirmation - Flownatty

Hi %s,

Your Flownatty password has been successfully changed.

If you did not perform this action, please contact our support team immediately.

--
Flowing with You, Naturally
Kenya's Social Commerce Super App
`, data.Name)

	return s.sendEmail(data.To, "Password Reset Confirmation - Flownatty", html, text)
}