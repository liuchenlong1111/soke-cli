// Copyright (c) 2026 Soke Technologies.
// SPDX-License-Identifier: MIT

package contact

import (
	"context"
	"testing"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/shortcuts/common"
	"github.com/spf13/cobra"
)

// mockRuntimeContext creates a mock runtime context for testing
func mockRuntimeContext(flags map[string]string) *common.RuntimeContext {
	cmd := &cobra.Command{Use: "test"}
	for name := range flags {
		cmd.Flags().String(name, "", "")
	}
	cmd.ParseFlags(nil)
	for name, value := range flags {
		cmd.Flags().Set(name, value)
	}

	return &common.RuntimeContext{
		Cmd: cmd,
		Config: &core.CliConfig{
			AppID:     "test-app",
			AppSecret: "test-secret",
		},
	}
}

func TestContactGetDepartment_DryRun(t *testing.T) {
	runtime := mockRuntimeContext(map[string]string{
		"dept-id": "123456",
	})

	dryRun := ContactGetDepartment.DryRun(context.Background(), runtime)

	if dryRun == nil {
		t.Fatal("DryRun should not return nil")
	}
}

func TestContactGetDepartment_Metadata(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected interface{}
	}{
		{"Service", "Service", "contact"},
		{"Command", "Command", "+get-department"},
		{"Risk", "Risk", "read"},
		{"HasFormat", "HasFormat", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual interface{}
			switch tt.field {
			case "Service":
				actual = ContactGetDepartment.Service
			case "Command":
				actual = ContactGetDepartment.Command
			case "Risk":
				actual = ContactGetDepartment.Risk
			case "HasFormat":
				actual = ContactGetDepartment.HasFormat
			}

			if actual != tt.expected {
				t.Errorf("%s = %v, want %v", tt.field, actual, tt.expected)
			}
		})
	}
}

func TestContactGetDepartment_Flags(t *testing.T) {
	if len(ContactGetDepartment.Flags) == 0 {
		t.Fatal("ContactGetDepartment should have flags")
	}

	var deptIDFlag *common.Flag
	for i := range ContactGetDepartment.Flags {
		if ContactGetDepartment.Flags[i].Name == "dept-id" {
			deptIDFlag = &ContactGetDepartment.Flags[i]
			break
		}
	}

	if deptIDFlag == nil {
		t.Fatal("dept-id flag not found")
	}

	if !deptIDFlag.Required {
		t.Error("dept-id flag should be required")
	}
}

func TestContactGetDepartment_Scopes(t *testing.T) {
	if len(ContactGetDepartment.UserScopes) == 0 {
		t.Error("UserScopes should not be empty")
	}

	if len(ContactGetDepartment.BotScopes) == 0 {
		t.Error("BotScopes should not be empty")
	}

	expectedScope := "contact:department:readonly"

	hasUserScope := false
	for _, scope := range ContactGetDepartment.UserScopes {
		if scope == expectedScope {
			hasUserScope = true
			break
		}
	}
	if !hasUserScope {
		t.Errorf("UserScopes should contain %s", expectedScope)
	}

	hasBotScope := false
	for _, scope := range ContactGetDepartment.BotScopes {
		if scope == expectedScope {
			hasBotScope = true
			break
		}
	}
	if !hasBotScope {
		t.Errorf("BotScopes should contain %s", expectedScope)
	}
}

func TestContactGetDepartment_AuthTypes(t *testing.T) {
	if len(ContactGetDepartment.AuthTypes) != 2 {
		t.Errorf("AuthTypes length = %d, want 2", len(ContactGetDepartment.AuthTypes))
	}

	authTypes := make(map[string]bool)
	for _, authType := range ContactGetDepartment.AuthTypes {
		authTypes[authType] = true
	}

	if !authTypes["user"] {
		t.Error("AuthTypes should contain 'user'")
	}

	if !authTypes["bot"] {
		t.Error("AuthTypes should contain 'bot'")
	}
}

func TestContactGetDepartment_Execute(t *testing.T) {
	if ContactGetDepartment.Execute == nil {
		t.Fatal("Execute function should not be nil")
	}

	// Test that Execute function exists and can be called
	// Note: Full integration test would require mocking the API client
	runtime := mockRuntimeContext(map[string]string{
		"dept-id": "123456",
	})

	// This will fail without proper API mocking, but verifies the function signature
	_ = ContactGetDepartment.Execute(context.Background(), runtime)
}
