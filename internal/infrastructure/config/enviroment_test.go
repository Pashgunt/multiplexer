package config

import (
	"os"
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment()

	if env == nil {
		t.Error("NewEnvironment() returned nil")
	}
}

func TestEnvironment_Init(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		expectError bool
	}{
		{
			name: "successful init with .env file",
			setupEnv: func() {
				content := []byte("TEST_KEY=test_value\n")
				err := os.WriteFile(".env", content, 0644)

				if err != nil {
					t.Fatalf("Failed to create test .env file: %v", err)
				}

			},
			expectError: false,
		},
		{
			name:        "init with missing .env file",
			setupEnv:    func() {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer os.Remove(".env")

			env := NewEnvironment()
			err := env.Init()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestEnvironment_Get(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	env := NewEnvironment()
	result := env.Get("TEST_VAR")

	if result != "test_value" {
		t.Errorf("Get() = %v, want %v", result, "test_value")
	}

	if env.Get("NON_EXISTENT_VAR") != "" {
		t.Errorf("Get() for non-existent var = %v, want empty string", result)
	}
}

func TestEnvironment_extractEnvName(t *testing.T) {
	tests := []struct {
		name      string
		envString string
		want      string
	}{
		{
			name:      "valid env pattern",
			envString: "%env(DB_PASSWORD)%",
			want:      "DB_PASSWORD",
		},
		{
			name:      "valid env pattern with multiple",
			envString: "%env(API_KEY)%_suffix",
			want:      "API_KEY",
		},
		{
			name:      "no env pattern",
			envString: "plain text",
			want:      EmptyEnvName,
		},
		{
			name:      "empty string",
			envString: "",
			want:      EmptyEnvName,
		},
		{
			name:      "incomplete pattern",
			envString: "%env(INCOMPLETE",
			want:      EmptyEnvName,
		},
		{
			name:      "pattern with spaces",
			envString: "%env( TEST_VAR )%",
			want:      " TEST_VAR ",
		},
		{
			name:      "nested pattern",
			envString: "%env(OUTER)%_%env(INNER)%",
			want:      "OUTER",
		},
	}

	env := NewEnvironment()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := env.extractEnvName(tt.envString)

			if got != tt.want {
				t.Errorf("extractEnvName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironment_Replace(t *testing.T) {
	os.Setenv("TEST_HOST", "localhost")
	os.Setenv("TEST_PORT", "8080")
	os.Setenv("TEST_USER", "admin")
	os.Setenv("TEST_PASS", "secret")
	defer func() {
		os.Unsetenv("TEST_HOST")
		os.Unsetenv("TEST_PORT")
		os.Unsetenv("TEST_USER")
		os.Unsetenv("TEST_PASS")
	}()

	tests := []struct {
		name      string
		data      map[string]interface{}
		wantError bool
		validate  func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "successful replacement",
			data: map[string]interface{}{
				KeyTopics: TransportOption{
					"topic1": TransportOption{
						KeyOptions: TransportOption{
							"option1": TransportOption{
								"host": "%env(TEST_HOST)%",
								"port": "%env(TEST_PORT)%",
							},
							"option2": TransportOption{
								"user": "%env(TEST_USER)%",
								"pass": "%env(TEST_PASS)%",
							},
						},
					},
				},
			},
			wantError: false,
			validate: func(t *testing.T, data map[string]interface{}) {
				topics := data[KeyTopics].(TransportOption)
				topic1 := topics["topic1"].(TransportOption)
				options := topic1[KeyOptions].(TransportOption)

				option1 := options["option1"].(TransportOption)
				if option1["host"] != "localhost" {
					t.Errorf("host = %v, want localhost", option1["host"])
				}
				if option1["port"] != "8080" {
					t.Errorf("port = %v, want 8080", option1["port"])
				}

				option2 := options["option2"].(TransportOption)
				if option2["user"] != "admin" {
					t.Errorf("user = %v, want admin", option2["user"])
				}
				if option2["pass"] != "secret" {
					t.Errorf("pass = %v, want secret", option2["pass"])
				}
			},
		},
		{
			name: "missing topics key",
			data: map[string]interface{}{
				"other_key": "value",
			},
			wantError: true,
			validate:  nil,
		},
		{
			name: "empty options",
			data: map[string]interface{}{
				KeyTopics: TransportOption{
					"topic1": TransportOption{
						KeyOptions: TransportOption{},
					},
				},
			},
			wantError: false,
			validate: func(t *testing.T, data map[string]interface{}) {
				// No validation needed, just ensure no panic
			},
		},
		{
			name: "values without env patterns",
			data: map[string]interface{}{
				KeyTopics: TransportOption{
					"topic1": TransportOption{
						KeyOptions: TransportOption{
							"option1": TransportOption{
								"static": "static_value",
								"number": "123",
							},
						},
					},
				},
			},
			wantError: false,
			validate: func(t *testing.T, data map[string]interface{}) {
				topics := data[KeyTopics].(TransportOption)
				topic1 := topics["topic1"].(TransportOption)
				options := topic1[KeyOptions].(TransportOption)
				option1 := options["option1"].(TransportOption)

				if option1["static"] != "static_value" {
					t.Errorf("static = %v, want static_value", option1["static"])
				}
				if option1["number"] != "123" {
					t.Errorf("number = %v, want 123", option1["number"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			err := env.Replace(tt.data)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.wantError && tt.validate != nil {
				tt.validate(t, tt.data)
			}
		})
	}
}

func TestEnvironment_Replace_Panics(t *testing.T) {
	tests := []struct {
		name        string
		data        map[string]interface{}
		expectPanic bool
	}{
		{
			name: "invalid topics type",
			data: map[string]interface{}{
				KeyTopics: "not_a_map",
			},
			expectPanic: true,
		},
		{
			name: "invalid topic value type",
			data: map[string]interface{}{
				KeyTopics: TransportOption{
					"topic1": "not_a_map",
				},
			},
			expectPanic: true,
		},
		{
			name: "invalid options type",
			data: map[string]interface{}{
				KeyTopics: TransportOption{
					"topic1": TransportOption{
						KeyOptions: "not_a_map",
					},
				},
			},
			expectPanic: true,
		},
		{
			name: "invalid option value type",
			data: map[string]interface{}{
				KeyTopics: TransportOption{
					"topic1": TransportOption{
						KeyOptions: TransportOption{
							"option1": "not_a_map",
						},
					},
				},
			},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()

				if tt.expectPanic && r == nil {
					t.Error("Expected panic but got none")
				}

				if !tt.expectPanic && r != nil {
					t.Errorf("Expected no panic but got: %v", r)
				}
			}()

			env := NewEnvironment()
			env.Replace(tt.data)
		})
	}
}

func TestEnvironment_Replace_WithEmptyEnvName(t *testing.T) {
	os.Setenv("TEST_SKIP", "should_not_be_used")
	defer os.Unsetenv("TEST_SKIP")

	data := map[string]interface{}{
		KeyTopics: TransportOption{
			"topic1": TransportOption{
				KeyOptions: TransportOption{
					"option1": TransportOption{
						"value1": "%env()%",
						"value2": "%env(   )%",
						"value3": "plain text %env()%",
					},
				},
			},
		},
	}

	env := NewEnvironment()
	err := env.Replace(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	topics := data[KeyTopics].(TransportOption)
	topic1 := topics["topic1"].(TransportOption)
	options := topic1[KeyOptions].(TransportOption)
	option1 := options["option1"].(TransportOption)

	if option1["value1"] != "%env()%" {
		t.Errorf("value1 = %v, want %%env()%%", option1["value1"])
	}

	if option1["value2"] != "" {
		t.Errorf("value2 = %v, want empty", option1["value2"])
	}

	if option1["value3"] != "plain text %env()%" {
		t.Errorf("value3 = %v, want 'plain text %%env()%%'", option1["value3"])
	}
}

func TestEnvironment_Integration(t *testing.T) {
	envContent := []byte(`
DB_HOST=localhost
REDIS_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret123
API_KEY=abc123xyz
`)

	err := os.WriteFile(".env", envContent, 0644)

	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	defer os.Remove(".env")

	env := NewEnvironment()
	err = env.Init()

	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	data := map[string]interface{}{
		KeyTopics: TransportOption{
			"database": TransportOption{
				KeyOptions: TransportOption{
					"connection": TransportOption{
						"host":     "%env(DB_HOST)%",
						"port":     "%env(DB_PORT)%",
						"user":     "%env(DB_USER)%",
						"password": "%env(DB_PASSWORD)%",
					},
					"api": TransportOption{
						"key": "%env(API_KEY)%",
					},
				},
			},
			"cache": TransportOption{
				KeyOptions: TransportOption{
					"redis": TransportOption{
						"host": "%env(REDIS_HOST)%",
						"port": "6379",
					},
				},
			},
		},
	}

	err = env.Replace(data)

	if err != nil {
		t.Errorf("Replace failed: %v", err)
	}

	topics := data[KeyTopics].(TransportOption)

	dbTopic := topics["database"].(TransportOption)
	dbOptions := dbTopic[KeyOptions].(TransportOption)
	connection := dbOptions["connection"].(TransportOption)

	if connection["host"] != "localhost" {
		t.Errorf("DB host = %v, want localhost", connection["host"])
	}

	if connection["port"] != "5432" {
		t.Errorf("DB port = %v, want 5432", connection["port"])
	}

	if connection["user"] != "postgres" {
		t.Errorf("DB user = %v, want postgres", connection["user"])
	}

	if connection["password"] != "secret123" {
		t.Errorf("DB password = %v, want secret123", connection["password"])
	}

	api := dbOptions["api"].(TransportOption)

	if api["key"] != "abc123xyz" {
		t.Errorf("API key = %v, want abc123xyz", api["key"])
	}

	cacheTopic := topics["cache"].(TransportOption)
	cacheOptions := cacheTopic[KeyOptions].(TransportOption)
	redis := cacheOptions["redis"].(TransportOption)

	if redis["host"] != "localhost" {
		t.Errorf("Redis host = %v, want localhost", redis["host"])
	}

	if redis["port"] != "6379" {
		t.Errorf("Redis port = %v, want 6379", redis["port"])
	}
}

func BenchmarkEnvironment_Replace(b *testing.B) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	env := NewEnvironment()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		testData := map[string]interface{}{
			KeyTopics: TransportOption{
				"topic1": TransportOption{
					KeyOptions: TransportOption{
						"option1": TransportOption{
							"key1": "%env(TEST_VAR)%",
							"key2": "%env(TEST_VAR)%",
							"key3": "%env(TEST_VAR)%",
						},
					},
				},
			},
		}
		env.Replace(testData)
	}
}

func BenchmarkEnvironment_extractEnvName(b *testing.B) {
	env := NewEnvironment()
	envString := "%env(TEST_VAR)%"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		env.extractEnvName(envString)
	}
}
