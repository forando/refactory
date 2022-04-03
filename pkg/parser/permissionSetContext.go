package parser

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

func buildContext(permissionSetNames *map[string]string) *hcl.EvalContext {
	moduleMap := make(map[string]cty.Value)
	for key, val := range *permissionSetNames {
		moduleMap[key] = cty.MapVal(map[string]cty.Value{
			"permission_set_name": cty.StringVal(val),
		})
	}
	ctx := hcl.EvalContext{
		Variables: map[string]cty.Value{
			"module": cty.MapVal(moduleMap),
		},
	}

	return &ctx
}
