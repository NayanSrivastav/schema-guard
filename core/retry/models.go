package retry

import (
	"context"
	"time"

	"github.com/schemaguard/schemaguard/core/validator"
)

// RetryStrategy defines how the LLM should be re-prompted upon validation failure
type RetryStrategy string

const (
	// SimpleRetry just repeats the exact same prompt
	SimpleRetry RetryStrategy = "simple"
	// ErrorAugmentedRetry appends the validation execution errors to the end of the prompt
	ErrorAugmentedRetry RetryStrategy = "error_augmented"
	// SchemaHintInjection appends both the errors and a hint containing the required schema
	SchemaHintInjection RetryStrategy = "schema_hint"
)

// TokenCost tracks LLM token usage and estimated cost
type TokenCost struct {
	TokensIn  int
	TokensOut int
	TotalCost float64 // Stored in USD ideally, requires knowing the model pricing
}

// RetryAttempt tracks metadata for a single LLM execution run
type RetryAttempt struct {
	AttemptNumber int
	Latency       time.Duration
	Tokens        TokenCost
	Success       bool
	ValidatorErrs int // Number of validation errors on this attempt
}

// RetryResult is the final output of the Engine execution loop
type RetryResult struct {
	FinalJSON string
	ValResult validator.ValidationResult
	Attempts  []RetryAttempt
	TotalCost TokenCost
	Model     string
}

// Config drives the engine execution rules
type Config struct {
	MaxRetries     int
	Strategy       RetryStrategy
	ValidationMode validator.ValidationMode
	SchemaJSON     string // The JSON Schema string standard
}

// LLMResponse normalizes an LLM output and its billing details
type LLMResponse struct {
	Content   string
	TokensIn  int
	TokensOut int
	ModelName string
}

// LLMClient represents an abstraction for generating LLM responses
type LLMClient interface {
	Generate(ctx context.Context, prompt string) (*LLMResponse, error)
}
