package util

import (
	"encoding/json"
	"testing"
)

func TestIsRawJSONValue(t *testing.T) {
	tests := []struct {
		val      string
		expected bool
	}{
		{"123", true},
		{"12.3", true},
		{"true", true},
		{"false", true},
		{"[1, 2]", true},
		{"{\"a\": 1}", true},
		{"string", false},
		{"\"string\"", false},
	}

	for _, test := range tests {
		if result := IsRawJSONValue(test.val); result != test.expected {
			t.Errorf("IsRawJSONValue(%q) = %v, expected %v", test.val, result, test.expected)
		}
	}
}

func TestInjectCIRulesUsesHTTPClientForRemoteRuleSets(t *testing.T) {
	input := `{
		"route": {
			"rule_set": [
				{"type": "remote", "tag": "remote", "download_detour": "proxy"},
				{"type": "local", "tag": "local", "path": "local.srs"}
			]
		}
	}`

	output, err := InjectCIRules(input)
	if err != nil {
		t.Fatalf("InjectCIRules() error = %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(output), &config); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}

	httpClients := config["http_clients"].([]interface{})
	httpClient := httpClients[len(httpClients)-1].(map[string]interface{})
	if httpClient["tag"] != "ci-rule-set-download" || httpClient["detour"] != "ci-direct-out" {
		t.Fatalf("unexpected CI HTTP client: %#v", httpClient)
	}

	route := config["route"].(map[string]interface{})
	if route["default_http_client"] != "ci-rule-set-download" {
		t.Errorf("default_http_client = %v", route["default_http_client"])
	}

	ruleSets := route["rule_set"].([]interface{})
	remoteRuleSet := ruleSets[0].(map[string]interface{})
	if _, exists := remoteRuleSet["download_detour"]; exists {
		t.Error("deprecated download_detour was not removed")
	}
	if remoteRuleSet["http_client"] != "ci-rule-set-download" {
		t.Errorf("remote http_client = %v", remoteRuleSet["http_client"])
	}

	localRuleSet := ruleSets[1].(map[string]interface{})
	if _, exists := localRuleSet["http_client"]; exists {
		t.Error("local rule-set must not receive an HTTP client")
	}
}
