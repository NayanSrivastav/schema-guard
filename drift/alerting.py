import requests
import logging

class DriftAlerter:
    """
    Fires webhooks to team communication channels (Slack/Email/PagerDuty) 
    when structural validation distributions collapse over statistical thresholds.
    """
    def __init__(self, webhook_url: str):
        self.webhook_url = webhook_url
        self.logger = logging.getLogger("DriftAlerter")
        logging.basicConfig(level=logging.INFO)

    def limit_reached(self, schema_name: str, field_name: str, p_value: float, drift_type: str):
        msg = (
            f"🚨 *Schema Drift Detected!* 🚨\n"
            f"*Schema:* `{schema_name}`\n"
            f"*Field:* `{field_name}`\n"
            f"*Type:* `{drift_type}` Shift\n"
            f"*Confidence (p-value):* `{p_value:.5f}`\n\n"
            f"_The overall LLM output generation distributions have structurally shifted. Please investigate the model version or the prompt context chains._"
        )
        self.logger.warning(f"Drift Alert Triggered: {schema_name}.{field_name}")

        if not self.webhook_url:
            self.logger.info("No webhook URL configured. Suppressing network post.")
            return

        payload = {"text": msg}
        try:
            r = requests.post(self.webhook_url, json=payload, timeout=5)
            r.raise_for_status()
            self.logger.info("Slack webhook delivery successful")
        except Exception as e:
            self.logger.error(f"Failed to post Webhook: {e}")
