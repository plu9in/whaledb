package lexer

// Dialect décrit les éléments variables d’un SQL dialecte (mots-clés, opérateurs).
// Objectif : ouvrir à l’extension (OCP) et inverser la dépendance (DIP).
type Dialect interface {
	// Keywords retourne une map des mots-clés (MAJUSCULES) vers un TokenKind existant.
	Keywords() map[string]TokenKind
	// Operators retourne la liste des opérateurs/ponctuations à reconnaître en priorité.
	// L’ordre compte si tu ajoutes plus tard des opérateurs multi-caractères (ex: "<=", "<>", "!=").
	Operators() []string // ordre: du plus long au plus court (pour <=, <>...)
	OperatorKinds() map[string]TokenKind
	Punctuators() []string // ponctuateurs unitaires non-opérateurs
	PunctuatorKinds() map[string]TokenKind
}
