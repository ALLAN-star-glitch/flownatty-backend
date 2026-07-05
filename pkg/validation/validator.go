package validation

// Validator provides all validation methods
type Validator struct {
	Phone    Phone
	Password Password
}

// New creates a new validator instance
func New() *Validator {
	return &Validator{
		Phone:    Phone{},
		Password: Password{},
	}
}