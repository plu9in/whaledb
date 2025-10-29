package sql

import (
	"testing"

	lex "github.com/plu9in/whaledb/internal/domain/sql/lexer"
)

func Test_String_SimpleQuotes_OK(t *testing.T) {
	d := lex.PostgresDialect{}
	lx := lex.NewLexerWithDialect("'It\\'s ok' FROM", d)

	tok := lx.Next()
	if tok.Kind != lex.TK_STRING || tok.Text != "'It\\'s ok'" {
		t.Fatalf("expected TK_STRING \"It's ok\", got kind=%v text=%q", tok.Kind, tok.Text)
	}
	tok = lx.Next()
	if tok.Kind != lex.KW_FROM {
		t.Fatalf("expected KW_FROM after string, got %v", tok.Kind)
	}
}

/*func Test_String_SimpleQuotes_Unterminated(t *testing.T) {
	d := PostgresDialect{}
	lx := NewLexerWithDialect("'unterminated", d)

	tok := lx.Next()
	if tok.Kind != TK_ERROR {
		t.Fatalf("expected TK_ERROR for unterminated string, got %v", tok.Kind)
	}
}
*/
