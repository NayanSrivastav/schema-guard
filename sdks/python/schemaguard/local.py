import json
import jsonschema
from typing import Dict, Any, Tuple

class LocalValidator:
    """
    Offline local MVP capability to test strict syntax validations without launching 
    the full Go binary. Useful for CI pipelines and quick testing natively via Python.
    """
    @staticmethod
    def coerce_heuristics(val: Any) -> Any:
        # Replicates Go's natively resilient coercion boundaries seamlessly 
        if isinstance(val, str):
            if val.isdigit():
                return int(val)
            try:
                return float(val)
            except ValueError:
                pass
            if val.lower() == "true":
                return True
            if val.lower() == "false":
                return False
            if (val.startswith('{') and val.endswith('}')) or (val.startswith('[') and val.endswith(']')):
                try:
                    parsed = json.loads(val)
                    return LocalValidator.coerce_heuristics(parsed)
                except Exception:
                    pass
            return val
        elif isinstance(val, dict):
            return {k: LocalValidator.coerce_heuristics(v) for k, v in val.items()}
        elif isinstance(val, list):
            return [LocalValidator.coerce_heuristics(v) for v in val]
        return val
    def __init__(self, schema: Dict[str, Any]):
        self.schema = schema

    def validate(self, payload: str) -> Tuple[bool, Any, list]:
        """
        Returns (is_valid, parsed_json_or_raw, [errors])
        """
        raw = payload
        
        # 1. Attempt to extract organic JSON block hidden within conversational markdown 
        if "```json" in payload:
            try:
                raw = payload.split("```json")[1].split("```")[0].strip()
            except IndexError:
                pass
        elif "```" in payload:
            try:
                raw = payload.split("```")[1].split("```")[0].strip()
            except IndexError:
                pass

        # 2. Syntax Check
        try:
            data = json.loads(raw)
        except json.JSONDecodeError as e:
            return False, raw, [f"Invalid JSON Format (Root Engine): {str(e)}"]

        # 3. Automatic Resilient Coercion matching Backend logic mapping
        data = self.coerce_heuristics(data)

        # 3. Ruleset Validation
        try:
            jsonschema.validate(instance=data, schema=self.schema)
            return True, data, []
        except jsonschema.exceptions.ValidationError as e:
            # Map precise path bounds (e.g. metadata.views object failed type constraints)
            path = ".".join([str(p) for p in e.path]) if e.path else "root"
            return False, data, [f"Field '{path}': {e.message}"]
