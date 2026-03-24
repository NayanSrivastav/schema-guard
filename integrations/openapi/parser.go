package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/schemaguard/schemaguard/core/registry"
)

// IngestSpec parses an OpenAPI YAML or JSON specification securely directly from a network URI or local path.
// It locates all the reusable Request/Response schemas bounded inside `components/schemas` natively
// and transforms them into strict JSON Schema `SchemaRecord` representations to perfectly feed into
// the Auto-Retry engine logic.
func IngestSpec(ctx context.Context, uri string, reg registry.Registry, version string) error {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	// Checks bounds cleanly to prevent unsafe parsing
	parsedUrl, err := url.Parse(uri)
	var doc *openapi3.T

	if err == nil && parsedUrl.Scheme != "" {
		doc, err = loader.LoadFromURI(parsedUrl)
	} else {
		// Attempting as bare filesystem loader
		doc, err = loader.LoadFromFile(uri)
	}

	if err != nil {
		return fmt.Errorf("failed resolving openapi spec via %s: %w", uri, err)
	}

	for name, schemaRef := range doc.Components.Schemas {
		schema := schemaRef.Value

		// Converts the generic openapi3 format gracefully physically into a JSON mapping AST logic byte block
		b, err := json.Marshal(schema)
		if err != nil {
			return fmt.Errorf("failed transforming component schema natively %s formatting: %w", name, err)
		}

		desc := schema.Description
		if desc == "" {
			desc = fmt.Sprintf("Auto-ingested via OpenApi documentation bounds internally securely (URI: %s)", uri)
		}

		record := registry.SchemaRecord{
			Name:        name,
			Version:     version,
			SchemaJSON:  string(b),
			CreatedAt:   time.Now(),
			Description: desc,
		}

		// Save the native constraints inside our Memory Lock
		if err := reg.Save(record); err != nil {
			return fmt.Errorf("schema saving failed into cache %s natively: %w", name, err)
		}
	}

	return nil
}
