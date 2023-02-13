package core

import (
	"os"
)

// ExecutionContext - type for the context of execution of the function including memory, function name and version...
type ExecutionContext struct {
	MemoryLimitInMB int    `json:"memoryLimitInMb"`
	FunctionName    string `json:"functionName"`
	FunctionVersion string `json:"functionVersion"`
}

// GetExecutionContext - retrieve the execution context of the current function.
func GetExecutionContext() ExecutionContext {
	return ExecutionContext{
		MemoryLimitInMB: 128,
		FunctionName:    os.Getenv("SCW_APPLICATION_NAME"),
		FunctionVersion: os.Getenv("SCW_APPLICATION_VERSION"),
	}
}
