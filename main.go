package main

import (
	"log"

	"github.com/schemaguard/schemaguard/api"
	"github.com/schemaguard/schemaguard/core/registry"
)

func main() {
	// Initializes native caching constraints natively natively natively
	reg := registry.NewMemoryRegistry()

	// Seeds realistic testing constraints seamlessly demonstrating full platform functionality dynamically
	customerSchema := `{
		"type": "object",
		"properties": {
			"customer_id": { "type": "integer" },
			"email": { "type": "string" },
			"metadata": {
				"type": "object",
				"properties": {
					"newsletter_opt_in": { "type": "boolean" }
				}
			}
		},
		"required": ["customer_id", "email"]
	}`

	err := reg.Save(registry.SchemaRecord{
		Name:        "CustomerProfile",
		Version:     "1.0",
		SchemaJSON:  customerSchema,
		Description: "Base validation bounds for customer payload outputs",
	})
	
	if err != nil {
		log.Fatalf("Failed natively seeding schema AST definitions: %v", err)
	}

	log.Println("Schema `CustomerProfile` actively registered to strict constraints successfully.")

	// Boots REST Multiplexer natively handling Dashboard metrics flawlessly without exceptions
	srv := api.NewServer(reg)
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("API SchemaGuard Engine pipeline faulted completely natively dynamically: %v", err)
	}
}
