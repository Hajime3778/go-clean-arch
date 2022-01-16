package auth

import "gopkg.in/go-playground/validator.v9"

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// IsSignUpRequestValid:
func (r SignUpRequest) IsSignUpRequestValid() (bool, error) {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return false, err
	}
	return true, nil
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// IsSignUpRequestValid:
func (r SignInRequest) IsSignInRequestValid() (bool, error) {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return false, err
	}
	return true, nil
}
