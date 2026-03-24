# SchemaGuard Python SDK

The official Python client for SchemaGuard — a runtime validation, monitoring, and coercion engine for protecting your LLM pipelines from bad JSON schemas.

## Installation

```bash
pip install schemaguard
```

## Quick Start (Offline Mode)

If you don't want to run the core Golang server, use the `LocalValidator` to catch schema drifts natively with exactly mapped formatting constraints dynamically directly in your Python code environments.

```python
from schemaguard import LocalValidator

schema = {
    "type": "object",
    "properties": { "id": {"type": "integer"} },
    "required": ["id"]
}

validator = LocalValidator(schema)
llm_response = '```json\n{"id": 123}\n```'

is_valid, data, errors = validator.validate(llm_response)
if is_valid:
    print("Clean format!", data["id"])
else:
    print("Errors:", errors)
```

## Production Mode (API Client)

In production, SchemaGuard manages schema distribution caching, active OpenTelemetry cost tracking, and circuit breakers directly locally or on your private clusters utilizing our blazingly fast Golang engine.

```python
from schemaguard import SchemaGuardClient

client = SchemaGuardClient(api_key="sg_...", base_url="http://localhost:8080/v1")

res = client.validate(schema_name="invoice_schema", version="latest", payload='{"total": "-10"}')
if not res.get("status") == "PASS":
    print(res["errors"]) # Maps dynamically cleanly for prompt re-injection strategies
```

## LangChain Integration

Using our native parser natively integrates into LangChain's existing RetryingOutputParser bounds without needing human configurations.

```python
from schemaguard_parser import SchemaGuardOutputParser

parser = SchemaGuardOutputParser(schema_dict=schema)
prompt = PromptTemplate(
    template="Extract invoice details.\n{format_instructions}\n{context}",
    input_variables=["context"],
    partial_variables={"format_instructions": parser.get_format_instructions()},
)

chain = prompt | llm | parser

# Automatically drops schema bounds into the prompt, extracts response payloads, validates constraints, and natively rejects outputs throwing `OutputParserException`.
```
