package config

import (
	"errors"
	"testing"

	"github.com/abitofoldtom/marks/mocks"

	"github.com/abitofoldtom/marks/marks"
)

type mockProvider struct {
	getStringFn      func(string) string
	getStringCalled  bool
	setDefaultCalled bool
}

func (p *mockProvider) GetString(s string) string {
	p.getStringCalled = true
	return p.getStringFn(s)
}

func (p *mockProvider) SetDefault(s string, i interface{}) {
	p.setDefaultCalled = true
}

func mockGetString(s string) string {
	return "mock string"
}

func newMockProvider() *mockProvider {
	return &mockProvider{
		getStringFn: mockGetString,
	}
}

type mockValidator struct {
	validateFn       func(*marks.Config, []validationRule) error
	validateFnCalled bool
}

func mockValidate(*marks.Config, []validationRule) error {
	return nil
}

func newMockValidator() *mockValidator {
	return &mockValidator{
		validateFn: mockValidate,
	}
}

func (v *mockValidator) validate(config *marks.Config, rules []validationRule) error {
	v.validateFnCalled = true
	return v.validateFn(config, rules)
}

func TestLoadConfig(t *testing.T) {
	p := newMockProvider()
	v := newMockValidator()
	loader := loader{p, v}
	config, err := loader.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !p.setDefaultCalled {
		t.Fatal("provider SetDefault not called")
	}
	if !p.getStringCalled {
		t.Fatal("provider GetString not called")
	}
	if config.ContentPath != "mock string" {
		t.Fatal("ContentPath not equal to mock string")
	}
	if config.MarksYamlFileMode == 0 {
		t.Fatal("MarkYamlFileMode equal to 0")
	}
	if !v.validateFnCalled {
		t.Fatal("validator validate not called")
	}
}

func TestLoadConfigWithValidationError(t *testing.T) {
	p := newMockProvider()
	v := newMockValidator()
	v.validateFn = func(c *marks.Config, r []validationRule) error {
		return errors.New("Config failed validation")
	}
	loader := loader{p, v}
	config, err := loader.Load()
	if err == nil {
		t.Fatal("Should return error")
	}
	if config != nil {
		t.Fatal("Config should be nil")
	}
	if !p.setDefaultCalled {
		t.Fatal("provider SetDefault not called")
	}
	if !p.getStringCalled {
		t.Fatal("provider GetString not called")
	}
	if !v.validateFnCalled {
		t.Fatal("validator validate not called")
	}
}

func TestValidatorPass(t *testing.T) {
	rules := []validationRule{
		func(*marks.Config) error {
			return nil
		},
		func(*marks.Config) error {
			return nil
		},
	}
	config := mocks.NewConfig()
	v := &concreteValidator{}
	err := v.validate(config, rules)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestValidatorFail(t *testing.T) {
	rules := []validationRule{
		func(*marks.Config) error {
			return nil
		},
		func(*marks.Config) error {
			return errors.New("validation error")
		},
	}
	config := mocks.NewConfig()
	v := &concreteValidator{}
	err := v.validate(config, rules)
	if err == nil {
		t.Fatal("Should cause error")
	}
}
