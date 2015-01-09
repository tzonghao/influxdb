package influxql

import (
	"strings"
)

// Token is a lexical token of the InfluxQL language.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	literal_beg
	// Literals
	IDENT        // main
	NUMBER       // 12345.67
	DURATION_VAL // 13h
	STRING       // "abc"
	BADSTRING    // "abc
	BADESCAPE    // \q
	TRUE         // true
	FALSE        // false
	literal_end

	operator_beg
	// Operators
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	AND // AND
	OR  // OR

	EQ  // =
	NEQ // !=
	LT  // <
	LTE // <=
	GT  // >
	GTE // >=
	operator_end

	LPAREN    // (
	RPAREN    // )
	COMMA     // ,
	SEMICOLON // ;
	DOT       // .

	keyword_beg
	// Keywords
	ALL
	ALTER
	AS
	ASC
	BEGIN
	BY
	CREATE
	CONTINUOUS
	DATABASE
	DATABASES
	DEFAULT
	DELETE
	DESC
	DROP
	DURATION
	END
	EXISTS
	EXPLAIN
	FIELD
	FROM
	GRANT
	GROUP
	IF
	INNER
	INSERT
	INTO
	KEYS
	LIMIT
	LIST
	MEASUREMENT
	MEASUREMENTS
	ON
	ORDER
	PASSWORD
	POLICY
	PRIVILEGES
	QUERIES
	QUERY
	READ
	REPLICATION
	RETENTION
	REVOKE
	SELECT
	SERIES
	TAG
	TO
	USER
	VALUES
	WHERE
	WITH
	WRITE
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	IDENT:        "IDENT",
	NUMBER:       "NUMBER",
	DURATION_VAL: "DURATION_VAL",
	STRING:       "STRING",
	TRUE:         "TRUE",
	FALSE:        "FALSE",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	AND: "AND",
	OR:  "OR",

	EQ:  "=",
	NEQ: "!=",
	LT:  "<",
	LTE: "<=",
	GT:  ">",
	GTE: ">=",

	LPAREN:    "(",
	RPAREN:    ")",
	COMMA:     ",",
	SEMICOLON: ";",
	DOT:       ".",

	ALL:          "ALL",
	ALTER:        "ALTER",
	AS:           "AS",
	ASC:          "ASC",
	BEGIN:        "BEGIN",
	BY:           "BY",
	CREATE:       "CREATE",
	CONTINUOUS:   "CONTINUOUS",
	DATABASE:     "DATABASE",
	DATABASES:    "DATABASES",
	DEFAULT:      "DEFAULT",
	DELETE:       "DELETE",
	DESC:         "DESC",
	DROP:         "DROP",
	DURATION:     "DURATION",
	END:          "END",
	EXISTS:       "EXISTS",
	EXPLAIN:      "EXPLAIN",
	FIELD:        "FIELD",
	FROM:         "FROM",
	GRANT:        "GRANT",
	GROUP:        "GROUP",
	IF:           "IF",
	INNER:        "INNER",
	INSERT:       "INSERT",
	INTO:         "INTO",
	KEYS:         "KEYS",
	LIMIT:        "LIMIT",
	LIST:         "LIST",
	MEASUREMENT:  "MEASUREMENT",
	MEASUREMENTS: "MEASUREMENTS",
	ON:           "ON",
	ORDER:        "ORDER",
	PASSWORD:     "PASSWORD",
	POLICY:       "POLICY",
	PRIVILEGES:   "PRIVILEGES",
	QUERIES:      "QUERIES",
	QUERY:        "QUERY",
	READ:         "READ",
	REPLICATION:  "REPLICATION",
	RETENTION:    "RETENTION",
	REVOKE:       "REVOKE",
	SELECT:       "SELECT",
	SERIES:       "SERIES",
	TAG:          "TAG",
	TO:           "TO",
	USER:         "USER",
	VALUES:       "VALUES",
	WHERE:        "WHERE",
	WITH:         "WITH",
	WRITE:        "WRITE",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for tok := keyword_beg + 1; tok < keyword_end; tok++ {
		keywords[strings.ToUpper(tokens[tok])] = tok
		keywords[strings.ToLower(tokens[tok])] = tok
	}
	for _, tok := range []Token{AND, OR} {
		keywords[strings.ToUpper(tokens[tok])] = tok
		keywords[strings.ToLower(tokens[tok])] = tok
	}
	keywords["true"] = TRUE
	keywords["false"] = FALSE
}

// String returns the string representation of the token.
func (tok Token) String() string {
	if tok >= 0 && tok < Token(len(tokens)) {
		return tokens[tok]
	}
	return ""
}

// Precedence returns the operator precedence of the binary operator token.
func (tok Token) Precedence() int {
	switch tok {
	case OR:
		return 1
	case AND:
		return 2
	case EQ, NEQ, LT, LTE, GT, GTE:
		return 3
	case ADD, SUB:
		return 4
	case MUL, DIV:
		return 5
	}
	return 0
}

// isOperator returns true for operator tokens.
func (tok Token) isOperator() bool { return tok > operator_beg && tok < operator_end }

// tokstr returns a literal if provided, otherwise returns the token string.
func tokstr(tok Token, lit string) string {
	if lit != "" {
		return lit
	}
	return tok.String()
}

// Lookup returns the token associated with a given string.
func Lookup(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// Pos specifies the line and character position of a token.
// The Char and Line are both zero-based indexes.
type Pos struct {
	Line int
	Char int
}
