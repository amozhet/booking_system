package data

import (
	"clientManage/internal/validator"
	"testing"
)

func TestUserModel_ValidateUser(t *testing.T) {
	user := &User{
		Fname:     "Primer",
		Sname:     "Primer",
		Email:     "primer@example.com",
		Password:  password{plaintext: new(string)},
		Activated: false,
	}
	err := user.Password.Set("validPassword123")
	if err != nil {
		t.Fatal(err)
	}

	v := validator.New()
	ValidateUser(v, user)
	if !v.Valid() {
		t.Errorf("Valid user marked as invalid: %v", v.Errors)
	}

	user.Fname = ""
	v = validator.New()
	ValidateUser(v, user)
	if v.Valid() {
		t.Error("User with empty fname should be invalid")
	} else if v.Errors["fname"] != "must be provided" {
		t.Errorf("Incorrect error message for empty fname: got %q, want %q", v.Errors["fname"], "must be provided")
	}

	user.Sname = " "
	v = validator.New()
	ValidateUser(v, user)
	if v.Valid() {
		t.Error("Empty sname should be invalid")
	} else if v.Errors["sname"] != "must be provided" {
		t.Errorf("Incorrect error message for empty sname: got %q, want %q", v.Errors["sname"], "must be provided")
	}

	user.Sname = "Averylongsurnamethatexceedsthe engthlimitof500bytesssssssssssssssss"
	v = validator.New()
	ValidateUser(v, user)
	if v.Valid() {
		t.Error("Too long sname should be invalid")
	} else if v.Errors["sname"] != "must not be more than 500 bytes long" {
		t.Errorf("Incorrect error message for too long sname: got %q, want %q", v.Errors["sname"], "must not be more than 500 bytes long")
	}

}

func TestUserModel_ValidateEmail(t *testing.T) {
	v := validator.New()

	ValidateEmail(v, "user@example.com")
	if !v.Valid() {
		t.Error("Valid email marked as invalid")
	}

	v = validator.New()
	ValidateEmail(v, "test---test@example.com")
	if !v.Valid() {
		t.Error("Valid email with '-' marked as invalid")
	}

	v = validator.New()
	ValidateEmail(v, "")
	if v.Valid() {
		t.Error("Empty email should be invalid")
	} else if v.Errors["email"] != "must be provided" {
		t.Errorf("Incorrect error message for empty email: got %q, want %q", v.Errors["email"], "must be provided")
	}

	v = validator.New()
	ValidateEmail(v, "user")
	if v.Valid() {
		t.Error("Invalid email format should be invalid")
	} else if v.Errors["email"] != "must be a valid email address" {
		t.Errorf("Incorrect error message for invalid email: got %q, want %q", v.Errors["email"], "must be a valid email address")
	}
}

func TestUserModel_ValidatePasswordPlainText(t *testing.T) {
	v := validator.New()

	ValidatePasswordPlainText(v, "Password123")
	if !v.Valid() {
		t.Error("Valid password marked as invalid")
	}

	v = validator.New()
	ValidatePasswordPlainText(v, "")
	if v.Valid() {
		t.Error("Empty password should be invalid")
	} else if v.Errors["password"] != "must be provided" {
		t.Errorf("Incorrect error message for empty password: got %q, want %q", v.Errors["password"], "must be provided")
	}

	v = validator.New()
	ValidatePasswordPlainText(v, "short")
	if v.Valid() {
		t.Error("Short password should be invalid")
	} else if v.Errors["password"] != "must be at least 8 bytes long" {
		t.Errorf("Incorrect error message for short password: got %q, want %q", v.Errors["password"], "must be at least 8 bytes long")
	}

	v = validator.New()
	ValidatePasswordPlainText(v, "longlonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonlonglonglonglonglonglonglonglonglonglonglonglonlonglong9")
	if v.Valid() {
		t.Error("Too long password should be invalid")
	} else if v.Errors["password"] != "must not be more than 71 bytes long" {
		t.Errorf("Incorrect error message for too long password: got %q, want %q", v.Errors["password"], "must not be more than 71 bytes long")
	}
}
