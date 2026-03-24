package validator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type ValidationMode string

const (
	StrictMode ValidationMode = "strict"
	CoerceMode ValidationMode = "coerce"
)

type ValidationResult struct {
	Status      string            `json:"status"` // "PASS", "FAIL", "PARTIAL"
	Errors      []ValidationError `json:"errors,omitempty"`
	CoercedJSON string            `json:"coerced_json,omitempty"`
	RawJSON     string            `json:"raw_json"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Validator struct {
	schema *jsonschema.Schema
}

// NewValidator builds a validator from a schema definition (as JSON string)
func NewValidator(schemaJSON string) (*Validator, error) {
	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", strings.NewReader(schemaJSON)); err != nil {
		return nil, fmt.Errorf("invalid schema format: %w", err)
	}
	sch, err := c.Compile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return &Validator{
		schema: sch,
	}, nil
}

// Validate takes LLM output (raw string or parsed JSON) and validates it against the schema.
func (v *Validator) Validate(input string, mode ValidationMode) (ValidationResult, error) {
	res := ValidationResult{
		Status:  "FAIL", // Default to FAIL until successful
		RawJSON: input,
		Errors:  []ValidationError{},
	}

	// 1. Try generic JSON unmarshal first to check if it's well-formed JSON
	var parsed interface{}
	err := json.Unmarshal([]byte(input), &parsed)
	if err != nil {
		// Possibly try to extract JSON from a markdown code block for the "raw string" support
		extracted, exErr := ExtractJSON(input)
		if exErr == nil {
			err = json.Unmarshal([]byte(extracted), &parsed)
			if err == nil {
				input = extracted
				res.RawJSON = extracted // Update to extracted raw for accuracy
			}
		}
		if err != nil {
			res.Errors = append(res.Errors, ValidationError{
				Field:   "root",
				Message: fmt.Sprintf("invalid JSON payload: %v", err),
			})
			return res, nil
		}
	}

	// For CoerceMode, we run a backup parse loop that tries heuristic coercion if validation fails initially
	var parsedCoerced interface{}
	if mode == CoerceMode {
	    // Apply heuristics to convert strings to numbers/bools where necessary
        // We will do a deep copy or marshaling clone so we don't pollute `parsed` if it succeeded
		b, _ := json.Marshal(parsed)
		var clone interface{}
		_ = json.Unmarshal(b, &clone)
		parsedCoerced = CoerceHeuristics(clone)
	} else {
	    parsedCoerced = parsed
	}

    // Try STRICT validation first
	validationErr := v.schema.Validate(parsed)
	if validationErr != nil {
		if mode == CoerceMode {
			// fallback to coerced value
			coerceErr := v.schema.Validate(parsedCoerced)
			if coerceErr != nil {
			    // Even coerced failed
			    errs := collectErrors(coerceErr)
			    res.Errors = errs
				res.Status = "FAIL"
			} else {
			    // Coerced Passed!
			    res.Status = "PASS"
			    parsed = parsedCoerced // switch to coerced
			}
		} else {
			// STRICT mode failure
			errs := collectErrors(validationErr)
			res.Errors = errs
			res.Status = "FAIL"
		}
	} else {
		res.Status = "PASS"
	}

	// 3. Serialize back if CoerceMode was used and something was modified (actually let's just serialize Coerced version)
	if mode == CoerceMode { //  and res.Status == "PASS" ?? actually the user may want to see what was coerced even on fail
		bb, _ := json.Marshal(parsed)
		res.CoercedJSON = string(bb)
	}

	return res, nil
}

func collectErrors(err error) []ValidationError {
    var validErrs []ValidationError
    if ve, ok := err.(*jsonschema.ValidationError); ok {
        return formatErrors(ve)
    }
    return append(validErrs, ValidationError{
        Field:   "schema",
        Message: err.Error(),
    })
}

func formatErrors(err *jsonschema.ValidationError) []ValidationError {
	var errs []ValidationError
	if len(err.Causes) == 0 {
	    loc := err.InstanceLocation
	    if loc == "" {
	        loc = "/"
	    }
		errs = append(errs, ValidationError{
			Field:   loc,
			Message: err.Message,
		})
		return errs
	}
	for _, cause := range err.Causes {
		errs = append(errs, formatErrors(cause)...)
	}
	return errs
}

// ExtractJSON tries to find a JSON block in a markdown text cleanly tracking bracket pairs
func ExtractJSON(input string) (string, error) {
	start := strings.Index(input, "```json")
	if start != -1 {
		input = input[start+7:]
	} else {
		start = strings.Index(input, "```")
		if start != -1 {
			input = input[start+3:]
		}
	}
	
	// Scan for the first array or object opening bounds natively
	firstCurly := strings.Index(input, "{")
	firstSquare := strings.Index(input, "[")
	
	if firstCurly == -1 && firstSquare == -1 {
		return "", fmt.Errorf("no json block bounds found gracefully natively")
	}
	
	idx := firstCurly
	if idx == -1 || (firstSquare != -1 && firstSquare < idx) {
		idx = firstSquare
	}
	
	// Brace counting structural evaluation
	if input[idx] == '{' {
		balance := 0
		for i := idx; i < len(input); i++ {
			if input[i] == '{' { balance++ }
			if input[i] == '}' { balance-- }
			if balance == 0 {
				return strings.TrimSpace(input[idx : i+1]), nil
			}
		}
	} else if input[idx] == '[' {
		balance := 0
		for i := idx; i < len(input); i++ {
			if input[i] == '[' { balance++ }
			if input[i] == ']' { balance-- }
			if balance == 0 {
				return strings.TrimSpace(input[idx : i+1]), nil
			}
		}
	}

	return "", fmt.Errorf("no properly terminating json block found completely gracefully")
}

// CoerceHeuristics walks the generic interface and converts typical LLM mistakes like string "123" to integer 123
func CoerceHeuristics(val interface{}) interface{} {
	switch v := val.(type) {
	case string:
		// Try to parse int
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return float64(i) // jsonschema typically handles numbers as float64 when unmarshaled generically
		}
		// Try to parse float
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
		// Try to parse bool
		lower := strings.ToLower(v)
		if lower == "true" {
			return true
		}
		if lower == "false" {
			return false
		}
		// Try json unmarshal if the string is json (nested json string bug)
		if (strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")) || (strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]")) {
		    var parsed interface{}
		    if err := json.Unmarshal([]byte(v), &parsed); err == nil {
		        return CoerceHeuristics(parsed)
		    }
		}
		return v
	case []interface{}:
		for i, el := range v {
			v[i] = CoerceHeuristics(el)
		}
		return v
	case map[string]interface{}:
		for key, mapVal := range v {
			v[key] = CoerceHeuristics(mapVal)
		}
		return v
	default:
		return v
	}
}
