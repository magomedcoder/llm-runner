package mcpclient

import "sync/atomic"

type mcpRPCStats struct {
	listToolsOK, listToolsFail         atomic.Uint64
	listResourcesOK, listResourcesFail atomic.Uint64
	listPromptsOK, listPromptsFail     atomic.Uint64
	readResourceOK, readResourceFail   atomic.Uint64
	getPromptOK, getPromptFail         atomic.Uint64
	callToolOK, callToolFail           atomic.Uint64
	callToolMCPError                   atomic.Uint64
	probeOK, probeFail                 atomic.Uint64
}

var rpcStats mcpRPCStats

func recordListTools(err error) {
	if err != nil {
		rpcStats.listToolsFail.Add(1)
	} else {
		rpcStats.listToolsOK.Add(1)
	}
}

func recordListResources(err error) {
	if err != nil {
		rpcStats.listResourcesFail.Add(1)
	} else {
		rpcStats.listResourcesOK.Add(1)
	}
}

func recordListPrompts(err error) {
	if err != nil {
		rpcStats.listPromptsFail.Add(1)
	} else {
		rpcStats.listPromptsOK.Add(1)
	}
}

func recordReadResource(err error) {
	if err != nil {
		rpcStats.readResourceFail.Add(1)
	} else {
		rpcStats.readResourceOK.Add(1)
	}
}

func recordGetPrompt(err error) {
	if err != nil {
		rpcStats.getPromptFail.Add(1)
	} else {
		rpcStats.getPromptOK.Add(1)
	}
}

func recordCallToolTransportErr() {
	rpcStats.callToolFail.Add(1)
}

func recordCallToolMCPError() {
	rpcStats.callToolMCPError.Add(1)
}

func recordCallToolOK() {
	rpcStats.callToolOK.Add(1)
}

func recordProbe(err error) {
	if err != nil {
		rpcStats.probeFail.Add(1)
	} else {
		rpcStats.probeOK.Add(1)
	}
}

func MCPCountersMap() map[string]uint64 {
	return map[string]uint64{
		"list_tools_ok":       rpcStats.listToolsOK.Load(),
		"list_tools_fail":     rpcStats.listToolsFail.Load(),
		"list_resources_ok":   rpcStats.listResourcesOK.Load(),
		"list_resources_fail": rpcStats.listResourcesFail.Load(),
		"list_prompts_ok":     rpcStats.listPromptsOK.Load(),
		"list_prompts_fail":   rpcStats.listPromptsFail.Load(),
		"read_resource_ok":    rpcStats.readResourceOK.Load(),
		"read_resource_fail":  rpcStats.readResourceFail.Load(),
		"get_prompt_ok":       rpcStats.getPromptOK.Load(),
		"get_prompt_fail":     rpcStats.getPromptFail.Load(),
		"call_tool_ok":        rpcStats.callToolOK.Load(),
		"call_tool_fail":      rpcStats.callToolFail.Load(),
		"call_tool_mcp_error": rpcStats.callToolMCPError.Load(),
		"probe_ok":            rpcStats.probeOK.Load(),
		"probe_fail":          rpcStats.probeFail.Load(),
	}
}
