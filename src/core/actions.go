package core

type ActionType struct {
	RETURN    AType
	NO_RETURN AType
}

var ATYPE = ActionType{RETURN: 1, NO_RETURN: 2}
