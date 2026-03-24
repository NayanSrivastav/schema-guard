package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/schemaguard/schemaguard/core/validator"
)

// Engine orchestrates LLM generation alongside dynamic retry/validation logic.
type Engine struct {
	Client   LLMClient
	Config   Config
	Val      *validator.Validator
	Breaker  *CircuitBreaker
}

// NewEngine constructs a schema-enforcing retry engine
func NewEngine(client LLMClient, circuit *CircuitBreaker, cfg Config) (*Engine, error) {
	v, err := validator.NewValidator(cfg.SchemaJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to load validator for engine: %w", err)
	}

	return &Engine{
		Client:  client,
		Config:  cfg,
		Val:     v,
		Breaker: circuit,
	}, nil
}

// Execute loops until MaxRetries are hit or a PASS status is achieved.
func (e *Engine) Execute(ctx context.Context, initialPrompt string) (*RetryResult, error) {
	if !e.Breaker.CanExecute() {
		return nil, ErrCircuitOpen
	}

	result := &RetryResult{
		TotalCost: TokenCost{},
	}

	currentPrompt := initialPrompt
	var lastErr error

	// To prevent infinite loops dynamically, enforce a global max loop bound + config boundary
	loops := e.Config.MaxRetries
	if loops < 1 {
		loops = 1
	}

	for i := 1; i <= loops; i++ {
		start := time.Now()
		
		// 1. Generate via LLM Client
		resp, llmErr := e.Client.Generate(ctx, currentPrompt)
		if llmErr != nil {
			e.Breaker.RecordFailure()
			lastErr = fmt.Errorf("LLM failure on attempt %d: %w", i, llmErr)
			continue
		}

		// 2. Add Token Costs
		attempt := RetryAttempt{
			AttemptNumber: i,
			Latency:       time.Since(start),
			Tokens: TokenCost{
				TokensIn:  resp.TokensIn,
				TokensOut: resp.TokensOut,
				// Rough $ estimate assuming generic average model cost
				TotalCost: float64(resp.TokensIn)*0.000005 + float64(resp.TokensOut)*0.000015,
			},
		}
		result.TotalCost.TokensIn += resp.TokensIn
		result.TotalCost.TokensOut += resp.TokensOut
		result.TotalCost.TotalCost += attempt.Tokens.TotalCost
		result.Model = resp.ModelName

		// 3. Validate response
		valRes, _ := e.Val.Validate(resp.Content, e.Config.ValidationMode)
		attempt.ValidatorErrs = len(valRes.Errors)

		if valRes.Status == "PASS" {
			attempt.Success = true
			result.Attempts = append(result.Attempts, attempt)
			result.FinalJSON = valRes.CoercedJSON
			if result.FinalJSON == "" { // fallback if strict mode passes without modifying payload
				resExtracted, errExt := validator.ExtractJSON(valRes.RawJSON)
				if errExt == nil {
					result.FinalJSON = resExtracted
				} else {
					result.FinalJSON = valRes.RawJSON
				}
			}
			result.ValResult = valRes
			
			e.Breaker.RecordSuccess()
			return result, nil // Early positive exit
		}

		// Validation Failed: Adjust Prompt via RetryStrategy and Loop over
		attempt.Success = false
		result.Attempts = append(result.Attempts, attempt)
		result.ValResult = valRes
		e.Breaker.RecordFailure()

		// Prepare next prompt if we aren't at the final iteration
		if i < loops {
			currentPrompt = e.buildRetryPrompt(initialPrompt, valRes)
		}
	}

	// Surpassed loops without early PASS
	if lastErr != nil && len(result.Attempts) == 0 {
		return result, lastErr
	}

	return result, fmt.Errorf("exhausted %d retries. final validation status: %s", loops, result.ValResult.Status)
}

func (e *Engine) buildRetryPrompt(base string, res validator.ValidationResult) string {
	if e.Config.Strategy == SimpleRetry {
		return base
	}

	errMsg := ""
	for _, vErr := range res.Errors {
		errMsg += fmt.Sprintf("- Field '%s': %s\n", vErr.Field, vErr.Message)
	}

	feedback := fmt.Sprintf("\n\nYour previous response failed structural validation. Please correct these errors:\n%s", errMsg)

	if e.Config.Strategy == SchemaHintInjection {
		hint := fmt.Sprintf("\nMake sure your output perfectly maps to this JSON Schema:\n```json\n%s\n```\n", e.Config.SchemaJSON)
		feedback += hint
	}

	return base + feedback
}
