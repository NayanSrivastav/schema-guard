import sys
import os
import requests
import time

# Automatically append the sdks directory to the Python path so module imports work dynamically
sys.path.append(os.path.join(os.path.dirname(__file__), "sdks", "python"))

try:
    from schemaguard.client import SchemaGuardClient
except ModuleNotFoundError:
    print("Error: Could not import SchemaGuardClient. Please ensure 'requests' is installed: pip install requests")
    sys.exit(1)

def main():
    print("=========================================================")
    print("SchemaGuard: Live Telemetry Demonstration Pipeline")
    print("=========================================================\n")

    # Connect to the local Go API running on port 8085
    print("1. Initializing remote client mapping to localhost:8085...")
    client = SchemaGuardClient(api_key="local-testing-key", base_url="http://localhost:8085/v1")

    # Example 1: Perfect Execution Stream
    print("\n---------------------------------------------------------")
    print("Scenario A: Validating a perfectly clean AI Markdown response")
    print("---------------------------------------------------------")
    success_response = '```json\n{"customer_id": 99, "email": "nayan@schemaguard.com", "metadata": { "newsletter_opt_in": true }}\n```'
    print(f"Intercepted Payload:\n{success_response}\n")
    
    try:
        valid_result = client.validate(schema_name="CustomerProfile", payload=success_response)
        print(f"✅ Validation Status: {valid_result.get('status', 'FAIL')}")
        print(f"📊 Extracted Data: {valid_result.get('data')}")
    except requests.exceptions.ConnectionError:
        print("❌ ERROR: Connection Refused. Please ensure the Docker Go backend is running `docker compose up --build`!")
        return

    time.sleep(1)

    # Example 2: Complete Structure Hallucination
    print("\n---------------------------------------------------------")
    print("Scenario B: Synthesizing an AI Hallucination & Type Mismatch")
    print("---------------------------------------------------------")
    fail_response = "Here you go: {'customer_id': 'string_instead_of_int', 'metadata': {'newsletter_opt_in': 'missing_email'}}"
    print(f"Intercepted Payload:\n{fail_response}\n")
    
    invalid_result = client.validate(schema_name="CustomerProfile", payload=fail_response)
    print(f"❌ Validation Status: {invalid_result.get('status', 'FAIL')}")
    print("🚨 Captured Errors Logging to Dashboard:")
    for err in invalid_result.get("errors", []):
        print(f"   -> {err}")

    print("\n=========================================================")
    print("✨ Evaluation Traces Complete! ✨")
    print("=========================================================")
    print("➡️ Check your React Dashboard at http://localhost:5173 to witness the telemetry matrices populate instantly natively!")

if __name__ == "__main__":
    main()
