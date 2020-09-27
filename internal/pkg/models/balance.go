package models

type Balance struct {
	UserID uint `json:"user_id"`
	Total float64 `json:"total"`
	Currency string `json:"currency"`
}

type ContextKey struct{}