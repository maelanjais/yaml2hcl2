package hcl

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	// Mettre une variable d'environnement pour le test
	os.Setenv("TEST_ENV_VAR", "my_test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	tests := []struct {
		name          string
		inputHCL      string
		expectedJSON  map[string]interface{}
		expectError   bool
		errorContains string
	}{
		{
			name: "Basic attributes evaluation",
			inputHCL: `
				foo = "bar"
				number = 42
				boolean = true
			`,
			expectedJSON: map[string]interface{}{
				"foo":     "bar",
				"number":  float64(42),
				"boolean": true,
			},
			expectError: false,
		},
		{
			name: "Function upper and lower evaluation",
			inputHCL: `
				uppercase = upper("hello")
				lowercase = lower("WORLD")
			`,
			expectedJSON: map[string]interface{}{
				"uppercase": "HELLO",
				"lowercase": "world",
			},
			expectError: false,
		},
		{
			name: "Function max and min evaluation",
			inputHCL: `
				maximum = max(10, 20, 5)
				minimum = min(10, 20, 5)
			`,
			expectedJSON: map[string]interface{}{
				"maximum": float64(20),
				"minimum": float64(5),
			},
			expectError: false,
		},
		{
			name: "Environment variables evaluation",
			inputHCL: `
				my_var = env.TEST_ENV_VAR
			`,
			expectedJSON: map[string]interface{}{
				"my_var": "my_test_value",
			},
			expectError: false,
		},
		{
			name: "Complex object evaluation",
			inputHCL: `
				tags = {
					name = upper("test-project")
					env  = "dev"
					count = max(1, 3)
				}
			`,
			expectedJSON: map[string]interface{}{
				"tags": map[string]interface{}{
					"name":  "TEST-PROJECT",
					"env":   "dev",
					"count": float64(3),
				},
			},
			expectError: false,
		},
		{
			name: "Syntax error in HCL",
			inputHCL: `
				foo = "bar
			`,
			expectError:   true,
			errorContains: "Erreur de parsing HCL",
		},
		{
			name: "Evaluation error in HCL (unknown function)",
			inputHCL: `
				foo = unknown_func("bar")
			`,
			expectError:   true,
			errorContains: "Erreur d'évaluation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultBytes, err := Evaluate([]byte(tt.inputHCL))

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errorContains), "Expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resultBytes)

			// Parse le JSON de retour pour comparer la structure
			var resultJSON map[string]interface{}
			err = json.Unmarshal(resultBytes, &resultJSON)
			assert.NoError(t, err, "Le résultat doit être un JSON valide")

			assert.Equal(t, tt.expectedJSON, resultJSON)
		})
	}
}
