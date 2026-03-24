import json
from typing import Any, Dict

try:
    from langchain_core.output_parsers import BaseOutputParser
    from langchain_core.exceptions import OutputParserException
except ImportError:
    raise ImportError("SchemaGuard integration requires 'langchain_core'. Please run `pip install langchain-core`.")

from schemaguard.local import LocalValidator

class SchemaGuardOutputParser(BaseOutputParser):
    """
    Native LangChain Output Parser wrapper. 
    Intercepts the generated text response, isolates organic JSON syntax, targets it
    against the SchemaGuard LocalValidator routines, and structurally rejects
    faulty formatting directly feeding Langchain RetryEngines.
    """
    schema_dict: Dict[str, Any]
    
    @property
    def _type(self) -> str:
        return "schema_guard_parser"

    def parse(self, text: str) -> Any:
        validator = LocalValidator(self.schema_dict)
        is_valid, parsed, errors = validator.validate(text)
        
        if not is_valid:
            error_msg = "\n".join(errors)
            # This triggers LangChain's native RetryWithErrorOutputParser flow securely
            raise OutputParserException(f"SchemaGuard Validation Failed:\n{error_msg}")
            
        return parsed

    def get_format_instructions(self) -> str:
        """
        Dynamically passes JSON-Schema hints structurally directly back to prompts natively.
        """
        schema_dump = json.dumps(self.schema_dict, indent=2)
        return (
            "You must format your exact output as a JSON blob. Do not append conversational text outside the code block."
            f"The generated JSON must perfectly map to the following schema structurally and typing boundaries:\n"
            f"```json\n{schema_dump}\n```\n"
        )
