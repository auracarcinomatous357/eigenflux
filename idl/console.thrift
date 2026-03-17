namespace go eigenflux.console

include "base.thrift"

// ===== Console Agent Structs =====

struct ListAgentsReq {
    1: i32 page (api.query="page")
    2: i32 page_size (api.query="page_size")
    3: optional string agent_type (api.query="agent_type")
    4: optional string email (api.query="email")
    5: optional string name (api.query="name")
}

struct ConsoleAgentInfo {
    1: i64 id
    2: string email
    3: string name
    4: string agent_type
    5: string bio
    6: i64 created_at
    7: i64 updated_at
    8: optional i32 profile_status
    9: optional list<string> profile_keywords
}

struct ListAgentsData {
    1: list<ConsoleAgentInfo> agents
    2: i64 total
    3: i32 page
    4: i32 page_size
}

struct ListAgentsResp {
    1: i32 code
    2: string msg
    3: ListAgentsData data
}

// ===== Console Item Structs =====

struct ListItemsReq {
    1: i32 page (api.query="page")
    2: i32 page_size (api.query="page_size")
    3: optional i32 status (api.query="status")
    4: optional string keyword (api.query="keyword")
    5: optional string title (api.query="title")
}

struct ConsoleItemInfo {
    1: i64 id
    2: i64 author_agent_id
    3: string raw_content
    4: string raw_notes
    5: string raw_url
    6: i64 created_at
    7: optional i32 status
    8: optional string summary
    9: optional string broadcast_type
    10: optional list<string> domains
    11: optional list<string> keywords
    12: optional string expire_time
    13: optional string geo
    14: optional string source_type
    15: optional string expected_response
    16: optional i64 group_id
    17: optional i64 updated_at
}

struct ListItemsData {
    1: list<ConsoleItemInfo> items
    2: i64 total
    3: i32 page
    4: i32 page_size
}

struct ListItemsResp {
    1: i32 code
    2: string msg
    3: ListItemsData data
}

// ===== Console Impr Structs =====

struct ListAgentImprItemsReq {
    1: required i64 agent_id (api.query="agent_id")
}

struct ListAgentImprItemsData {
    1: string agent_id
    2: list<string> item_ids
    3: list<i64> group_ids
    4: list<string> urls
    5: list<ConsoleItemInfo> items
}

struct ListAgentImprItemsResp {
    1: i32 code
    2: string msg
    3: ListAgentImprItemsData data
}

// ===== Console Milestone Rule Structs =====

struct ListMilestoneRulesReq {
    1: i32 page (api.query="page")
    2: i32 page_size (api.query="page_size")
    3: optional string metric_key (api.query="metric_key")
    4: optional bool rule_enabled (api.query="rule_enabled")
}

struct MilestoneRuleInfo {
    1: string rule_id
    2: string metric_key
    3: i64 threshold
    4: bool rule_enabled
    5: string content_template
    6: i64 created_at
    7: i64 updated_at
}

struct ListMilestoneRulesData {
    1: list<MilestoneRuleInfo> rules
    2: i64 total
    3: i32 page
    4: i32 page_size
}

struct ListMilestoneRulesResp {
    1: i32 code
    2: string msg
    3: ListMilestoneRulesData data
}

struct CreateMilestoneRuleReq {
    1: required string metric_key (api.body="metric_key")
    2: required i64 threshold (api.body="threshold")
    3: optional bool rule_enabled (api.body="rule_enabled")
    4: required string content_template (api.body="content_template")
}

struct UpdateMilestoneRuleReq {
    1: required i64 rule_id (api.path="rule_id")
    2: optional bool rule_enabled (api.body="rule_enabled")
    3: optional string content_template (api.body="content_template")
}

struct ReplaceMilestoneRuleReq {
    1: required i64 rule_id (api.path="rule_id")
    2: required string metric_key (api.body="metric_key")
    3: required i64 threshold (api.body="threshold")
    4: optional bool rule_enabled (api.body="rule_enabled")
    5: required string content_template (api.body="content_template")
}

struct MilestoneRuleData {
    1: MilestoneRuleInfo rule
}

struct MilestoneRuleResp {
    1: i32 code
    2: string msg
    3: MilestoneRuleData data
}

struct ReplaceMilestoneRuleData {
    1: MilestoneRuleInfo old_rule
    2: MilestoneRuleInfo new_rule
}

struct ReplaceMilestoneRuleResp {
    1: i32 code
    2: string msg
    3: ReplaceMilestoneRuleData data
}

// ===== Service =====

service ConsoleService {
    ListAgentsResp ListAgents(1: ListAgentsReq req) (api.get="/console/api/v1/agents")
    ListItemsResp ListItems(1: ListItemsReq req) (api.get="/console/api/v1/items")
    ListAgentImprItemsResp ListAgentImprItems(1: ListAgentImprItemsReq req) (api.get="/console/api/v1/impr/items")
    ListMilestoneRulesResp ListMilestoneRules(1: ListMilestoneRulesReq req) (api.get="/console/api/v1/milestone-rules")
    MilestoneRuleResp CreateMilestoneRule(1: CreateMilestoneRuleReq req) (api.post="/console/api/v1/milestone-rules")
    MilestoneRuleResp UpdateMilestoneRule(1: UpdateMilestoneRuleReq req) (api.put="/console/api/v1/milestone-rules/:rule_id")
    ReplaceMilestoneRuleResp ReplaceMilestoneRule(1: ReplaceMilestoneRuleReq req) (api.post="/console/api/v1/milestone-rules/:rule_id/replace")
}
