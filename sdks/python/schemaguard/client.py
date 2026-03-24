import requests
from typing import Dict, Any, Optional

class SchemaGuardClient:
    """
    Remote client for delegating LLM response validation to the Golang core platform.
    Designed to interact securely with the API engine metrics and circuit breaking layers.
    """
    def __init__(self, api_key: str, base_url: str = "http://localhost:8080/v1"):
        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self.session = requests.Session()
        self.session.headers.update({"Authorization": f"Bearer {self.api_key}"})

    def validate(self, schema_name: str, payload: str, version: Optional[str] = "latest") -> Dict[str, Any]:
        """
        Submits unstructured LLM payload text entirely securely targeting the Schema registry constraints.
        """
        response = self.session.post(
            f"{self.base_url}/validate",
            json={"schema_name": schema_name, "version": version, "payload": payload}
        )
        response.raise_for_status()
        return response.json()
