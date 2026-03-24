package retry

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/schemaguard/schemaguard/core/validator"
)

// MockLLMClient allows deterministic simulation of LLM responses
type MockLLMClient struct {
	Responses   []LLMResponse
	ExecuteErrs []error
	Calls       int
}

func (m *MockLLMClient) Generate(ctx context.Context, prompt string) (*LLMResponse, error) {
	if m.Calls < len(m.ExecuteErrs) && m.ExecuteErrs[m.Calls] != nil {
		err := m.ExecuteErrs[m.Calls]
		m.Calls++
		return nil, err
	}
	if m.Calls < len(m.Responses) {
		resp := m.Responses[m.Calls]
		m.Calls++
		return &resp, nil
	}
	// Fallback to avoid panics if not enough mocks provided
	return &LLMResponse{Content: "{}", TokensIn: 10, TokensOut: 10, ModelName: "mock-fallback"}, nil
}

func TestEngine_SuccessOnFirstTry(t *testing.T) {
	schema := `{
		"type": "object",
		"properties": { "id": { "type": "integer" } },
		"required": ["id"]
	}`

	mock := &MockLLMClient{
		Responses: []LLMResponse{
			{Content: `{"id": 123}`, TokensIn: 50, TokensOut: 10, ModelName: "gpt-4o-mini"},
		},
	}

	cb := NewCircuitBreaker(5, time.Minute)
	eng, err := NewEngine(mock, cb, Config{
		MaxRetries:     3,
		Strategy:       ErrorAugmentedRetry,
		ValidationMode: validator.StrictMode,
		SchemaJSON:     schema,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, err := eng.Execute(context.Background(), "Extract an ID")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if res.ValResult.Status != "PASS" {
		t.Errorf("expected PASS, got %s", res.ValResult.Status)
	}
	if len(res.Attempts) != 1 {
		t.Errorf("expected 1 attempt, got %d", len(res.Attempts))
	}
	if mock.Calls != 1 {
		t.Errorf("expected 1 LLM call, got %d", mock.Calls)
	}
}

func TestEngine_RetryUntilSuccess(t *testing.T) {
	schema := `{
		"type": "object",
		"properties": { "id": { "type": "integer" } },
		"required": ["id"]
	}`

	mock := &MockLLMClient{
		Responses: []LLMResponse{
			{Content: `{"missing": "id field"}`, TokensIn: 50, TokensOut: 10, ModelName: "mock"}, // Fail
			{Content: `invalid generic text`, TokensIn: 60, TokensOut: 10, ModelName: "mock"}, // Fail
			{Content: `{"id": 456}`, TokensIn: 70, TokensOut: 10, ModelName: "mock"},          // Pass
		},
	}

	cb := NewCircuitBreaker(5, time.Minute)
	eng, err := NewEngine(mock, cb, Config{
		MaxRetries:     3,
		Strategy:       SchemaHintInjection,
		ValidationMode: validator.StrictMode,
		SchemaJSON:     schema,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, err := eng.Execute(context.Background(), "Get id")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if res.ValResult.Status != "PASS" {
		t.Errorf("expected PASS, got %s", res.ValResult.Status)
	}
	if len(res.Attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", len(res.Attempts))
	}
	
	totalIn := 50 + 60 + 70
	totalOut := 30
	if res.TotalCost.TokensIn != totalIn {
		t.Errorf("expected token tracking sum to match %d, got %d", totalIn, res.TotalCost.TokensIn)
	}
}

func TestCircuitBreaker_OpensAndBlocks(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)

	// Simulate failures
	cb.RecordFailure()
	if !cb.CanExecute() {
		t.Error("Expected to stay open under threshold")
	}

	cb.RecordFailure()
	if cb.CanExecute() {
		t.Error("Expected to block after hitting failure bound")
	}

	time.Sleep(100 * time.Millisecond)
	if !cb.CanExecute() {
		t.Error("Expected to re-open after timeout lapsed")
	}
}

func TestEngine_StrategyInjections(t *testing.T) {
	schema := `{
		"type": "object",
		"properties": { "id": { "type": "integer" } },
		"required": ["id"]
	}`

	// Just need a dummy that fails to inspect the prompt builder logic
	eng, _ := NewEngine(nil, nil, Config{
		MaxRetries:     3,
		Strategy:       SchemaHintInjection,
		ValidationMode: validator.StrictMode,
		SchemaJSON:     schema,
	})

	basePrompt := "Extract data"
	valErr := validator.ValidationResult{
		Status: "FAIL",
		Errors: []validator.ValidationError{{Field: "id", Message: "is missing"}},
	}
	
	prompt := eng.buildRetryPrompt(basePrompt, valErr)
	if !strings.Contains(prompt, "is missing") {
		t.Errorf("expected Prompt to contain validation errors, got %s", prompt)
	}
	if !strings.Contains(prompt, "Make sure your output perfectly maps to this JSON Schema") {
		t.Errorf("expected Prompt to contain Schema string hint")
	}
}
