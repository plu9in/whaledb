package sql_test

import (
	"testing"

	"github.com/plu9in/whaledb/internal/domain/sql/lexer"
)

func TestLexer_Recognizes_SELECT(t *testing.T) {
	lx := lexer.NewLexer("SELECT")
	tok := lx.Next()
	if tok.Kind != lexer.KW_SELECT {
		t.Fatalf("expected KW_SELECT, got kind=%v text=%q", tok.Kind, tok.Text)
	}
	tok = lx.Next()
	if tok.Kind != lexer.EOF {
		t.Fatalf("expected EOF, got kind=%v text=%q", tok.Kind, tok.Text)
	}
}

func TestLexer_DoesNotMatch_PREFIX(t *testing.T) {
	lx := lexer.NewLexer("SELECTX")
	tok := lx.Next()
	if tok.Kind == lexer.KW_SELECT {
		t.Fatalf("should not be KW_SELECT on SELECTX")
	}
	if tok.Kind != lexer.IDENT || tok.Text != "SELECTX" {
		t.Fatalf("expected IDENT 'SELECTX', got kind=%v text=%q", tok.Kind, tok.Text)
	}
}

func TestLexer_Skips_Leading_Whitespace_Before_SELECT(t *testing.T) {
	lx := lexer.NewLexer("   \t\nSELECT")
	tok := lx.Next()
	if tok.Kind != lexer.KW_SELECT {
		t.Fatalf("expected KW_SELECT with leading spaces, got %v (%q)", tok.Kind, tok.Text)
	}
}

func TestLexer_Recognizes_SELECT_CaseInsensitive(t *testing.T) {
	lx := lexer.NewLexer("select")
	tok := lx.Next()
	if tok.Kind != lexer.KW_SELECT {
		t.Fatalf("expected KW_SELECT for lowercase, got %v (%q)", tok.Kind, tok.Text)
	}
}

func TestLexer_Recognizes_FROM(t *testing.T) {
	lx := lexer.NewLexer("FROM")
	tok := lx.Next()
	if tok.Kind != lexer.KW_FROM {
		t.Fatalf("expected KW_FROM, got %v (%q)", tok.Kind, tok.Text)
	}
}

func TestLexer_TwoTokens_SELECT_FROM(t *testing.T) {
	lx := lexer.NewLexer("SELECT FROM")
	tok1 := lx.Next()
	if tok1.Kind != lexer.KW_SELECT {
		t.Fatalf("want KW_SELECT, got %v", tok1.Kind)
	}
	tok2 := lx.Next()
	if tok2.Kind != lexer.KW_FROM {
		t.Fatalf("want KW_FROM, got %v", tok2.Kind)
	}
	tok3 := lx.Next()
	if tok3.Kind != lexer.EOF {
		t.Fatalf("want EOF, got %v", tok3.Kind)
	}
}

func TestLexer_SELECT_Order(t *testing.T) {
	lx := lexer.NewLexer("SELECT * FROM schema.table as b;")

	lx.Next() // Skip SELECT

	tok := lx.Next() // Read *
	if tok.Kind != lexer.OP_MUL {
		t.Fatalf("want OP_MUL, got %v", tok.Kind)
	}

	lx.Next() // Skip FROM

	tok = lx.Next() // Read schema
	if tok.Kind != lexer.IDENT || tok.Text != "schema" {
		t.Fatalf("want IDENT 'schema', got %v (%q)", tok.Kind, tok.Text)
	}

	tok = lx.Next() // Read .
	if tok.Kind != lexer.TK_POINT || tok.Text != "." {
		t.Fatalf("want TK_POINT '.', got %v (%q)", tok.Kind, tok.Text)
	}

	tok = lx.Next() // Read table
	if tok.Kind != lexer.IDENT || tok.Text != "table" {
		t.Fatalf("want IDENT 'table', got %v (%q)", tok.Kind, tok.Text)
	}

	tok = lx.Next() // Read as
	if tok.Kind != lexer.KW_AS || tok.Text != "AS" {
		t.Fatalf("want KW_AS 'as', got %v (%q)", tok.Kind, tok.Text)
	}
	tok = lx.Next() // Read as
	if tok.Kind != lexer.IDENT || tok.Text != "b" {
		t.Fatalf("want KW_AS 'as', got %v (%q)", tok.Kind, tok.Text)
	}
	tok = lx.Next()
	if tok.Kind != lexer.TK_SEMICOLON || tok.Text != ";" {
		t.Fatalf("want KW_SEMICOLON ';', got %v (%q)", tok.Kind, tok.Text)
	}
}

func TestLexer_SELECT2_Order(t *testing.T) {
	lx := lexer.NewLexer("SELECT a.*, b as quote FROM schema.table;")
	lx.Next()        // Skip SELECT
	lx.Next()        // Read a
	lx.Next()        // Read .
	tok := lx.Next() // Read *
	if tok.Kind != lexer.OP_MUL {
		t.Fatalf("want KW_WILDCARD, got %v", tok.Kind)
	}
	tok = lx.Next() // Read ,
	if tok.Kind != lexer.TK_COMMA {
		t.Fatalf("want TK_COMMA, got %v", tok.Kind)
	}
	tok = lx.Next() // Read b
	if tok.Kind != lexer.IDENT || tok.Text != "b" {
		t.Fatalf("want IDENT 'b', got %v (%q)", tok.Kind, tok.Text)
	}
	tok = lx.Next() // Read as
	if tok.Kind != lexer.KW_AS || tok.Text != "AS" {
		t.Fatalf("want KW_AS 'AS', got %v (%q)", tok.Kind, tok.Text)
	}
	tok = lx.Next() // Read quote
	if tok.Kind != lexer.IDENT || tok.Text != "quote" {
		t.Fatalf("want IDENT 'quote', got %v (%q)", tok.Kind, tok.Text)
	}

}
