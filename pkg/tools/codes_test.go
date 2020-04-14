package tools

import (
	"fmt"
	"regexp"
	"testing"
)

func TestGenerateCode(t *testing.T) {
	code := GenerateCode()

	if GenerateCode() == code {
		t.Error("Codes are not unique")
	}

	matched, _ := regexp.MatchString("^[A-Z0-9]{3}-[A-Z0-9]{3}$", code)
	if !matched {
		t.Errorf("Code has wrong format: %s", code)
	}
}

func TestFoo(t *testing.T) {
	for i := 0; i < 50; i++ {
		fmt.Println(GenerateCode())
	}
}
