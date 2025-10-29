package lexer

type TokenKind int

const (
	EOF       TokenKind = iota // 0
	KW_AS                      // 1
	KW_FROM                    // 2
	KW_SELECT                  // 3
	KW_WITH                    // 4

	// Opérateurs
	OP_NOT_ILIKE        // 5
	OP_JSON_GET_TEXT    // 6
	OP_JSON_PATH_TEXT   // 7
	OP_INET_SUBNET_EQ_L // 8
	OP_INET_SUBNET_EQ_R // 9
	OP_DISTANCE         // 10
	OP_PERP             // 11
	OP_PARALLEL         // 12
	OP_LEFT_OF_SEG      // 13
	OP_RIGHT_OF_SEG     // 14
	OP_BELOW_SEG        // 15
	OP_ABOVE_SEG        // 16
	OP_CBRT             // 17
	OP_ILIKE            // 18
	OP_NOT_LIKE         // 19
	OP_NOT_REGEX_I      // 20
	OP_JSON_GET         // 21
	OP_JSON_PATH        // 22
	OP_CONTAINS         // 23
	OP_CONTAINED_BY     // 24
	OP_OVERLAP          // 25
	OP_CONCAT           // 26
	OP_CAST             // 27
	OP_SHIFT_L          // 28
	OP_SHIFT_R          // 29
	OP_LE               // 30
	OP_GE               // 31
	OP_NEQ              // 32
	OP_NEQ_ALT          // 33
	OP_REGEX_I          // 34
	OP_NOT_REGEX        // 35
	OP_LIKE             // 36
	OP_PLUS_EQ          // 37
	OP_MINUS_EQ         // 38
	OP_MUL_EQ           // 39
	OP_DIV_EQ           // 40
	OP_PLUS             // 41
	OP_MINUS            // 42
	OP_MUL              // 43
	OP_DIV              // 44
	OP_MOD              // 45
	OP_EQ               // 46
	OP_LT               // 47
	OP_GT               // 48
	OP_TILDE            // 49
	OP_BANG             // 50
	OP_AMP              // 51
	OP_BAR              // 52
	OP_CARET            // 53
	OP_QUESTION         // 54

	// Ponctuations
	TK_POINT     // 55
	TK_COMMA     // 56
	TK_SEMICOLON // 57
	TK_LPAREN    // 58
	TK_RPAREN    // 59
	TK_LBRACKET  // 60
	TK_RBRACKET  // 61
	TK_LBRACE    // 62
	TK_RBRACE    // 63
	TK_COLON     // 64

	// Autres tokens
	TK_STRING // 65
	TK_NUMBER // 66
	IDENT     // 67
)

type Token struct {
	Kind TokenKind
	Text string
}
