package validator

import (
	"testing"
)

func TestValidate_ValidStruct(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
		Age   int    `validate:"required,gte=0,lte=150"`
	}

	valid := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	err := Validate(valid)
	if err != nil {
		t.Errorf("Validate() unexpected error for valid struct: %v", err)
	}
}

func TestValidate_MissingRequiredField(t *testing.T) {
	type TestStruct struct {
		Name string `validate:"required"`
	}

	invalid := TestStruct{
		Name: "",
	}

	err := Validate(invalid)
	if err == nil {
		t.Error("Validate() expected error for missing required field")
	}

	if !contains(err.Error(), "Name") || !contains(err.Error(), "required") {
		t.Errorf("Validate() error message should mention field and validation: %v", err)
	}
}

func TestValidate_InvalidEmail(t *testing.T) {
	type TestStruct struct {
		Email string `validate:"required,email"`
	}

	invalid := TestStruct{
		Email: "not-an-email",
	}

	err := Validate(invalid)
	if err == nil {
		t.Error("Validate() expected error for invalid email")
	}
}

func TestValidate_NumberRange(t *testing.T) {
	type TestStruct struct {
		Age int `validate:"gte=0,lte=150"`
	}

	tests := []struct {
		name    string
		age     int
		wantErr bool
	}{
		{"valid age", 30, false},
		{"minimum age", 0, false},
		{"maximum age", 150, false},
		{"negative age", -1, true},
		{"too old", 151, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TestStruct{Age: tt.age}
			err := Validate(s)

			if tt.wantErr && err == nil {
				t.Error("Validate() expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestValidate_MinMaxLength(t *testing.T) {
	type TestStruct struct {
		Items []string `validate:"min=1,max=5"`
	}

	tests := []struct {
		name    string
		items   []string
		wantErr bool
	}{
		{"valid length", []string{"a", "b"}, false},
		{"minimum length", []string{"a"}, false},
		{"maximum length", []string{"a", "b", "c", "d", "e"}, false},
		{"empty array", []string{}, true},
		{"too many items", []string{"a", "b", "c", "d", "e", "f"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TestStruct{Items: tt.items}
			err := Validate(s)

			if tt.wantErr && err == nil {
				t.Error("Validate() expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestValidate_OneOf(t *testing.T) {
	type TestStruct struct {
		Role string `validate:"oneof=admin user guest"`
	}

	tests := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{"valid role admin", "admin", false},
		{"valid role user", "user", false},
		{"valid role guest", "guest", false},
		{"invalid role", "superadmin", true},
		{"empty role", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TestStruct{Role: tt.role}
			err := Validate(s)

			if tt.wantErr && err == nil {
				t.Error("Validate() expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestValidate_NestedStruct(t *testing.T) {
	type Address struct {
		City string `validate:"required"`
	}

	type Person struct {
		Name    string  `validate:"required"`
		Address Address `validate:"required"`
	}

	tests := []struct {
		name    string
		person  Person
		wantErr bool
	}{
		{
			name: "valid nested struct",
			person: Person{
				Name:    "John",
				Address: Address{City: "NYC"},
			},
			wantErr: false,
		},
		{
			name: "missing nested field",
			person: Person{
				Name:    "John",
				Address: Address{City: ""},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.person)

			if tt.wantErr && err == nil {
				t.Error("Validate() expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestFormatFieldError(t *testing.T) {
	// This is tested indirectly through Validate tests
	// but we can add specific tests if needed
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
