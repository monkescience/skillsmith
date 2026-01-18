package service

// AllTools returns all supported tools.
func AllTools() []Tool {
	return []Tool{ToolOpenCode, ToolClaude}
}

// AllScopes returns all supported scopes.
func AllScopes() []Scope {
	return []Scope{ScopeLocal, ScopeGlobal}
}

// ListTools returns all tools that have at least one compatible item in the registry.
func (s *Service) ListTools() []Tool {
	seen := make(map[Tool]bool)

	for _, item := range s.registry.Items {
		for _, tool := range item.Compatibility {
			seen[Tool(tool)] = true
		}
	}

	var tools []Tool

	for _, tool := range AllTools() {
		if seen[tool] {
			tools = append(tools, tool)
		}
	}

	return tools
}
