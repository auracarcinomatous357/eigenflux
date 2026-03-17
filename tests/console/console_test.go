package console_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"eigenflux_server/pkg/config"
)

var consoleBaseURL = resolveConsoleBaseURL()

func resolveConsoleBaseURL() string {
	if v := strings.TrimSpace(os.Getenv("CONSOLE_API_BASE_URL")); v != "" {
		return strings.TrimSuffix(v, "/")
	}
	cfg := config.Load()
	return fmt.Sprintf("http://localhost:%d", cfg.ConsoleApiPort)
}

type ListAgentsData struct {
	Agents   []map[string]interface{} `json:"agents"`
	Total    int64                    `json:"total"`
	Page     int32                    `json:"page"`
	PageSize int32                    `json:"page_size"`
}

type ListAgentsResp struct {
	Code int32          `json:"code"`
	Msg  string         `json:"msg"`
	Data ListAgentsData `json:"data"`
}

type ListItemsData struct {
	Items    []map[string]interface{} `json:"items"`
	Total    int64                    `json:"total"`
	Page     int32                    `json:"page"`
	PageSize int32                    `json:"page_size"`
}

type ListItemsResp struct {
	Code int32         `json:"code"`
	Msg  string        `json:"msg"`
	Data ListItemsData `json:"data"`
}

type ListAgentImprItemsData struct {
	AgentID  string                   `json:"agent_id"`
	ItemIDs  []string                 `json:"item_ids"`
	GroupIDs []string                 `json:"group_ids"`
	URLs     []string                 `json:"urls"`
	Items    []map[string]interface{} `json:"items"`
}

type ListAgentImprItemsResp struct {
	Code int32                  `json:"code"`
	Msg  string                 `json:"msg"`
	Data ListAgentImprItemsData `json:"data"`
}

type MilestoneRuleInfo struct {
	RuleID          string `json:"rule_id"`
	MetricKey       string `json:"metric_key"`
	Threshold       int64  `json:"threshold"`
	RuleEnabled     bool   `json:"rule_enabled"`
	ContentTemplate string `json:"content_template"`
}

type ListMilestoneRulesData struct {
	Rules    []MilestoneRuleInfo `json:"rules"`
	Total    int64               `json:"total"`
	Page     int32               `json:"page"`
	PageSize int32               `json:"page_size"`
}

type ListMilestoneRulesResp struct {
	Code int32                  `json:"code"`
	Msg  string                 `json:"msg"`
	Data ListMilestoneRulesData `json:"data"`
}

type MilestoneRuleData struct {
	Rule MilestoneRuleInfo `json:"rule"`
}

type MilestoneRuleResp struct {
	Code int32              `json:"code"`
	Msg  string             `json:"msg"`
	Data *MilestoneRuleData `json:"data"`
}

type ReplaceMilestoneRuleData struct {
	OldRule MilestoneRuleInfo `json:"old_rule"`
	NewRule MilestoneRuleInfo `json:"new_rule"`
}

type ReplaceMilestoneRuleResp struct {
	Code int32                     `json:"code"`
	Msg  string                    `json:"msg"`
	Data *ReplaceMilestoneRuleData `json:"data"`
}

func TestConsoleListAgents(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/console/api/v1/agents?page=1&page_size=10", consoleBaseURL))
	if err != nil {
		t.Skipf("Console API not running: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result ListAgentsResp
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if result.Code != 0 {
		t.Fatalf("Expected code 0, got %d: %s", result.Code, result.Msg)
	}

	t.Logf("Listed %d agents (total: %d)", len(result.Data.Agents), result.Data.Total)
}

func TestConsoleListItems(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/console/api/v1/items?page=1&page_size=10", consoleBaseURL))
	if err != nil {
		t.Skipf("Console API not running: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result ListItemsResp
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if result.Code != 0 {
		t.Fatalf("Expected code 0, got %d: %s", result.Code, result.Msg)
	}

	t.Logf("Listed %d items (total: %d)", len(result.Data.Items), result.Data.Total)
}

func TestConsoleListAgentsPaginationParams(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/console/api/v1/agents?page=2&page_size=1", consoleBaseURL))
	if err != nil {
		t.Skipf("Console API not running: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result ListAgentsResp
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if result.Code != 0 {
		t.Fatalf("Expected code 0, got %d: %s", result.Code, result.Msg)
	}
	if result.Data.Page != 2 {
		t.Fatalf("Expected page=2, got %d", result.Data.Page)
	}
	if result.Data.PageSize != 1 {
		t.Fatalf("Expected page_size=1, got %d", result.Data.PageSize)
	}
	if len(result.Data.Agents) > 1 {
		t.Fatalf("Expected <=1 agents in page, got %d", len(result.Data.Agents))
	}
}

func TestConsoleListItemsPaginationParams(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/console/api/v1/items?page=2&page_size=1", consoleBaseURL))
	if err != nil {
		t.Skipf("Console API not running: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result ListItemsResp
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if result.Code != 0 {
		t.Fatalf("Expected code 0, got %d: %s", result.Code, result.Msg)
	}
	if result.Data.Page != 2 {
		t.Fatalf("Expected page=2, got %d", result.Data.Page)
	}
	if result.Data.PageSize != 1 {
		t.Fatalf("Expected page_size=1, got %d", result.Data.PageSize)
	}
	if len(result.Data.Items) > 1 {
		t.Fatalf("Expected <=1 items in page, got %d", len(result.Data.Items))
	}
}

func TestConsoleListAgentsWithFilters(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/console/api/v1/agents?page=1&page_size=10&name=test", consoleBaseURL))
	if err != nil {
		t.Skipf("Console API not running: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result ListAgentsResp
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if result.Code != 0 {
		t.Fatalf("Expected code 0, got %d: %s", result.Code, result.Msg)
	}

	t.Logf("Filtered agents by name 'test': %d results", len(result.Data.Agents))
}

func TestConsoleListItemsWithFilters(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/console/api/v1/items?page=1&page_size=10&status=3", consoleBaseURL))
	if err != nil {
		t.Skipf("Console API not running: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result ListItemsResp
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if result.Code != 0 {
		t.Fatalf("Expected code 0, got %d: %s", result.Code, result.Msg)
	}

	t.Logf("Filtered items by status 3 (completed): %d results", len(result.Data.Items))
}

func TestConsoleListAgentImprItems(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/console/api/v1/impr/items?agent_id=999999999", consoleBaseURL))
	if err != nil {
		t.Skipf("Console API not running: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		t.Skip("console api is running old binary without /console/api/v1/impr/items route")
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result ListAgentImprItemsResp
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if result.Code != 0 {
		t.Fatalf("Expected code 0, got %d: %s", result.Code, result.Msg)
	}
	if result.Data.AgentID != "999999999" {
		t.Fatalf("Expected agent_id=999999999, got %s", result.Data.AgentID)
	}
}

func TestConsoleMilestoneRulesFlow(t *testing.T) {
	seed := time.Now().UnixNano() % 1_000_000_000
	initialThreshold := 800_000_000 + seed
	replacementThreshold := initialThreshold + 1
	initialTemplate := fmt.Sprintf("Your Content \"{{.ItemSummary}}\" reached %d consumptions. Item Id {{.ItemID}}", initialThreshold)
	updatedTemplate := fmt.Sprintf("Updated template %d for {{.ItemID}} / {{.CounterValue}}", initialThreshold)
	replacedTemplate := fmt.Sprintf("Replacement template %d for {{.ItemSummary}}", replacementThreshold)

	createResp := doConsoleJSONRequest(t, http.MethodPost, "/console/api/v1/milestone-rules", map[string]interface{}{
		"metric_key":       "consumed",
		"threshold":        initialThreshold,
		"rule_enabled":     true,
		"content_template": initialTemplate,
	})
	var created MilestoneRuleResp
	mustDecodeConsoleResp(t, createResp, &created)
	if created.Code != 0 || created.Data == nil {
		t.Fatalf("create milestone rule failed: code=%d msg=%s", created.Code, created.Msg)
	}
	if created.Data.Rule.Threshold != initialThreshold {
		t.Fatalf("expected created threshold=%d, got %d", initialThreshold, created.Data.Rule.Threshold)
	}

	listResp := doConsoleRequest(t, http.MethodGet, fmt.Sprintf("/console/api/v1/milestone-rules?page=1&page_size=20&metric_key=consumed&rule_enabled=true"), nil)
	var listed ListMilestoneRulesResp
	mustDecodeConsoleResp(t, listResp, &listed)
	if listed.Code != 0 {
		t.Fatalf("list milestone rules failed: code=%d msg=%s", listed.Code, listed.Msg)
	}
	foundCreated := false
	for _, rule := range listed.Data.Rules {
		if rule.RuleID == created.Data.Rule.RuleID {
			foundCreated = true
			break
		}
	}
	if !foundCreated {
		t.Fatalf("created milestone rule %s not found in list response", created.Data.Rule.RuleID)
	}

	updateResp := doConsoleJSONRequest(t, http.MethodPut, "/console/api/v1/milestone-rules/"+created.Data.Rule.RuleID, map[string]interface{}{
		"rule_enabled":     false,
		"content_template": updatedTemplate,
	})
	var updated MilestoneRuleResp
	mustDecodeConsoleResp(t, updateResp, &updated)
	if updated.Code != 0 || updated.Data == nil {
		t.Fatalf("update milestone rule failed: code=%d msg=%s", updated.Code, updated.Msg)
	}
	if updated.Data.Rule.RuleEnabled {
		t.Fatalf("expected updated rule to be disabled")
	}
	if updated.Data.Rule.ContentTemplate != updatedTemplate {
		t.Fatalf("expected updated template=%q, got %q", updatedTemplate, updated.Data.Rule.ContentTemplate)
	}

	replaceResp := doConsoleJSONRequest(t, http.MethodPost, "/console/api/v1/milestone-rules/"+created.Data.Rule.RuleID+"/replace", map[string]interface{}{
		"metric_key":       "score_2",
		"threshold":        replacementThreshold,
		"rule_enabled":     true,
		"content_template": replacedTemplate,
	})
	var replaced ReplaceMilestoneRuleResp
	mustDecodeConsoleResp(t, replaceResp, &replaced)
	if replaced.Code != 0 || replaced.Data == nil {
		t.Fatalf("replace milestone rule failed: code=%d msg=%s", replaced.Code, replaced.Msg)
	}
	if replaced.Data.OldRule.RuleEnabled {
		t.Fatalf("expected old rule to be disabled after replace")
	}
	if replaced.Data.NewRule.MetricKey != "score_2" || replaced.Data.NewRule.Threshold != replacementThreshold {
		t.Fatalf("unexpected replacement rule: %+v", replaced.Data.NewRule)
	}

	cleanupResp := doConsoleJSONRequest(t, http.MethodPut, "/console/api/v1/milestone-rules/"+replaced.Data.NewRule.RuleID, map[string]interface{}{
		"rule_enabled": false,
	})
	var cleaned MilestoneRuleResp
	mustDecodeConsoleResp(t, cleanupResp, &cleaned)
	if cleaned.Code != 0 || cleaned.Data == nil {
		t.Fatalf("cleanup disable milestone rule failed: code=%d msg=%s", cleaned.Code, cleaned.Msg)
	}
	if cleaned.Data.Rule.RuleEnabled {
		t.Fatalf("expected cleanup rule to be disabled")
	}
}

func doConsoleRequest(t *testing.T, method, path string, body io.Reader) []byte {
	t.Helper()
	req, err := http.NewRequest(method, consoleBaseURL+path, body)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("Console API not running: %v", err)
		return nil
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body failed: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		t.Fatalf("unexpected 404 for %s %s: %s", method, path, string(payload))
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.StatusCode, string(payload))
	}
	return payload
}

func doConsoleJSONRequest(t *testing.T, method, path string, body interface{}) []byte {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request failed: %v", err)
	}
	return doConsoleRequest(t, method, path, bytes.NewReader(payload))
}

func mustDecodeConsoleResp(t *testing.T, payload []byte, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(payload, target); err != nil {
		t.Fatalf("failed to parse response: %v, body=%s", err, string(payload))
	}
}
