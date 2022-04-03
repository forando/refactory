package parser

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

type ExpressionParser struct {
	Context *hcl.EvalContext
}

func (p *ExpressionParser) ParseObjectExpr(name string, bytes []byte) (*map[string][]byte, error) {

	nativeExpr, diags := hclsyntax.ParseExpression(bytes, "dummy.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, errors.Errorf("%s expression cannot be parsed:\n%s", name, diags.Error())
	}

	switch expr := nativeExpr.(type) {
	case *hclsyntax.ObjectConsExpr:
		res := make(map[string][]byte)
		for _, item := range expr.Items {
			if keyVal, diags := item.KeyExpr.Value(p.Context); diags.HasErrors() {
				return nil, errors.Errorf("%s expression cannot be parsed:\n%s", name, diags.Error())
			} else {
				res[keyVal.AsString()] = item.ValueExpr.Range().SliceBytes(bytes)
			}
		}
		return &res, nil
	default:
		return nil, errors.Errorf("%s expression is not of object type", name)
	}
}

func (p *ExpressionParser) ParseArrayExpr(name string, bytes []byte) ([]string, error) {

	nativeExpr, diags := hclsyntax.ParseExpression(bytes, "dummy.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, errors.Errorf("%s expression cannot be parsed:\n%s", name, diags.Error())
	}

	var exprValue cty.Value
	if exprValue, diags = nativeExpr.Value(p.Context); diags != nil {
		return nil, errors.Errorf("%s expression cannot be parsed:\n%s", name, diags.Error())
	}

	res := make([]string, 0)

	for _, val := range exprValue.AsValueSlice() {
		res = append(res, val.AsString())
	}

	return res, nil
}
