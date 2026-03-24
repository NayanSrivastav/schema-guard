package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/schemaguard/schemaguard/core/registry"
	"github.com/schemaguard/schemaguard/core/validator"
)

// Server acts as the core Web API Middleware matching routing logically to the Go Execution validation environments natively.
type Server struct {
	reg *registry.MemoryRegistry

	mu                sync.RWMutex
	passCount         int
	failCount         int
	fieldFailureCount map[string]int
	dailyStats        map[string]map[string]int // "Mon" -> {"success": 10, "errors": 2}
}

// NewServer boots the native Memory tracking metrics arrays smoothly mapping the telemetry values
func NewServer(r *registry.MemoryRegistry) *Server {
	return &Server{
		reg:               r,
		fieldFailureCount: make(map[string]int),
		dailyStats:        make(map[string]map[string]int),
	}
}

// Start hosts the HTTP server dynamically cleanly matching mapping endpoints natively without external SDK noise
func (s *Server) Start(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/validate", s.handleValidatePayload)
	mux.HandleFunc("/v1/stats", s.handleGetTelemetryStats)
	
	log.Printf("Starting SchemaGuard Validation Engine natively on %s", port)
	return http.ListenAndServe(port, mux)
}

func (s *Server) enableCORS(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (s *Server) handleValidatePayload(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w)
	if r.Method == http.MethodOptions { return }

	var req ValidatePayloadRequestObject
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body"})
		return
	}

	vers := req.Version
	if vers == "" || vers == "latest" {
		vers = "latest"
	}

	var schemaRec *registry.SchemaRecord
	var err error
	if vers == "latest" {
		schemaRec, err = s.reg.GetLatest(req.SchemaName)
	} else {
		schemaRec, err = s.reg.Get(req.SchemaName, vers)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Schema configuration not found in internal cache"})
		return
	}

	// Natively executes genuine validation Engine testing AST boundary trees natively natively natively
	val, err := validator.NewValidator(schemaRec.SchemaJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, _ := val.Validate(req.Payload, validator.CoerceMode) // Execute auto-fixing routines safely

	// Dynamically records raw metrics data arrays internally to broadcast real Dashboard updates safely
	s.mu.Lock()
	day := time.Now().Format("Mon")
	if s.dailyStats[day] == nil {
		s.dailyStats[day] = map[string]int{"success": 0, "errors": 0}
	}

	if res.Status == "PASS" {
		s.passCount++
		s.dailyStats[day]["success"]++
	} else {
		s.failCount++
		s.dailyStats[day]["errors"]++
		for _, e := range res.Errors {
			s.fieldFailureCount[e.Field]++
		}
	}
	s.mu.Unlock()

	var coerced *string
	if res.CoercedJSON != "" {
		c := res.CoercedJSON
		coerced = &c
	} else if res.Status == "PASS" {
		c := res.RawJSON
		coerced = &c
	}

	var errs []ValidationError
	for _, e := range res.Errors {
		errs = append(errs, ValidationError{Field: e.Field, Message: e.Message})
	}

	json.NewEncoder(w).Encode(ValidatePayloadResponseObject{
		Status:      res.Status,
		CoercedJson: coerced,
		RawJson:     res.RawJSON,
		Errors:      errs,
	})
}

func (s *Server) handleGetTelemetryStats(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w)
	if r.Method == http.MethodOptions { return }

	s.mu.RLock()
	defer s.mu.RUnlock()

	total := s.passCount + s.failCount
	passRate := 0.0
	if total > 0 {
		passRate = (float64(s.passCount) / float64(total)) * 100.0
	}

	heatmap := make([]map[string]interface{}, 0)
	for field, count := range s.fieldFailureCount {
		heatmap = append(heatmap, map[string]interface{}{"name": field, "failures": count})
	}

	timeseries := make([]map[string]interface{}, 0)
	for day, stats := range s.dailyStats {
		timeseries = append(timeseries, map[string]interface{}{
			"day": day, "success": stats["success"], "errors": stats["errors"],
		})
	}

    // Default zero-states effectively ensuring charts boot correctly gracefully natively flawlessly
	if len(timeseries) == 0 {
		timeseries = append(timeseries, map[string]interface{}{"day": time.Now().Format("Mon"), "success": s.passCount, "errors": s.failCount})
	}

	json.NewEncoder(w).Encode(TelemetryStats{
		KPI: map[string]interface{}{
			"pass_rate": passRate,
			"validations_month": total,
			"cost_saved": float64(s.failCount) * 0.04, // $0.04 token execution cascades aggressively prevented accurately per failure
		},
		Heatmap: heatmap,
		Timeseries: timeseries,
	})
}
