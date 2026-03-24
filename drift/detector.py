import numpy as np
from scipy import stats
from typing import List, Tuple

class DriftDetector:
    """
    Passively monitors LLM output distributions and detects schema semantic drift.
    Numeric shifts use Kolmogorov-Smirnov test.
    Categorical shifts use Chi-Square checks.
    """
    def __init__(self, p_value_threshold: float = 0.05):
        self.threshold = p_value_threshold

    def check_numeric_drift(self, baseline: List[float], recent: List[float]) -> Tuple[bool, float]:
        """Runs a two-sample KS test to evaluate if numeric distributions have fundamentally shifted."""
        if len(baseline) < 5 or len(recent) < 5:
            return False, 1.0 # Minimum sample size not met for statistical significance

        # Compare cumulative distributions
        stat, p_value = stats.ks_2samp(baseline, recent)
        is_drifted = p_value < self.threshold
        return is_drifted, p_value

    def check_categorical_drift(self, baseline: List[str], recent: List[str]) -> Tuple[bool, float]:
        """Runs a relative Chi-Square approximation to track string token structural boundaries."""
        if len(baseline) < 5 or len(recent) < 5:
            return False, 1.0

        all_cats = set(baseline).union(set(recent))
        
        # Count occurrences
        base_counts = np.array([baseline.count(c) for c in all_cats]) + 1 # +1 Laplacian smoothing
        rec_counts = np.array([recent.count(c) for c in all_cats]) + 1
        
        # Normalize baseline to match recent window size dynamically
        base_freq = base_counts / sum(base_counts)
        rec_expected = base_freq * sum(rec_counts)

        stat, p_value = stats.chisquare(f_obs=rec_counts, f_exp=rec_expected)
        is_drifted = p_value < self.threshold
        return is_drifted, p_value
