package factory

import (
	"github.com/forando/refactory/pkg/schema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func buildInlinePolicyTokens(document *schema.PolicyDocument) *hclwrite.Tokens {
	tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenEOF, Bytes: []byte("<<EOF")},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		{Type: hclsyntax.TokenOBrace, Bytes: []byte{'{'}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}
	tokens = append(tokens, *buildPolicyVersion()...)
	tokens = append(tokens, hclwrite.Tokens{
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("Statement")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenColon, Bytes: []byte(": ")},
		{Type: hclsyntax.TokenOBrack, Bytes: []byte{'['}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}...)
	for _, statement := range document.Statements {
		tokens = append(tokens, *buildPolicyStatement(&statement)...)
	}
	tokens = append(tokens, hclwrite.Tokens{
		{Type: hclsyntax.TokenCBrack, Bytes: []byte{']'}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		{Type: hclsyntax.TokenCBrace, Bytes: []byte{'}'}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		{Type: hclsyntax.TokenEOF, Bytes: []byte("EOF")},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}...)
	return &tokens
}

func buildPolicyStatement(statement *schema.Statement) *hclwrite.Tokens {
	tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenOBrace, Bytes: []byte{'{'}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}
	if len(statement.Sid) > 0 {
		tokens = append(tokens, *buildPolicySid(statement.Sid)...)
	}
	tokens = append(tokens, *buildPolicyEffect(statement.Effect)...)
	tokens = append(tokens, *buildPolicyAction(statement.Actions)...)
	tokens = append(tokens, *buildPolicyResource(statement.Resources)...)
	closingTokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenCBrace, Bytes: []byte{'}'}},
		{Type: hclsyntax.TokenComma, Bytes: []byte{','}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}
	tokens = append(tokens, closingTokens...)
	return &tokens
}

func buildPolicyVersion() *hclwrite.Tokens {
	return &hclwrite.Tokens{
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("Version")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenColon, Bytes: []byte(": ")},
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("2012-10-17")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenComma, Bytes: []byte{','}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}
}

func buildPolicySid(sid string) *hclwrite.Tokens {
	return &hclwrite.Tokens{
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("Sid")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenColon, Bytes: []byte(": ")},
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(sid)},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenComma, Bytes: []byte{','}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}
}

func buildPolicyEffect(effect string) *hclwrite.Tokens {
	if len(effect) == 0 {
		effect = schema.EffectAllow
	}
	return &hclwrite.Tokens{
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("Effect")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenColon, Bytes: []byte(": ")},
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte(effect)},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenComma, Bytes: []byte{','}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}
}

func buildPolicyAction(actions []string) *hclwrite.Tokens {
	tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("Action")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenColon, Bytes: []byte(": ")},
	}
	tokens = append(tokens, *buildArrayOfStrings(actions)...)
	return &tokens
}

func buildPolicyResource(resources []string) *hclwrite.Tokens {
	tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenStringLit, Bytes: []byte("Resource")},
		{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		{Type: hclsyntax.TokenColon, Bytes: []byte(": ")},
	}
	tokens = append(tokens, *buildArrayOfStrings(resources)...)
	return &tokens
}

func buildArrayOfStrings(values []string) *hclwrite.Tokens {
	tokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenOBrack, Bytes: []byte{'['}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}
	for _, action := range values {
		actionTokens := hclwrite.Tokens{
			{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
			{Type: hclsyntax.TokenStringLit, Bytes: []byte(action)},
			{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
			{Type: hclsyntax.TokenComma, Bytes: []byte{','}},
			{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		}
		tokens = append(tokens, actionTokens...)
	}
	closingTokens := hclwrite.Tokens{
		{Type: hclsyntax.TokenCBrack, Bytes: []byte{']'}},
		{Type: hclsyntax.TokenComma, Bytes: []byte{','}},
		{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}
	tokens = append(tokens, closingTokens...)
	return &tokens
}
