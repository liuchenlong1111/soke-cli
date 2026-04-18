package common

import (
	"context"
	"fmt"
	"io"

	"codeup.aliyun.com/5edbc121d1d1abe63b55f1c7/soke/soke-cli/internal/core"
)

type Flag struct {
	Name     string
	Type     string
	Default  string
	Desc     string
	Required bool
}

type Shortcut struct {
	Service     string
	Command     string
	Description string
	Risk        string
	Scopes      []string
	UserScopes  []string
	BotScopes   []string
	AuthTypes   []string
	Flags       []Flag
	HasFormat   bool
	Tips        []string
	DryRun      func(ctx context.Context, runtime *RuntimeContext) *DryRunAPI
	Validate    func(ctx context.Context, runtime *RuntimeContext) error
	Execute     func(ctx context.Context, runtime *RuntimeContext) error
}

func (s *Shortcut) ScopesForIdentity(identity string) []string {
	switch identity {
	case "user":
		if len(s.UserScopes) > 0 {
			return s.UserScopes
		}
	case "bot":
		if len(s.BotScopes) > 0 {
			return s.BotScopes
		}
	}
	return s.Scopes
}

type RuntimeContext struct {
	ctx    context.Context
	Config *core.CliConfig
	Cmd    interface{}
	Format string
	JqExpr string
}

func (r *RuntimeContext) Str(name string) string {
	return ""
}

func (r *RuntimeContext) Int(name string) int {
	return 0
}

func (r *RuntimeContext) Bool(name string) bool {
	return false
}

func (r *RuntimeContext) IsBot() bool {
	return false
}

func (r *RuntimeContext) CallAPI(method, path string, params, body interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (r *RuntimeContext) OutFormat(data map[string]interface{}, err error, fn func(io.Writer)) {
	if fn != nil {
		fn(io.Discard)
	}
}

type DryRunAPI struct {
	method      string
	path        string
	description string
	params      map[string]interface{}
	body        interface{}
	mode        string
}

func NewDryRunAPI() *DryRunAPI {
	return &DryRunAPI{}
}

func (d *DryRunAPI) GET(path string) *DryRunAPI {
	d.method = "GET"
	d.path = path
	return d
}

func (d *DryRunAPI) POST(path string) *DryRunAPI {
	d.method = "POST"
	d.path = path
	return d
}

func (d *DryRunAPI) Desc(desc string) *DryRunAPI {
	d.description = desc
	return d
}

func (d *DryRunAPI) Params(params map[string]interface{}) *DryRunAPI {
	d.params = params
	return d
}

func (d *DryRunAPI) Body(body interface{}) *DryRunAPI {
	d.body = body
	return d
}

func (d *DryRunAPI) Set(key, value string) *DryRunAPI {
	return d
}

func FlagErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
