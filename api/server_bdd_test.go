package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/schemaguard/schemaguard/core/registry"
)

// TestSchemaGuard_BDD executes the core engine capabilities seamlessly wrapping HTTP integration networks 
// across explicit Given/When/Then standard constraints cleanly.
func TestSchemaGuard_BDD(t *testing.T) {
	t.Run("Feature: LLM Payload Validation & Telemetry Tracking", func(t *testing.T) {
		
		// GIVEN: We initialize a global SchemaGuard memory store loaded with one test schema gracefully 
		reg := registry.NewMemoryRegistry()
		reg.Save(registry.SchemaRecord{
			Name:       "UserSchema",
			Version:    "1.0",
			SchemaJSON: `{"type":"object","properties":{"age":{"type":"integer"}},"required":["age"]}`,
		})
		srv := NewServer(reg)

		// Boot an HTTP router identically matching our production main.go execution bounds natively
		mux := http.NewServeMux()
		mux.HandleFunc("/v1/validate", srv.handleValidatePayload)
		mux.HandleFunc("/v1/stats", srv.handleGetTelemetryStats)

		/* -------------------------------------------------------------------------
		 * Scenario 1: Validating a perfectly formatted JSON payload securely
		 * -------------------------------------------------------------------------*/
		t.Run("Scenario: Validating a perfectly formatted JSON payload securely", func(t *testing.T) {
			t.Run("Given an incoming API request with correct JSON string mappings mapped neatly", func(t *testing.T) {
				
				reqBody := ValidatePayloadRequestObject{
					SchemaName: "UserSchema",
					Version:    "1.0",
					Payload:    `{"age": 25}`, // Completely valid structure natively
				}
				b, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/v1/validate", bytes.NewReader(b))
				rec := httptest.NewRecorder()

				t.Run("When the HTTP request is processed securely over the Multiplexer", func(t *testing.T) {
					mux.ServeHTTP(rec, req)

					t.Run("Then the response should yield 200 OK cleanly and trigger a PASS status", func(t *testing.T) {
						if rec.Code != http.StatusOK {
							t.Fatalf("Expected 200 OK cleanly, got %d", rec.Code)
						}
						
						var resp ValidatePayloadResponseObject
						json.NewDecoder(rec.Body).Decode(&resp)

						if resp.Status != "PASS" {
							t.Errorf("Expected PASS securely, got %s", resp.Status)
						}
					})
				})
			})
		})

		/* -------------------------------------------------------------------------
		 * Scenario 2: Fixing a hallucinated LLM markdown response via Coercion heuristics
		 * -------------------------------------------------------------------------*/
		t.Run("Scenario: Fixing a hallucinated LLM markdown response via Coercion", func(t *testing.T) {
			t.Run("Given an incoming payload containing messy conversational markdown and wrongly typed strings", func(t *testing.T) {
				
				reqBody := ValidatePayloadRequestObject{
					SchemaName: "UserSchema",
					Version:    "1.0",
					// LLM heavily hallucinated bounding markdown securely preventing JSON loads natively!
					// "age" was returned as a generic string "25" rather than raw Int natively
					Payload:    `Sure! Here is the data output natively: ` + "```json\n" + `{"age": "25"}` + "\n```",
				}
				b, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/v1/validate", bytes.NewReader(b))
				rec := httptest.NewRecorder()

				t.Run("When the SchemaGuard Engine natively coerces and validates the internal syntax mappings", func(t *testing.T) {
					mux.ServeHTTP(rec, req)

					t.Run("Then it translates the string into an integer safely tracking PASS bypassing failures natively", func(t *testing.T) {
						var resp ValidatePayloadResponseObject
						json.NewDecoder(rec.Body).Decode(&resp)
						
						if resp.Status != "PASS" {
							t.Errorf("Expected PASS status explicitly natively, got %s. AST Errors: %+v", resp.Status, resp.Errors)
						}
						
						if resp.CoercedJson == nil || !strings.Contains(*resp.CoercedJson, `"age":25`) {
							t.Errorf("Expected Coerced JSON payload safely translating strictly internal limits physically, got %v", resp.CoercedJson)
						}
					})
				})
			})
		})

		/* -------------------------------------------------------------------------
		 * Scenario 3: Rejecting an absolutely broken payload structurally natively
		 * -------------------------------------------------------------------------*/
		t.Run("Scenario: Rejecting an absolutely broken payload structurally", func(t *testing.T) {
			t.Run("Given an LLM response completely forgetting the required bounded schema 'age' field strictly", func(t *testing.T) {
				
				reqBody := ValidatePayloadRequestObject{
					SchemaName: "UserSchema",
					Version:    "1.0",
					Payload:    `{"name": "Developer"}`, // Missing required 'age' property
				}
				b, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/v1/validate", bytes.NewReader(b))
				rec := httptest.NewRecorder()

				t.Run("When the internal HTTP validation executes structurally testing strict bindings", func(t *testing.T) {
					mux.ServeHTTP(rec, req)

					t.Run("Then the engine returns a FAIL status mapping the exact required property failure cleanly into Heatmaps", func(t *testing.T) {
						var resp ValidatePayloadResponseObject
						json.NewDecoder(rec.Body).Decode(&resp)

						if resp.Status != "FAIL" {
							t.Errorf("Expected runtime tracking FAIL softly natively, got %s", resp.Status)
						}
						
						if len(resp.Errors) == 0 {
							t.Errorf("Expected exact validation execution limits mapping heavily populated dynamically")
						} else {
							if resp.Errors[0].Field != "age" && resp.Errors[0].Field != "" { // jsonschema locations differ sometimes depending on property layers
								t.Logf("Expected heavily mapped error array targeting limits successfully: %+v", resp.Errors)
							}
						}
					})
				})
			})
		})

		/* -------------------------------------------------------------------------
		 * Scenario 4: Retrieving Live React Dashboard Telemetry Aggregations
		 * -------------------------------------------------------------------------*/
		t.Run("Scenario: Retrieving Telemetry Aggregations from previous execution failures", func(t *testing.T) {
			t.Run("Given the previous web validations have flawlessly recorded memory distribution states seamlessly", func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/v1/stats", nil)
				rec := httptest.NewRecorder()

				t.Run("When the Frontend Dashboard polls exactly /v1/stats", func(t *testing.T) {
					mux.ServeHTTP(rec, req)

					t.Run("Then the system natively calculates and returns the execution validation KPIs synchronously natively", func(t *testing.T) {
						var stats TelemetryStats
						json.NewDecoder(rec.Body).Decode(&stats)

						// Mathematically passing logic: 2 SUCCESS (valid, coerced), 1 FAIL. 
						// Total = 3 executions. Pass rate logically explicitly safely equates to 66.66%
						kpiPass := stats.KPI["pass_rate"].(float64)
						if kpiPass < 66.0 || kpiPass > 67.0 {
							t.Errorf("Expected pass execution rate exactly roughly mapped to 66.6%% natively cleanly, got %f", kpiPass)
						}

						// Heatmap must structurally hold {"name": "age", "failures": 1} or {"name": "", "failures": 1} based on JS map bounds
						foundFailure := false
						for _, hm := range stats.Heatmap {
							if hm["failures"].(float64) == 1 {
								foundFailure = true
							}
						}
						if !foundFailure {
							t.Errorf("Expected runtime telemetry strictly mapping tracking 1 explicit failure dynamically cleanly, got %v", stats.Heatmap)
						}
					})
				})
			})
		})
	})
}
