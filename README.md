<div align="center">

# 🛡️ SchemaGuard
### Structured Output Validator & Schema Enforcer for LLM Pipelines

[![CI Pipeline](https://github.com/NayanSrivastav/schema-guard/actions/workflows/ci.yml/badge.svg)](https://github.com/NayanSrivastav/schema-guard/actions)
[![PyPI](https://img.shields.io/pypi/v/schemaguard)](https://pypi.org/project/schemaguard/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

</div>

LLMs returning structured data (JSON, XML) to backend systems fail silently in ways that are hard to debug. Schema streams silently drift across weights models, uncontrolled prompt hallucination loops cause massive financial Token cascades without latency limits, and there's fundamentally zero absolute operational visibility natively into reliability scaling.

SchemaGuard is an extremely fast runtime validation layer operating entirely locally directly physically placed directly between your LLM calls and architecture.

## 🚀 The Reliability Stack Architecture

*   **Go Engine:** High-performance core validation (`jsonschema`), circuit breakers cutting runaway generation bounds dynamically, telemetry token caching.
*   **Python SDK:** `pip install schemaguard` acting securely mapped cleanly between Langchain pipelines safely.
*   **React Dashboard:** Fast local Glassmorphic web observability maps field heatmaps natively natively into browser metrics directly.

## ⚡ Quickstart

### 1. Start the System locally

```bash
docker compose up --build
```

The Go Validation Engine now sits on `localhost:8080` internally piping structural states and statistics mapped securely into your UI tracking Dashboard natively accessible physically at `http://localhost:5173`.

### 2. Validate your LLM generation pipelines natively via Python.

```bash
cd sdks/python
pip install .
```

```python
from schemaguard import SchemaGuardClient

# Directly interfaces your schema bounds into the cluster
client = SchemaGuardClient(api_key="your_api_key", base_url="http://localhost:8080/v1")
result = client.validate(
    schema_name="receipts", 
    version="v1", 
    payload='{"price": "10.00"}' 
)
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
