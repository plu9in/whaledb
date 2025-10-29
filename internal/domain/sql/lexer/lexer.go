package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input   string
	dialect Dialect
}

// Conserve l’API existante : par défaut on utilise un dialecte Postgres minimal.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, dialect: PostgresDialect{}}
}

// Permet d’injecter un autre dialecte (MySQL, Oracle, …) sans toucher au Lexer.
func NewLexerWithDialect(input string, d Dialect) *Lexer {
	return &Lexer{input: input, dialect: d}
}

func (lx *Lexer) TrimLeft(cutset string) {
	lx.input = strings.TrimLeft(lx.input, cutset)
}

func (lx *Lexer) readOneLineComment() Token {
	// Handle SQL comments starting with --
	if strings.HasPrefix(lx.input, "--") {
		i := strings.Index(lx.input, "\n")
		println("one line comment: ", i)
		if i == -1 {
			// Comment goes until the end of input
			lx.input = ""
			return Token{Kind: EOF, Text: ""}
		}
		// Skip the comment and continue
		lx.input = lx.input[i+1:]
	}
	return Token{}
}

func (lx *Lexer) readMultiLineComment() Token {
	// Handle SQL comments starting with /* and ending with */
	if strings.HasPrefix(lx.input, "/*") {
		endIdx := strings.Index(lx.input, "*/")
		println("multilines comment: ", endIdx)
		if endIdx == -1 {
			// Unterminated comment, consume the rest of the input
			lx.input = ""
			return Token{Kind: EOF, Text: ""}
		}
		// Skip the comment and continue
		lx.input = lx.input[endIdx+2:]
	}
	return Token{}
}

func (lx *Lexer) Next() Token {
	lx.TrimLeft(" \t\r\n")

	for strings.HasPrefix(lx.input, "--") || strings.HasPrefix(lx.input, "/*") {
		if lx.input == "" {
			return Token{Kind: EOF, Text: ""}
		}

		if strings.HasPrefix(lx.input, "--") {
			if commentTok := lx.readOneLineComment(); commentTok.Kind != 0 {
				return commentTok
			}
		} else if strings.HasPrefix(lx.input, "/*") {
			if commentTok := lx.readMultiLineComment(); commentTok.Kind != 0 {
				return commentTok
			}
		}

		lx.TrimLeft(" \t\r\n")
	}

	if tok, ok := lx.readNumberToken(); ok {
		return tok
	}

	if tok, ok := lx.readStringToken(); ok {
		return tok
	}

	if op := lx.readOperator(); op != "" {
		// ⬇️ lookup comme pour les mots-clés
		if kind, ok := lx.dialect.OperatorKinds()[op]; ok {
			if strings.HasPrefix(op, "--") || strings.HasPrefix(op, "/*") {
				// Si c’est un commentaire, on l’ignore et on continue
				return lx.Next()
			}
			return Token{Kind: kind, Text: op}
		}
		// fallback si non mappé (optionnel : crée un TokenKind OP/UNKNOWN)
		return Token{Kind: IDENT, Text: op}
	}

	if op := lx.readPunctuator(); op != "" {
		// ⬇️ lookup comme pour les mots-clés
		if kind, ok := lx.dialect.PunctuatorKinds()[op]; ok {
			return Token{Kind: kind, Text: op}
		}
		// fallback si non mappé (optionnel : crée un TokenKind OP/UNKNOWN)
		return Token{Kind: IDENT, Text: op}
	}

	if lx.input == "" {
		return Token{Kind: EOF, Text: ""}
	}

	word := lx.readWord(lx.input)
	// Mapping mots-clés via le dialecte (OCP-ready)
	if kwKind, ok := lx.dialect.Keywords()[strings.ToUpper(word)]; ok {
		return Token{Kind: kwKind, Text: strings.ToUpper(word)}
	}

	return Token{Kind: IDENT, Text: word}
}

func (lx *Lexer) readPunctuator() string {
	if len(lx.input) == 0 {
		return ""
	}
	// Délègue la liste des ponctuations au dialecte
	for _, p := range lx.dialect.Punctuators() {
		if strings.HasPrefix(lx.input, p) {
			lx.input = lx.input[len(p):]
			return p
		}
	}
	return ""
}

func (lx *Lexer) readOperator() string {
	if len(lx.input) == 0 {
		return ""
	}
	// Délègue la liste des opérateurs
	for _, p := range lx.dialect.Operators() {
		if strings.HasPrefix(lx.input, p) {
			lx.input = lx.input[len(p):]
			return p
		}
	}
	return ""
}

func (lx *Lexer) readNumberToken() (Token, bool) {
	i := 0
	hasDigits := false
	for i < len(lx.input) {
		r, size := utf8.DecodeRuneInString(lx.input[i:])
		if unicode.IsDigit(r) {
			hasDigits = true
			i += size
		} else {
			break
		}
	}
	if hasDigits {
		numberText := lx.input[:i]
		lx.input = lx.input[i:]
		return Token{Kind: TK_NUMBER, Text: numberText}, true
	}
	return Token{}, false
}

func (lx *Lexer) readStringToken() (Token, bool) {
	if len(lx.input) == 0 {
		return Token{}, false
	}
	quoteChar := lx.input[0]
	if quoteChar != '\'' && quoteChar != '"' {
		return Token{}, false
	}
	i := 1
	for i < len(lx.input) {
		r, size := utf8.DecodeRuneInString(lx.input[i:])
		if byte(r) == quoteChar {
			// Check if the quoteChar is escaped
			if i > 0 && lx.input[i-1] == '\\' {
				i += size
				continue
			}
			// fin de chaîne
			strText := lx.input[:i+size]
			lx.input = lx.input[i+size:]
			return Token{Kind: TK_STRING, Text: strText}, true
		}
		i += size
	}
	// chaîne non terminée
	return Token{}, false
}

// Un mot se termine sur espace OU dès que la sous-chaîne courante commence par un opérateur du dialecte
func (lx *Lexer) readWord(input string) string {
	i := 0
	for i < len(input) {
		r, size := utf8.DecodeRuneInString(input[i:])
		if unicode.IsSpace(r) {
			lx.input = input[i:]
			return input[:i]
		}
		// si on croise un opérateur (., ;, (), *, ||, <=, <>, !=, etc.) OU un début de ponctuateur, on s’arrête avant
		if lx.startsWithAnyOperator(input[i:]) || lx.startsWithAnyPunctuator(input[i:]) {
			lx.input = input[i:]
			return input[:i]
		}
		i += size
	}
	// fin de chaîne
	lx.input = ""
	return input
}

// Renvoie true si s commence par un opérateur déclaré par le dialecte
func (lx *Lexer) startsWithAnyOperator(s string) bool {
	for _, op := range lx.dialect.Operators() {
		if strings.HasPrefix(s, op) {
			return true
		}
	}
	return false
}

func (lx *Lexer) startsWithAnyPunctuator(s string) bool {
	for _, p := range lx.dialect.Punctuators() {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
