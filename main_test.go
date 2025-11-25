package main

import (
	"reflect"
	"testing"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hello"`, "hello"},
		{`"hello world"`, "hello world"},
		{`"hello\nworld"`, "hello\nworld"},
		{`"hello\tworld"`, "hello\tworld"},
		{`"quote: \"test\""`, `quote: "test"`},
		{`"unicode: \u0048\u0065\u006C\u006C\u006F"`, "unicode: Hello"},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.input, err)
			continue
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`123`, 123.0},
		{`-123`, -123.0},
		{`123.456`, 123.456},
		{`-123.456`, -123.456},
		{`1e10`, 1e10},
		{`1.5e-10`, 1.5e-10},
		{`0`, 0.0},
		{`-0`, -0.0},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.input, err)
			continue
		}
		if num, ok := result.(float64); !ok || num != tt.expected {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestParseBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`true`, true},
		{`false`, false},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.input, err)
			continue
		}
		if b, ok := result.(bool); !ok || b != tt.expected {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestParseNull(t *testing.T) {
	result, err := Parse(`null`)
	if err != nil {
		t.Errorf("Parse(null) error: %v", err)
	}
	if result != nil {
		t.Errorf("Parse(null) = %v, want nil", result)
	}
}

func TestParseArray(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`[]`, []interface{}{}},
		{`[1, 2, 3]`, []interface{}{1.0, 2.0, 3.0}},
		{`["a", "b", "c"]`, []interface{}{"a", "b", "c"}},
		{`[1, "two", true, null, false]`, []interface{}{1.0, "two", true, nil, false}},
		{`[[1, 2], [3, 4]]`, []interface{}{
			[]interface{}{1.0, 2.0},
			[]interface{}{3.0, 4.0},
		}},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.input, err)
			continue
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestParseObject(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		{`{}`, map[string]interface{}{}},
		{`{"name": "John"}`, map[string]interface{}{"name": "John"}},
		{`{"age": 30}`, map[string]interface{}{"age": 30.0}},
		{`{"active": true, "verified": false}`, map[string]interface{}{
			"active":   true,
			"verified": false,
		}},
		{`{"data": null}`, map[string]interface{}{"data": nil}},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.input, err)
			continue
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestParseNestedStructures(t *testing.T) {
	input := `{
		"user": {
			"name": "Alice",
			"age": 25,
			"scores": [95, 87, 92],
			"active": true
		},
		"metadata": {
			"created": "2024-01-01",
			"tags": ["important", "verified"]
		}
	}`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", result)
	}

	// Verify nested user object
	user := obj["user"].(map[string]interface{})
	if user["name"] != "Alice" {
		t.Errorf("Expected name Alice, got %v", user["name"])
	}
	if user["age"] != 25.0 {
		t.Errorf("Expected age 25, got %v", user["age"])
	}
	if user["active"] != true {
		t.Errorf("Expected active true, got %v", user["active"])
	}

	// Verify nested array
	scores := user["scores"].([]interface{})
	if len(scores) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(scores))
	}
}

func TestParseErrors(t *testing.T) {
	tests := []string{
		`{`,                  // Unterminated object
		`[`,                  // Unterminated array
		`"hello`,             // Unterminated string
		`{"key": }`,          // Missing value
		`{"key" "value"}`,    // Missing colon
		`[1 2 3]`,            // Missing commas
		`{123: "value"}`,     // Non-string key
		`invalid`,            // Invalid identifier
		`{"key": undefined}`, // Invalid value
	}

	for _, input := range tests {
		_, err := Parse(input)
		if err == nil {
			t.Errorf("Parse(%q) should have returned error", input)
		}
	}
}

func TestParseWhitespace(t *testing.T) {
	tests := []string{
		`  {"name": "John"}  `,
		`{"name"  :  "John"}`,
		`[  1  ,  2  ,  3  ]`,
		"{\n\t\"name\": \"John\"\n}",
	}

	for _, input := range tests {
		_, err := Parse(input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", input, err)
		}
	}
}

func BenchmarkParseSimpleObject(b *testing.B) {
	input := `{"name": "John", "age": 30, "active": true}`
	for i := 0; i < b.N; i++ {
		Parse(input)
	}
}

func BenchmarkParseArray(b *testing.B) {
	input := `[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]`
	for i := 0; i < b.N; i++ {
		Parse(input)
	}
}

func BenchmarkParseNestedStructure(b *testing.B) {
	input := `{"user": {"name": "Alice", "scores": [95, 87, 92]}, "metadata": {"tags": ["a", "b"]}}`
	for i := 0; i < b.N; i++ {
		Parse(input)
	}
}
