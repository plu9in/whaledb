package lexer

// PostgresDialect : implémentation minimale par défaut,
// limitée à ce que ton lexer gérait déjà (SELECT, FROM, AS et .,;()*).
type PostgresDialect struct{}

func (PostgresDialect) Keywords() map[string]TokenKind {
	return map[string]TokenKind{
		"AS":     KW_AS,
		"FROM":   KW_FROM,
		"SELECT": KW_SELECT,
		"WITH":   KW_WITH,
	}
}

func (PostgresDialect) Operators() []string {
	// plus tard tu pourras mettre: "<=", "<>", "!=", "||", ...
	return []string{
		// 4 caractères
		"!~~*", // NOT ILIKE (forme opérateur interne)
		// 3 caractères
		"->>", "#>>", // JSON/JSONB (extraction texte)
		"<<=", ">>=", // inet/cidr inclusion réseau
		"<->",        // distance (géométrie)
		"?-|", "?||", // perpendiculaire / parallèles (géométrie)
		"<<|", "|>>", "&<|", "|&>", // relations géométriques
		"||/",               // racine cubique
		"~~*", "!~~", "!~*", // LIKE/ILIKE négations (formes opérateurs)
		// 2 caractères
		"->", "#>", "@>", "<@", // JSON/JSONB et tableaux (contenance)
		"&&",       // overlap (plages, tableaux, géométrie)
		"||",       // concat chaîne / JSON
		"::",       // cast
		"<<", ">>", // décalage bitwise ou géométrie
		"<=", ">=", "<>", "!=", // comparaisons
		"~*", "!~", // regex (case-insensitive, négatif)
		"~~",                   // LIKE (forme opérateur)
		"+=", "-=", "*=", "/=", // (au cas où pour extensions / custom)
		// 1 caractère (unitaires et ponctuation)
		"+", "-", "*", "/", "%", // arithmétiques
		"=", "<", ">", // comparaisons
		"~", "!", "&", "|", "^", "?", // bitwise / regex / divers
	}
}

func (PostgresDialect) OperatorKinds() map[string]TokenKind {
	return map[string]TokenKind{
		// 4
		"!~~*": OP_NOT_ILIKE,
		// 3
		"->>": OP_JSON_GET_TEXT, "#>>": OP_JSON_PATH_TEXT,
		"<<=": OP_INET_SUBNET_EQ_L, ">>=": OP_INET_SUBNET_EQ_R,
		"<->": OP_DISTANCE,
		"?-|": OP_PERP, "?||": OP_PARALLEL,
		"<<|": OP_LEFT_OF_SEG, "|>>": OP_RIGHT_OF_SEG, "&<|": OP_BELOW_SEG, "|&>": OP_ABOVE_SEG,
		"||/": OP_CBRT,
		"~~*": OP_ILIKE, "!~~": OP_NOT_LIKE, "!~*": OP_NOT_REGEX_I,
		// 2
		"->": OP_JSON_GET, "#>": OP_JSON_PATH, "@>": OP_CONTAINS, "<@": OP_CONTAINED_BY,
		"&&": OP_OVERLAP,
		"||": OP_CONCAT,
		"::": OP_CAST,
		"<<": OP_SHIFT_L, ">>": OP_SHIFT_R,
		"<=": OP_LE, ">=": OP_GE, "<>": OP_NEQ, "!=": OP_NEQ_ALT,
		"~*": OP_REGEX_I, "!~": OP_NOT_REGEX,
		"~~": OP_LIKE,
		"+=": OP_PLUS_EQ, "-=": OP_MINUS_EQ, "*=": OP_MUL_EQ, "/=": OP_DIV_EQ,
		// 1
		"+": OP_PLUS, "-": OP_MINUS, "*": OP_MUL /* ou OP_MUL */, "/": OP_DIV, "%": OP_MOD,
		"=": OP_EQ, "<": OP_LT, ">": OP_GT,
		"~": OP_TILDE, "!": OP_BANG, "&": OP_AMP, "|": OP_BAR, "^": OP_CARET, "?": OP_QUESTION,
	}
}

func (PostgresDialect) Punctuators() []string {
	return []string{".", ",", ";", "(", ")", "[", "]", "{", "}", ":"}
}

func (PostgresDialect) PunctuatorKinds() map[string]TokenKind {
	return map[string]TokenKind{
		".": TK_POINT,
		",": TK_COMMA,
		";": TK_SEMICOLON,
		"(": TK_LPAREN,
		")": TK_RPAREN,
		"[": TK_LBRACKET,
		"]": TK_RBRACKET,
		"{": TK_LBRACE,
		"}": TK_RBRACE,
		":": TK_COLON, // attention: l’opérateur '::' reste côté Operators()
	}
}
