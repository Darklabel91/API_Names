package models

//InputBody struct for validation middleware
type InputBody struct {
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}
