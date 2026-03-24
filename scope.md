# AI Pipeline Reliability Tools Scope

## Overview
Two complementary tools forming a reliability platform for AI pipelines, targeting acquisition in 9–12 months.

**1. Product A — SchemaGuard:** Structured Output Validator & Schema Enforcer for LLM Pipelines (#21)
**2. Product B — AgentPost:** Agent Incident Response & Post-Mortem Automation (#20)

Build order: SchemaGuard is the fast MVP (4–6 weeks). AgentPost is built on top of its learnings and becomes the acquisition story.
*Brand name suggestion:* Reliably.ai or Pipecheck.dev

---

## PRODUCT A: SCHEMAGUARD
Structured Output Validator & Schema Enforcer for LLM Pipelines

### Problem Statement
LLMs returning structured data fail silently in ways difficult to debug (missing fields, schema drifts, uncontrolled retries). SchemaGuard acts as a runtime validation layer sitting directly between the LLM call and application logic.

### Core Features (MVP — Weeks 1–6)
*   **F1. Schema Definition Layer:** Define JSON/Pydantic/Zod schemas, versioning registry, CLI generator.
*   **F2. Runtime Validator:** Intercept LLM output, execute Strict/Coercion mode validation, return `PASS/FAIL/PARTIAL` with detailed field-level error mapping.
*   **F3. Auto-Retry Engine:** Error-augmented prompt engineering retry loops, circuit breaking, and retry token cost/latency tracking.
*   **F4. Schema Drift Detector:** Passive time-series output distribution monitoring and model version change behavior alerts.
*   **F5. Observability Dashboard:** Schema pass/fail rates, field failure heatmaps, CSV/OpenTelemetry exports.
*   **F6. SDK & Integration Layer:** Python/TypeScript SDKs, LangChain/LlamaIndex parsers.

### Post-MVP Features (Weeks 7–12)
*   **F7. Multi-model Schema Comparison:** A/B strictness tests across GPT-4o, Claude, Gemini.
*   **F8. Schema Marketplace:** Public registry of community schemas.
*   **F9. Prompt Template Linter:** LLM prompt heuristic improvements to maximize output schema conformance rate.
*   **F10. Compliance Mode:** PII scrubbing, SOC2 automated audit trails, strict GDPR no-raw-data modes.

---

### Tech Stack
*   **Core Runtime:** Go (Performance, low-latency, backend logic)
*   **SDKs:** Python + TypeScript (Primary LLM developer interfaces)
*   **Schema Store:** PostgreSQL + Redis cache
*   **Metrics Pipeline:** OpenTelemetry → ClickHouse
*   **Drift Detection:** Python (numpy/scipy)
*   **Dashboard:** React + Recharts / Grafana plugin
*   **Deployment:** Docker + Kubernetes Helm charts (Ready day 1)

---

### Week-by-Week Build Plan (SchemaGuard)
*   **Week 1:** Set up monorepo, CI/CD, linting, build core validator in Go, write 50+ unit tests.
*   **Week 2:** Auto-retry engine (F3) involving cost tracking, circuit breakers, and mock LLM integrations.
*   **Week 3:** Schema registry with versioning, LangChain output parser integration, Python SDK alpha on PyPI.
*   **Week 4:** Drift detection Python service, OpenTelemetry metrics pipeline, Slack webhook drift alerting.
*   **Week 5:** Dashboard MVP, Go REST API middleware bindings, Docker Compose orchestration for local environment runs.
*   **Week 6:** TypeScript SDK, comprehensive README, full quickstart guides, architecture docs, public demo video release.
