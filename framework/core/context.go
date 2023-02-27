package core

// ExecutionContext type for the context of execution of the function including memory, function name and version...
type ExecutionContext struct {
	MemoryLimitInMB int    `json:"memoryLimitInMb"`
	FunctionName    string `json:"functionName"`
	FunctionVersion string `json:"functionVersion"`
}

// GetExecutionContext is used to create a new execution context and make it available, for offline testing thoses
// values are definied by default and does not affect functions performance.
func GetExecutionContext() ExecutionContext {
	return ExecutionContext{
		MemoryLimitInMB: 128,
		FunctionName:    "handler",
		FunctionVersion: "0.0.0",
	}
}
