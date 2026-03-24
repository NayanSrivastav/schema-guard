package validator

import (
	"testing"
)

var schemaStr = `{
	"type": "object",
	"properties": {
		"id": { "type": "integer" },
		"name": { "type": "string" },
		"isActive": { "type": "boolean" },
		"tags": {
			"type": "array",
			"items": { "type": "string" }
		},
		"metadata": {
			"type": "object",
			"properties": {
				"views": { "type": "number" }
			}
		}
	},
	"required": ["id", "name"]
}`

func TestValidator_StrictMode(t *testing.T) {
	v, err := NewValidator(schemaStr)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	tests := []struct {
		name       string
		input      string
		wantStatus string
	}{
		// 1-10: Basic functionality & required fields
		{"valid full object", `{"id": 1, "name": "test", "isActive": true, "tags": ["a"], "metadata": {"views": 100}}`, "PASS"},
		{"valid min object", `{"id": 1, "name": "test"}`, "PASS"},
		{"missing name", `{"id": 1}`, "FAIL"},
		{"missing id", `{"name": "test"}`, "FAIL"},
		{"empty object", `{}`, "FAIL"},
		{"null object", `null`, "FAIL"},
		{"empty string", `""`, "FAIL"},
		{"number instead of object", `123`, "FAIL"},
		{"array instead of object", `[]`, "FAIL"},
		{"malformed json", `{"id": 1, "name": "test"`, "FAIL"},
		// 11-20: Type mismatch validations
		{"string for integer", `{"id": "1", "name": "test"}`, "FAIL"},
		{"float for integer", `{"id": 1.5, "name": "test"}`, "FAIL"},
		{"boolean for integer", `{"id": true, "name": "test"}`, "FAIL"},
		{"integer for string", `{"id": 1, "name": 2}`, "FAIL"},
		{"null for string", `{"id": 1, "name": null}`, "FAIL"},
		{"string for boolean", `{"id": 1, "name": "test", "isActive": "true"}`, "FAIL"},
		{"string true for boolean", `{"id": 1, "name": "test", "isActive": "True"}`, "FAIL"},
		{"integer for boolean", `{"id": 1, "name": "test", "isActive": 1}`, "FAIL"},
		{"null for boolean", `{"id": 1, "name": "test", "isActive": null}`, "FAIL"},
		{"string for array", `{"id": 1, "name": "test", "tags": "a,b,c"}`, "FAIL"},
		// 21-25: Complex type mismatches
		{"array items wrong type", `{"id": 1, "name": "test", "tags": [1, 2, 3]}`, "FAIL"},
		{"null for array", `{"id": 1, "name": "test", "tags": null}`, "FAIL"},
		{"string for object", `{"id": 1, "name": "test", "metadata": "views: 1"}`, "FAIL"},
		{"nested string for number", `{"id": 1, "name": "test", "metadata": {"views": "100"}}`, "FAIL"},
		{"nested boolean for number", `{"id": 1, "name": "test", "metadata": {"views": true}}`, "FAIL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := v.Validate(tt.input, StrictMode)
			if res.Status != tt.wantStatus {
				t.Errorf("got status %q, want %q", res.Status, tt.wantStatus)
			}
		})
	}
}

func TestValidator_CoerceMode(t *testing.T) {
	v, err := NewValidator(schemaStr)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	tests := []struct {
		name       string
		input      string
		wantStatus string
	}{
		// 1-10: Coerced successes
		{"coerce string to int", `{"id": "1", "name": "test"}`, "PASS"},
		{"coerce float string to float", `{"id": 1, "name": "test", "metadata": {"views": "100.5"}}`, "PASS"},
		{"coerce int string to float", `{"id": 1, "name": "test", "metadata": {"views": "100"}}`, "PASS"},
		{"coerce string true to boolean", `{"id": 1, "name": "test", "isActive": "true"}`, "PASS"},
		{"coerce string false to boolean", `{"id": 1, "name": "test", "isActive": "false"}`, "PASS"},
		{"coerce string True to boolean", `{"id": 1, "name": "test", "isActive": "True"}`, "PASS"},
		{"coerce string FALSE to boolean", `{"id": 1, "name": "test", "isActive": "FALSE"}`, "PASS"},
		{"nested json string to object", `{"id": 1, "name": "test", "metadata": "{\"views\": 50}"}`, "PASS"},
		{"nested json string to array", `{"id": 1, "name": "test", "tags": "[\"a\", \"b\"]"}`, "PASS"},
		{"coerce root json extraction", "```\n{\"id\":\"1\", \"name\":\"test\"}\n```", "PASS"},
		// 11-15: Coerced misses (where coercion doesn't magically fix meaning)
		{"fail wrong int type", `{"id": "abc", "name": "test"}`, "FAIL"},
		{"fail missing required after coercion", `{"id": "1"}`, "FAIL"},
		{"fail array item conversion to wrong type", `{"id": 1, "name": "test", "tags": [1, 2]}`, "FAIL"}, 
		{"fail string to boolean invalid", `{"id": 1, "name": "test", "isActive": "yes"}`, "FAIL"},
		{"fail obj to array invalid", `{"id": 1, "name": "test", "tags": {"a":1}}`, "FAIL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := v.Validate(tt.input, CoerceMode)
			if res.Status != tt.wantStatus {
				t.Errorf("got status %q, want %q", res.Status, tt.wantStatus)
			}
		})
	}
}

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"json codeblock", "```json\n{\"id\": 1}\n```", `{"id": 1}`, false},
		{"json codeblock spacing", "```json \n {\"id\": 1}\n```", `{"id": 1}`, false},
		{"plan codeblock", "```\n{\"id\": 1}\n```", `{"id": 1}`, false},
		{"no codeblock curly", "some text \n {\"id\": 1} \n text", `{"id": 1}`, false},
		{"no codeblock square", "some text \n [{\"id\": 1}] \n text", `[{"id": 1}]`, false},
		{"pure json object", `{"id": 1}`, `{"id": 1}`, false},
		{"pure json array", `[1, 2, 3]`, `[1, 2, 3]`, false},
		{"malformed codeblock", "```json{\"id\": 1}", "", true},
		{"no brackets", `just some plain string without brackets`, "", true},
		{"json prefix suffix", `prefix {"key":"val"} suffix`, `{"key":"val"}`, false},
		{"array prefix suffix", `prefix [1,2] suffix`, `[1,2]`, false},
		{"multiple json structures takes first", `prefix {"id":1} mid {"id":2}`, `{"id":1}`, false},
		{"mixed brackets takes first sequence", `some text { [ ] } eof`, `{ [ ] }`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractJSON() = %q, want %q", got, tt.want)
			}
		})
	}
}
