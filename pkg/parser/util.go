package parser

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"strings"
)

func getExpressionAsString(attr *hclwrite.Attribute) string {
	return removeQuotes(string(attr.Expr().BuildTokens(nil).Bytes()))
}

func removeQuotes(s string) string {
	out := strings.TrimSpace(s)
	if len(out) < 3 {
		return s
	}
	if out[len(out)-1] == '"' {
		out = out[:len(out)-1]
	}
	if out[0] == '"' {
		out = out[1:]
	}
	return out
}
