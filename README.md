<div align="center">

# 🛡️ SchemaGuard
### Structured Output Validator & Schema Enforcer for LLM Pipelines
**Official Documentation:** [https://nayansrivastav.github.io/schema-guard/](https://nayansrivastav.github.io/schema-guard/)

[![CI Pipeline](https://github.com/NayanSrivastav/schema-guard/actions/workflows/ci.yml/badge.svg)](https://github.com/NayanSrivastav/schema-guard/actions)
![Coverage](https://img.shields.io/badge/Coverage-97.3%25-brightgreen.svg)
[![PyPI](https://img.shields.io/pypi/v/schema-guard-core)](https://pypi.org/project/schema-guard-core/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

</div>

LLMs returning structured data (JSON, XML) to backend systems fail silently in ways that are hard to debug. Schema streams silently drift across weights models, uncontrolled prompt hallucination loops cause massive financial Token cascades without latency limits, and there's fundamentally zero absolute operational visibility natively into reliability scaling.

SchemaGuard is an extremely fast runtime validation layer operating entirely locally directly physically placed directly between your LLM calls and architecture.

## 🚀 The Reliability Stack Architecture

*   **Go Engine:** High-performance core validation (`jsonschema`), circuit breakers cutting runaway generation bounds dynamically, telemetry token caching.
*   **Python SDK:** `pip install schema-guard-core` — offline validation or cloud telemetry mode for LangChain pipelines.
*   **React Dashboard:** Fast local Glassmorphic web observability maps field heatmaps natively into browser metrics directly.

## ⚡ Quickstart

### 1. Install the Python SDK

```bash
pip install schema-guard-core
```

### 2. Define your Schema

SchemaGuard uses standard JSON Schema for defining the structure and types of your expected output.

```python
# my_schema.py
customer_schema = {
    "type": "object",
    "properties": {
        "customer_id": {"type": "integer", "description": "Unique identifier for the customer"},
        "email": {"type": "string", "format": "email", "description": "Customer's email address"},
        "name": {"type": "string", "description": "Customer's full name"},
        "address": {
            "type": "object",
            "properties": {
                "street": {"type": "string"},
                "city": {"type": "string"},
                "zip_code": {"type": "string"}
            },
            "required": ["street", "city", "zip_code"]
        },
        "is_active": {"type": "boolean", "description": "Whether the customer account is active"}
    },
    "required": ["customer_id", "email", "name", "address", "is_active"]
}
```

### 3. Validate LLM Output

SchemaGuard offers two primary modes for validation:

#### Option A: Local Validation (Zero Latency, No Server Needed)

This mode performs validation entirely within your Python application, ideal for high-performance, low-latency scenarios.

```python
from schemaguard import LocalValidator
from my_schema import customer_schema

# Example LLM response (can be a string or a dict)
ai_response_valid = '{"customer_id": 123, "email": "test@example.com", "name": "John Doe", "address": {"street": "123 Main St", "city": "Anytown", "zip_code": "12345"}, "is_active": true}'
ai_response_invalid = '{"customer_id": "abc", "email": "invalid-email", "name": "Jane Doe", "address": {"street": "456 Oak Ave", "city": "Otherville"}, "is_active": "yes"}'

validator = LocalValidator(customer_schema)

# Validate a valid response
is_valid, data, errors = validator.validate(ai_response_valid)
if is_valid:
    print("✅ Clean payload:", data)
else:
    print("❌ Errors:", errors)

# Validate an invalid response
is_valid, data, errors = validator.validate(ai_response_invalid)
if is_valid:
    print("✅ Clean payload:", data)
else:
    print("❌ Errors:", errors)
```

#### Option B: Cloud Mode (Telemetry to Go Engine Dashboard)

For centralized monitoring and advanced features like circuit breakers and dynamic prompt morphing, use the `SchemaGuardClient`.

First, start the Go Engine locally:

```bash
docker compose up --build
```

The Go Validation Engine now sits on `localhost:8080` internally piping structural states and statistics mapped securely into your UI tracking Dashboard natively accessible physically at `http://localhost:5173`.

Then, use the client in your Python code:

```python
from schemaguard import SchemaGuardClient
from my_schema import customer_schema
import json

client = SchemaGuardClient(base_url="http://localhost:8080/v1")

# Register your schema with a name (only needs to be done once per schema)
# This allows the Go engine to reference it
client.register_schema(schema_name="CustomerProfile", schema_definition=customer_schema)

# Example LLM response
ai_response = '{"customer_id": 99, "email": "user@example.com", "name": "Alice", "address": {"street": "789 Pine St", "city": "Sometown", "zip_code": "67890"}, "is_active": false}'

result = client.validate(
    schema_name="CustomerProfile",
    payload=ai_response # Payload can be a JSON string or a Python dict
)

if result.is_valid:
    print("✅ Cloud validated payload:", result.validated_data)
else:
    print("❌ Cloud validation errors:", result.errors)
    print("Suggested fix:", result.suggested_fix) # Only available in Cloud Mode with certain validation modes
```

### 4. LangChain Integration

SchemaGuard seamlessly integrates with LangChain, allowing you to add robust validation to your LLM chains and agents.

```python
from schemaguard import LocalValidator
from my_schema import customer_schema
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import JsonOutputParser
from langchain_openai import ChatOpenAI

# 1. Define your schema (as above)
# from my_schema import customer_schema

# 2. Initialize LocalValidator
validator = LocalValidator(customer_schema)

# 3. Define your LLM and prompt
llm = ChatOpenAI(model="gpt-3.5-turbo", temperature=0)
parser = JsonOutputParser(pydantic_object=None) # We'll validate with SchemaGuard instead of Pydantic

prompt = ChatPromptTemplate.from_messages([
    ("system", "You are a helpful assistant that extracts customer information in JSON format."),
    ("human", "Extract customer details for a user named {name} with email {email} and address {address}. Make sure to include an active status."),
    ("human", "Format your response as a JSON object conforming to the following schema:\n{schema}")
]).partial(schema=json.dumps(customer_schema, indent=2))

chain = prompt | llm | parser

# 4. Invoke the chain and validate the output
raw_llm_output = chain.invoke({
    "name": "Bob Smith",
    "email": "bob.smith@example.com",
    "address": "101 Elm Street, Villagetown, 54321"
})

# The raw_llm_output is already parsed into a Python dict by JsonOutputParser
# Now, validate it with SchemaGuard
is_valid, validated_data, errors = validator.validate(raw_llm_output)

if is_valid:
    print("✅ LangChain output is valid:", validated_data)
else:
    print("❌ LangChain output validation errors:", errors)
    print("Raw LLM Output:", raw_llm_output)

# Example with an intentionally malformed prompt to show validation failure
malformed_prompt = ChatPromptTemplate.from_messages([
    ("system", "You are a helpful assistant that extracts customer information. Do not include an email or address."),
    ("human", "Extract customer details for a user named {name}. Make sure to include an active status."),
    ("human", "Format your response as a JSON object conforming to the following schema:\n{schema}")
]).partial(schema=json.dumps(customer_schema, indent=2))

malformed_chain = malformed_prompt | llm | parser

raw_llm_output_malformed = malformed_chain.invoke({
    "name": "Charlie Brown"
})

is_valid_malformed, validated_data_malformed, errors_malformed = validator.validate(raw_llm_output_malformed)

if is_valid_malformed:
    print("✅ LangChain malformed output is valid (unexpected!):", validated_data_malformed)
else:
    print("❌ LangChain malformed output validation errors:", errors_malformed)
    print("Raw LLM Output (malformed):", raw_llm_output_malformed)
```

## Core Validation Mechanics 🛠️

SchemaGuard dynamically operates upon 2 rulesets configurations logic.
* **StrictMode:** Performs rigorous typing bounds matches natively. 
* **CoerceMode:** Operates specifically applying heuristic translations fixing normal trivial hallucinations entirely directly internally across memory bounds. (e.g. converting nested numerical formatted strings into ints without breaking parsing).

```go
engine, _ := retry.NewEngine(llmClient, circuitBreaker, retry.Config{
    MaxRetries: 3,
    Strategy: retry.SchemaHintInjection, // Intelligently morphs the Prompt securely
    ValidationMode: validator.CoerceMode,
    SchemaJSON: schemaDefinitionStr,
})
```
