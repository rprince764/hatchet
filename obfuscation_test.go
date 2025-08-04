/*
 * Copyright 2022-present Kuei-chun Chen. All rights reserved.
 * obfuscation_test.go
 */

package hatchet

import (
	"regexp"
	"testing"
)

func TestObfuscateInt(t *testing.T) {
	// Initialize the Obfuscation struct
	o := &Obfuscation{
		intMap:      map[int]int{},
		Coefficient: 0.5,
	}

	// Test case 1: Obfuscating a new integer
	input1 := 10
	expectedOutput1 := 5
	actualOutput1 := o.ObfuscateInt(input1)
	if actualOutput1 != expectedOutput1 {
		t.Errorf("Expected %d but got %d for input %d", expectedOutput1, actualOutput1, input1)
	}

	// Test case 2: Obfuscating the same integer as test case 1
	// The function should return the cached result instead of recalculating it
	input2 := 10
	expectedOutput2 := expectedOutput1
	actualOutput2 := o.ObfuscateInt(input2)
	if actualOutput2 != expectedOutput2 {
		t.Errorf("Expected %d but got %d for input %d", expectedOutput2, actualOutput2, input2)
	}
}

func TestObfuscateNumber(t *testing.T) {
	// Initialize the Obfuscation struct
	o := &Obfuscation{
		numberMap:   map[string]float64{},
		Coefficient: 0.5,
	}

	// Test case 1: Obfuscating a new positive number
	input1 := 10.5
	expectedOutput1 := 5.25
	actualOutput1 := o.ObfuscateNumber(input1)
	if actualOutput1 != expectedOutput1 {
		t.Errorf("Expected %f but got %f for input %f", expectedOutput1, actualOutput1, input1)
	}

	// Test case 2: Obfuscating the same positive number as test case 1
	// The function should return the cached result instead of recalculating it
	input2 := 10.5
	expectedOutput2 := expectedOutput1
	actualOutput2 := o.ObfuscateNumber(input2)
	if actualOutput2 != expectedOutput2 {
		t.Errorf("Expected %f but got %f for input %f", expectedOutput2, actualOutput2, input2)
	}
}

func TestObfuscateCreditCardNo(t *testing.T) {
	o := &Obfuscation{}
	testCases := []struct {
		name  string
		input string
	}{
		{"with hyphens", "1234-5678-9012-3456"},
		{"with spaces", "1234 5678 9012 3456"},
		{"only digits", "1234567890123456"},
		{"no spaces or hyphens", "123456789012345"},
		{"empty", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.input == "" {
				if o.ObfuscateCreditCardNo(tc.input) != "" {
					t.Errorf("expected empty string for empty input")
				}
				return
			}
			obfuscated := o.ObfuscateCreditCardNo(tc.input)
			if len(obfuscated) != len(tc.input) {
				t.Errorf("expected length %d, but got %d", len(tc.input), len(obfuscated))
			}
			re := regexp.MustCompile(`[\d\s-]`)
			lastFour := re.ReplaceAllString(tc.input[len(tc.input)-4:], "")
			if obfuscated[len(obfuscated)-4:] != tc.input[len(tc.input)-4:] {
				t.Errorf("expected last 4 digits to be '%s', but got '%s'", lastFour, obfuscated[len(obfuscated)-4:])
			}
		})
	}
}

func TestObfuscateEmail(t *testing.T) {
	// Initialize the Obfuscation struct
	o := &Obfuscation{
		NameMap: make(map[string]string),
	}

	// Test case 1: Obfuscating a valid email address
	input1 := "john.doe@example.com"
	expectedOutput1Regex := regexp.MustCompile(`^[a-z]+@[a-z]+\.com$`)
	actualOutput1 := o.ObfuscateEmail(input1)
	if !expectedOutput1Regex.MatchString(actualOutput1) {
		t.Errorf("Expected obfuscated email to match pattern %s, but got %s", expectedOutput1Regex.String(), actualOutput1)
	}

	// Test case 2: Obfuscating an email address that is already obfuscated
	input2 := input1
	expectedOutput2 := actualOutput1
	actualOutput2 := o.ObfuscateEmail(input2)
	if actualOutput2 != expectedOutput2 {
		t.Errorf("Expected output to be %s but got %s for input %s", expectedOutput2, actualOutput2, input2)
	}
}

func TestObfuscateIP(t *testing.T) {
	// Initialize the Obfuscation struct
	o := &Obfuscation{
		IPMap: make(map[string]string),
	}

	// Test case 1: Obfuscating a valid IP address
	input1 := "192.168.0.1"
	expectedOutput1Regex := regexp.MustCompile(`^192\.[0-9]+\.[0-9]+\.[0-9]+$`)
	actualOutput1 := o.ObfuscateIP(input1)
	if !expectedOutput1Regex.MatchString(actualOutput1) {
		t.Errorf("Expected obfuscated IP to match pattern %s, but got %s", expectedOutput1Regex.String(), actualOutput1)
	}

	// Test case 2: Obfuscating the same IP address as in test case 1
	expectedOutput2 := actualOutput1
	actualOutput2 := o.ObfuscateIP(input1)
	if actualOutput2 != expectedOutput2 {
		t.Errorf("Expected output to be %s but got %s for input %s", expectedOutput2, actualOutput2, input1)
	}

	// Test case 5: Obfuscating an empty IP address
	input5 := ""
	expectedOutput5 := ""
	actualOutput5 := o.ObfuscateIP(input5)
	if actualOutput5 != expectedOutput5 {
		t.Errorf("Expected output to be %s but got %s for input %s", expectedOutput5, actualOutput5, input5)
	}
}

func TestObfuscateFQDN(t *testing.T) {
	// Initialize the Obfuscation struct
	o := &Obfuscation{
		NameMap: make(map[string]string),
	}

	// Test case 1: Obfuscating a valid FQDN with 2 parts
	input1 := "example.com"
	expectedOutputRegex := regexp.MustCompile(`([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}`)
	actualOutput1 := o.ObfuscateFQDN(input1)
	if !expectedOutputRegex.MatchString(actualOutput1) {
		t.Errorf("Expected obfuscated FQDN to match pattern %s, but got %s", expectedOutputRegex.String(), actualOutput1)
	}

	// Test case 2: Obfuscating a valid FQDN with more than 2 parts
	input2 := "www.example.co.uk"
	actualOutput2 := o.ObfuscateFQDN(input2)
	if !expectedOutputRegex.MatchString(actualOutput2) {
		t.Errorf("Expected obfuscated FQDN to match pattern %s, but got %s", expectedOutputRegex.String(), actualOutput2)
	}

	// Test case 3: Obfuscating an empty FQDN
	input3 := ""
	expectedOutput3 := ""
	actualOutput3 := o.ObfuscateFQDN(input3)
	if actualOutput3 != expectedOutput3 {
		t.Errorf("Expected output to be %s but got %s for input %s", expectedOutput3, actualOutput3, input3)
	}
}

func TestObfuscateNS(t *testing.T) {
	ptr := &Obfuscation{
		NameMap: make(map[string]string),
	}

	// Test case 1: Obfuscate a valid FQDN with two labels
	for _, ns := range []string{"example.com", "mail.example.com"} {
		obfuscated := ptr.ObfuscateNS(ns)
		if obfuscated == ns || !IsNamespace(obfuscated) {
			t.Errorf("ObfuscateNS(%q) returned %q, expected %q", ns, obfuscated, ns)
		}
	}

	// Test case 1: Obfuscate a valid FQDN with two labels
	for _, ns := range []string{"user@example.com", "user@mail.example.com"} {
		obfuscated := ptr.ObfuscateNS(ns)
		if obfuscated != ns {
			t.Errorf("ObfuscateNS(%q) returned %q, expected %q", ns, obfuscated, ns)
		}
	}
}

func TestObfuscateSSN(t *testing.T) {
	o := &Obfuscation{
		SSNMap: make(map[string]string),
	}
	testCases := []struct {
		name  string
		input string
	}{
		{"with hyphens", "123-45-6789"},
		{"only digits", "123456789"},
		{"invalid", "12345-6789"},
		{"empty", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.input == "" {
				if o.ObfuscateSSN(tc.input) != "" {
					t.Errorf("expected empty string for empty input")
				}
				return
			}
			if !IsSSN(tc.input) {
				if o.ObfuscateSSN(tc.input) != tc.input {
					t.Errorf("expected no change for invalid SSN")
				}
				return
			}
			obfuscated := o.ObfuscateSSN(tc.input)
			if obfuscated == tc.input {
				t.Errorf("expected obfuscated SSN to be different from input")
			}
			re := regexp.MustCompile(`^\d{3}-\d{2}-\d{4}$`)
			if !re.MatchString(obfuscated) {
				t.Errorf("expected obfuscated SSN to match pattern %s, but got %s", re.String(), obfuscated)
			}
		})
	}
}

func TestObfuscatePhoneNo(t *testing.T) {
	// Initialize the Obfuscation struct
	o := &Obfuscation{
		PhoneMap: make(map[string]string),
	}

	// Test case 1: Obfuscating a valid phone number with 10 digits
	input1 := "1234567890"
	expectedOutputRegex := regexp.MustCompile(`(?:\+?\d{1,3}[- ]?)?\d{10,14}|(\+\d{1,3}\s?)?\(\d{3}\)\s?\d{3}[- ]?\d{4}|\d{3}[- ]?\d{3}[- ]?\d{4}`)
	actualOutput1 := o.ObfuscatePhoneNo(input1)
	if !expectedOutputRegex.MatchString(actualOutput1) {
		t.Errorf("Expected obfuscated phone number to match pattern %s, but got %s", expectedOutputRegex.String(), actualOutput1)
	}

	// Test case 2: Obfuscating a valid phone number with 10 digits
	input2 := "123-456-7890"
	expectedOutput2Regex := regexp.MustCompile(`^(\d{3})[-\.\s]?(\d{3})[-\.\s]?(\d{4})$`)
	actualOutput2 := o.ObfuscatePhoneNo(input2)
	if !expectedOutputRegex.MatchString(actualOutput2) {
		t.Errorf("Expected obfuscated phone number to match pattern %s, but got %s", expectedOutput2Regex.String(), actualOutput2)
	}

	// Test case 3: Obfuscating an empty phone number
	input3 := ""
	expectedOutput3 := ""
	actualOutput3 := o.ObfuscatePhoneNo(input3)
	if actualOutput3 != expectedOutput3 {
		t.Errorf("Expected output to be %s but got %s for input %s", expectedOutput3, actualOutput3, input3)
	}

	// Test case 4: Obfuscating an empty phone number
	input4 := "+1 (123) 456-7890"
	actualOutput4 := o.ObfuscatePhoneNo(input4)
	if !expectedOutputRegex.MatchString(actualOutput4) {
		t.Errorf("Expected obfuscated phone number to match pattern %s, but got %s", expectedOutputRegex.String(), actualOutput4)
	}
}
