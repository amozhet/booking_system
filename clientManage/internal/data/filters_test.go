package data

import (
	"clientManage/internal/validator"
	"testing"
)

func TestValidateFilters(t *testing.T) {
	v := validator.New()
	f := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "id",
		SortSafelist: []string{"id", "-created_at"},
	}

	ValidateFilters(v, f)
	assertValid(t, v)

	f.Page = 0
	ValidateFilters(v, f)
	assertInvalid(t, v, "page", "must be greater than zero")

	f.Page = 1
	f.PageSize = 0
	ValidateFilters(v, f)
	assertInvalid(t, v, "pagesize", "must be greater than zero")

	f.PageSize = 20
	f.Sort = "invalid"
	ValidateFilters(v, f)
	assertInvalid(t, v, "sort", "invalid sort value")
}

func assertValid(t *testing.T, v *validator.Validator) {
	if !v.Valid() {
		t.Errorf("Validation failed: %v", v.Errors)
	}
}

func assertInvalid(t *testing.T, v *validator.Validator, key string, message string) {
	if v.Valid() {
		t.Error("Validation should have failed")
	}

	if v.Errors[key] != message {
		t.Errorf("Expected error message '%s' for key '%s', got '%s'", message, key, v.Errors[key])
	}
}
