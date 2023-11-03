package json

import (
	"encoding/json"
	"testing"
)

// MockJSONHelper is a mock implementation of the Helper interface for testing purposes.
type MockJSONHelper struct {
	MarshalFunc   func(v interface{}) ([]byte, error)
	UnmarshalFunc func(data []byte, v interface{}) error
}

func (m *MockJSONHelper) Marshal(v interface{}) ([]byte, error) {
	if m.MarshalFunc != nil {
		return m.MarshalFunc(v)
	}
	return nil, nil
}

func (m *MockJSONHelper) Unmarshal(data []byte, v interface{}) error {
	if m.UnmarshalFunc != nil {
		return m.UnmarshalFunc(data, v)
	}
	return nil
}

func TestJSON_Marshal(t *testing.T) {
	// Create a JSON instance with a mock JSON helper
	mockHelper := &MockJSONHelper{
		MarshalFunc: func(v interface{}) ([]byte, error) {
			return []byte(`{"mocked":true}`), nil
		},
	}
	jsonHelper := NewJSON(mockHelper)

	// Test JSON.Marshal method
	result, err := jsonHelper.Marshal(map[string]interface{}{"key": "value"})
	if err != nil {
		t.Errorf("JSON.Marshal should not return an error, but got: %v", err)
	}
	expectedJSON := `{"mocked":true}`
	if string(result) != expectedJSON {
		t.Errorf("Expected JSON: %s, but got: %s", expectedJSON, string(result))
	}
}

func TestJSON_Unmarshal(t *testing.T) {
	// Create a JSON instance with a mock JSON helper
	mockHelper := &MockJSONHelper{
		UnmarshalFunc: func(data []byte, v interface{}) error {
			return json.Unmarshal([]byte(`{"mocked":true}`), v)
		},
	}
	jsonHelper := NewJSON(mockHelper)

	// Test JSON.Unmarshal method
	var output map[string]interface{}
	err := jsonHelper.Unmarshal([]byte(`{"key": "value"}`), &output)
	if err != nil {
		t.Errorf("JSON.Unmarshal should not return an error, but got: %v", err)
	}
	expectedOutput := map[string]interface{}{"mocked": true}
	if !compareJSON(output, expectedOutput) {
		t.Errorf("Expected output: %v, but got: %v", expectedOutput, output)
	}
}

func TestNativeJSON_Marshal(t *testing.T) {
	nativeJSONHelper := &NativeJSON{}

	// Test NativeJSON.Marshal method
	result, err := nativeJSONHelper.Marshal(map[string]interface{}{"key": "value"})
	if err != nil {
		t.Errorf("NativeJSON.Marshal should not return an error, but got: %v", err)
	}
	expectedJSON := `{"key":"value"}`
	if string(result) != expectedJSON {
		t.Errorf("Expected JSON: %s, but got: %s", expectedJSON, string(result))
	}
}

func TestNativeJSON_Unmarshal(t *testing.T) {
	nativeJSONHelper := &NativeJSON{}

	// Test NativeJSON.Unmarshal method
	var output map[string]interface{}
	err := nativeJSONHelper.Unmarshal([]byte(`{"key": "value"}`), &output)
	if err != nil {
		t.Errorf("NativeJSON.Unmarshal should not return an error, but got: %v", err)
	}
	expectedOutput := map[string]interface{}{"key": "value"}
	if !compareJSON(output, expectedOutput) {
		t.Errorf("Expected output: %v, but got: %v", expectedOutput, output)
	}
}

// compareJSON compares two JSON objects for equality.
func compareJSON(a, b map[string]interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}
