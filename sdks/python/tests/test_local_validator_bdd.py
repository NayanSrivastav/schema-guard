import pytest
from schemaguard.local import LocalValidator

class TestSchemaGuardLocalValidation_BDD:
    """Feature: Offline Schema Verification & Markdown Extraction via Python SDK natively"""

    def setup_method(self):
        # GIVEN: The developer initializes a heavily strict schema offline cache natively seamlessly
        self.schema = {
            "type": "object",
            "properties": {
                "age": {"type": "integer"}
            },
            "required": ["age"]
        }
        self.validator = LocalValidator(self.schema)

    # -------------------------------------------------------------------------
    # Scenario 1: Validating a perfectly structured JSON payload securely
    # -------------------------------------------------------------------------
    def test_scenario_validating_perfect_json_payload(self):
        # GIVEN an incoming AI request containing correct JSON mapping bindings internally
        valid_llm_response = '{"age": 25}'

        # WHEN the Python SDK parses the generated response locally isolating constraints accurately
        is_valid, data, errors = self.validator.validate(valid_llm_response)

        # THEN the system should trigger a PASS cleanly and extract the integer functionally
        assert is_valid is True, "Expected valid execution natively mapped truthfully"
        assert data["age"] == 25
        assert len(errors) == 0

    # -------------------------------------------------------------------------
    # Scenario 2: Fixing a hallucinated LLM markdown response natively
    # -------------------------------------------------------------------------
    def test_scenario_fixing_hallucinated_markdown(self):
        # GIVEN a conversational payload containing messy markdown and implicitly wrong type formatting heavily natively
        hallucinated_response = "Here is the extracted information requested:\n```json\n{\"age\": \"25\"}\n```\nLet me know if you need anything else!"

        # WHEN the Python SDK dynamically truncates the markdown formatting attempting fallback logic internally
        is_valid, data, errors = self.validator.validate(hallucinated_response)

        # THEN it translates the parsed string `"25"` natively into the rigid integer `25` exactly bypassing LLM limits safely
        assert is_valid is True, f"Expected Coerced String mapping natively internally structurally. Errors: {errors}"
        assert data["age"] == 25, "Expected AST translation mapping the string inherently cleanly!"
        assert len(errors) == 0

    # -------------------------------------------------------------------------
    # Scenario 3: Rejecting an absolutely broken payload logically dynamically
    # -------------------------------------------------------------------------
    def test_scenario_rejecting_missing_parameters(self):
        # GIVEN an LLM output completely hallucinating past the critically required 'age' variable logically
        broken_response = '{"name": "Enterprise Engineer"}'

        # WHEN the internal Local SDK strict evaluation runs directly targeting formatting structurally
        is_valid, data, errors = self.validator.validate(broken_response)

        # THEN the engine traps a FAIL securely catching the exact AST path failure seamlessly mapping back error arrays logically
        assert is_valid is False
        assert len(errors) > 0
        
        # Verify that the missing 'age' property was correctly flagged inherently 
        error_message_contains_age = any("age" in err for err in errors)
        assert error_message_contains_age is True, "Expected specific error message dictating 'age' explicitly!"
