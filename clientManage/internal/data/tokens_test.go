package data

import (
	"clientManage/internal/validator"
	"testing"
)

func TestTokenModel_ValidateTokenPlaintext(t *testing.T) {
	v := validator.New()

	ValidateTokenPlaintext(v, "abcdefghijklmnopqrstuvwxyz")
	if !v.Valid() {
		t.Error("Valid token marked as invalid")
	}

	v = validator.New()
	ValidateTokenPlaintext(v, "")
	if v.Valid() {
		t.Error("Empty token should be invalid")
	} else if v.Errors["token"] != "must be provided" {
		t.Errorf("Incorrect error message for empty token: got %q, want %q", v.Errors["token"], "must be provided")
	}

	v = validator.New()
	ValidateTokenPlaintext(v, "short")
	if v.Valid() {
		t.Error("Short token should be invalid")
	} else if v.Errors["token"] != "must be 26 bytes long" {
		t.Errorf("Incorrect error message for short token: got %q, want %q", v.Errors["token"], "must be 26 bytes long")
	}
}
