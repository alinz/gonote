package lexer

import "fmt"

type TokenType int

const (
	//DocStart when a new document is loaded or parsed
	tokenDocStart TokenType = iota
	tokenArrayStart
	tokenMap
	tokenNumber
	tokenString
	tokenNull
	tokenTrue
	tokenFalse
	tokenSpace
	tokenEnter
	tokenConfig
	tokenError
	tokenEnd
)

type Token struct {
	typ TokenType
	val string
}

func (t Token) String() string {
	tokType := ""
	switch t.typ {
	case tokenEnter:
		tokType = "tokenEnter"
	case tokenConfig:
		tokType = "tokenConfig"
	case tokenError:
		tokType = "tokenError"
	case tokenArrayStart:
		tokType = "tokenArrayStart"
	case tokenMap:
		tokType = "tokenMap"
	case tokenString:
		tokType = "tokenString"
	case tokenNull:
		tokType = "tokenNull"
	case tokenTrue:
		tokType = "tokenTrue"
	case tokenFalse:
		tokType = "tokenFalse"
	case tokenSpace:
		tokType = "tokenSpace"
	case tokenNumber:
		tokType = "tokenNumber"
	case tokenEnd:
		tokType = "tokenEnd"
	}

	if len(t.val) > 20 {
		return fmt.Sprintf("%s: %.20q...", tokType, t.val)
	}
	return fmt.Sprintf("%s: %q", tokType, t.val)
}
