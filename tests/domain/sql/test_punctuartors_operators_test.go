package sql

import (
	"testing"

	"github.com/plu9in/whaledb/internal/domain/sql/lexer"
	"github.com/stretchr/testify/require"
)

func Test_Punctuators_And_Operators(t *testing.T) {
	d := lexer.PostgresDialect{}
	lx := lexer.NewLexerWithDialect("a[1:3]::int, --ceci est un commentaire\n b--ceci est un commentaire\n/*ceci est un commentaire*/||'WITH\\'';", d)

	// a
	tok := lx.Next()
	require.Equal(t, lexer.IDENT, tok.Kind)

	// [
	tok = lx.Next()
	require.Equal(t, lexer.TK_LBRACKET, tok.Kind)
	// 1
	tok = lx.Next()
	require.Equal(t, lexer.TK_NUMBER, tok.Kind)
	// :
	tok = lx.Next()
	require.Equal(t, lexer.TK_COLON, tok.Kind)
	// 3
	tok = lx.Next()
	require.Equal(t, lexer.TK_NUMBER, tok.Kind)
	// ]
	tok = lx.Next()
	require.Equal(t, lexer.TK_RBRACKET, tok.Kind)
	// ::
	tok = lx.Next()
	require.Equal(t, lexer.OP_CAST, tok.Kind)
	// int
	tok = lx.Next()
	require.Equal(t, lexer.IDENT, tok.Kind)
	// ,
	tok = lx.Next()
	require.Equal(t, lexer.TK_COMMA, tok.Kind)
	// b
	tok = lx.Next()
	require.Equal(t, lexer.IDENT, tok.Kind)
	// ||
	tok = lx.Next()
	require.Equal(t, lexer.OP_CONCAT, tok.Kind)
	// 'WITH'
	tok = lx.Next()
	require.Equal(t, lexer.TK_STRING, tok.Kind)
	require.Equal(t, "'WITH\\''", tok.Text)
	// ;
	tok = lx.Next()
	require.Equal(t, lexer.TK_SEMICOLON, tok.Kind)
}
