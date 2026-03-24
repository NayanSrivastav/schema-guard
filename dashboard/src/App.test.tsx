import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import App from './App';

// BDD Mock Overrides isolating Frontend Execution from external Go API routing
global.fetch = vi.fn();

describe('Feature: SchemaGuard Telemetry Operational Dashboard', () => {
  
  beforeEach(() => {
    vi.resetAllMocks();
  });

  /* -------------------------------------------------------------------------
   * Scenario 1: Booting the Dashboard while fetching Live Engine Metrics
   * -------------------------------------------------------------------------*/
  describe('Scenario: Booting the Dashboard while fetching Live Engine Metrics', () => {
    
    describe('Given the Go validation engine has healthy timeseries configurations natively mapped', () => {
      beforeEach(() => {
        // Mocks successful API execution natively identical to actual REST endpoints
        (global.fetch as any).mockResolvedValue({
          json: async () => ({
            kpi: {
              pass_rate: 98.2,
              validations_month: 25000,
              cost_saved: 420.50
            },
            heatmap: [
              { name: "user.role", failures: 120 }
            ],
            timeseries: [
              { day: "Mon", success: 100, errors: 5 }
            ]
          })
        });
      });

      describe('When the React UI Component organically mounts to the DOM perfectly', () => {
        it('Then it correctly initializes loading states and populates the Operational Telemetry components natively', async () => {
          render(<App />);

          // Explicitly verifies the Initializing Screen loads perfectly without flashing directly initially
          expect(screen.getByText(/Initializing Telemetry Engine/i)).toBeDefined();

          // Wait natively for standard JavaScript Fetch API lifecycle constraints resolving safely
          await waitFor(() => {
            // Evaluates KPI Blocks organically matching exact syntax precision logic natively
            expect(screen.getByText(/98.2%/i)).toBeDefined();
            expect(screen.getByText(/25,000/i)).toBeDefined();
            expect(screen.getByText(/\$420.50/i)).toBeDefined();
            
            // Evaluates Top Heatmap Failing components structurally populated
            expect(screen.getByText(/user.role/i)).toBeDefined();
          });
        });
      });
    });
  });

  /* -------------------------------------------------------------------------
   * Scenario 2: Processing Empty Database Mappings without crashing the DOM
   * -------------------------------------------------------------------------*/
  describe('Scenario: Handling Empty Database Mappings natively seamlessly', () => {
    
    describe('Given the Go Engine has essentially completely zero data registered (On Day-1 Installation)', () => {
      beforeEach(() => {
        // Mocks completely barren API environment 
        (global.fetch as any).mockResolvedValue({
          json: async () => ({
            kpi: { pass_rate: 0, validations_month: 0, cost_saved: 0.00 },
            heatmap: [], 
            timeseries: []
          })
        });
      });

      describe('When the React Application executes parsing the empty HTTP responses safely cleanly', () => {
        it('Then the heatmap displays safe generic fallback boundaries dynamically entirely stopping DOM crashes', async () => {
          render(<App />);

          await waitFor(() => {
            // Verifies the exact UI Failsafe hooks we built previously natively execute organically
            expect(screen.getByText(/No errors recorded yet/i)).toBeDefined();
            expect(screen.getByText(/0%/i)).toBeDefined();
          });
        });
      });
    });
  });

});
