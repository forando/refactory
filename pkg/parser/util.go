package parser

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
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

func WriterTokens(nativeTokens hclsyntax.Tokens) hclwrite.Tokens {
	// Ultimately we want a slice of token _pointers_, but since we can
	// predict how much memory we're going to devote to tokens we'll allocate
	// it all as a single flat buffer and thus give the GC less work to do.
	tokBuf := make([]hclwrite.Token, len(nativeTokens))
	var lastByteOffset int
	for i, mainToken := range nativeTokens {
		// Create a copy of the bytes so that we can mutate without
		// corrupting the original token stream.
		bytes := make([]byte, len(mainToken.Bytes))
		copy(bytes, mainToken.Bytes)

		tokBuf[i] = hclwrite.Token{
			Type:  mainToken.Type,
			Bytes: bytes,

			// We assume here that spaces are always ASCII spaces, since
			// that's what the scanner also assumes, and thus the number
			// of bytes skipped is also the number of space characters.
			SpacesBefore: mainToken.Range.Start.Byte - lastByteOffset,
		}

		lastByteOffset = mainToken.Range.End.Byte
	}

	// Now make a slice of pointers into the previous slice.
	ret := make(hclwrite.Tokens, len(tokBuf))
	for i := range ret {
		ret[i] = &tokBuf[i]
	}

	return ret
}
