// Code generated by goyacc - DO NOT EDIT.

package sql

import __yyfmt__ "fmt"

import (
	"strconv"
)

type yySymType struct {
	yys            int
	str            string
	strList        []string
	boolean        bool
	fieldType      FieldType
	compareOperate CompareOperate

	stmt     Statement
	stmtList []Statement

	createStmt        *CreateStmt
	createTable       *CreateTable
	createField       *CreateField
	createIndex       *CreateIndex
	createTableOption *CreateTableOption

	insertStmt *InsertStmt

	updateStmt  *UpdateStmt
	updateValue map[string]string

	deleteStmt *DeleteStmt

	selectStmt      *SelectStmt
	selectFieldList []*SelectField
	selectWhereList []SelectWhere
	selectOrderList []*SelectOrder
	selectLimit     *SelectLimit
}

type yyXError struct {
	state, xsym int
}

const (
	yyDefault = 57375
	yyEofCode = 57344
	AND       = 57364
	ASC       = 57367
	BY        = 57366
	COMP_GE   = 57373
	COMP_LE   = 57372
	COMP_NE   = 57371
	CREATE    = 57346
	DEFAULT   = 57352
	DELETE    = 57359
	DESC      = 57368
	FROM      = 57361
	INDEX     = 57351
	INSERT    = 57354
	INTO      = 57355
	KEY       = 57348
	LIMIT     = 57369
	NOT       = 57349
	NULL      = 57350
	OFFSET    = 57370
	OR        = 57363
	ORDER     = 57365
	PRIMARY   = 57353
	SELECT    = 57360
	SET       = 57358
	TABLE     = 57347
	UPDATE    = 57357
	VALUE     = 57356
	VARIABLE  = 57374
	WHERE     = 57362
	yyErrCode = 57345

	yyMaxDepth = 200
	yyTabOfs   = -70
)

var (
	yyPrec = map[int]int{
		OR:  0,
		AND: 1,
		'+': 2,
		'-': 2,
		'*': 3,
		'/': 3,
	}

	yyXLAT = map[int]int{
		57374: 0,  // VARIABLE (41x)
		44:    1,  // ',' (37x)
		41:    2,  // ')' (36x)
		59:    3,  // ';' (35x)
		57387: 4,  // Expr (31x)
		57369: 5,  // LIMIT (20x)
		57344: 6,  // $end (15x)
		57346: 7,  // CREATE (15x)
		57359: 8,  // DELETE (15x)
		57354: 9,  // INSERT (15x)
		57360: 10, // SELECT (15x)
		57357: 11, // UPDATE (15x)
		57364: 12, // AND (9x)
		57363: 13, // OR (9x)
		57365: 14, // ORDER (9x)
		40:    15, // '(' (8x)
		61:    16, // '=' (6x)
		57352: 17, // DEFAULT (6x)
		57362: 18, // WHERE (6x)
		57361: 19, // FROM (5x)
		57350: 20, // NULL (5x)
		60:    21, // '<' (4x)
		62:    22, // '>' (4x)
		57373: 23, // COMP_GE (4x)
		57372: 24, // COMP_LE (4x)
		57371: 25, // COMP_NE (4x)
		57367: 26, // ASC (3x)
		57378: 27, // CompareOperate (3x)
		57368: 28, // DESC (3x)
		57349: 29, // NOT (3x)
		57399: 30, // SelectWhere (3x)
		57400: 31, // SelectWhereList (3x)
		57377: 32, // Ascend (2x)
		57379: 33, // CreateField (2x)
		57380: 34, // CreateIndex (2x)
		57381: 35, // CreatePrimary (2x)
		57382: 36, // CreateStmt (2x)
		57386: 37, // DeleteStmt (2x)
		57351: 38, // INDEX (2x)
		57391: 39, // InsertStmt (2x)
		57353: 40, // PRIMARY (2x)
		57395: 41, // SelectLimit (2x)
		57398: 42, // SelectStmt (2x)
		57358: 43, // SET (2x)
		57401: 44, // Stmt (2x)
		57403: 45, // UpdateStmt (2x)
		57356: 46, // VALUE (2x)
		57405: 47, // VaribleList (2x)
		57376: 48, // AllowNull (1x)
		57366: 49, // BY (1x)
		57383: 50, // CreateTable (1x)
		57384: 51, // CreateTableOption (1x)
		57385: 52, // DefaultVal (1x)
		57388: 53, // FieldType (1x)
		57389: 54, // InsertField (1x)
		57390: 55, // InsertFieldList (1x)
		57392: 56, // InsertValue (1x)
		57393: 57, // InsertValueList (1x)
		57355: 58, // INTO (1x)
		57348: 59, // KEY (1x)
		57370: 60, // OFFSET (1x)
		57394: 61, // SelectFieldList (1x)
		57396: 62, // SelectOrder (1x)
		57397: 63, // SelectOrderList (1x)
		57406: 64, // start (1x)
		57402: 65, // StmtList (1x)
		57347: 66, // TABLE (1x)
		57404: 67, // UpdateValue (1x)
		57375: 68, // $default (0x)
		42:    69, // '*' (0x)
		43:    70, // '+' (0x)
		45:    71, // '-' (0x)
		47:    72, // '/' (0x)
		57345: 73, // error (0x)
	}

	yySymNames = []string{
		"VARIABLE",
		"','",
		"')'",
		"';'",
		"Expr",
		"LIMIT",
		"$end",
		"CREATE",
		"DELETE",
		"INSERT",
		"SELECT",
		"UPDATE",
		"AND",
		"OR",
		"ORDER",
		"'('",
		"'='",
		"DEFAULT",
		"WHERE",
		"FROM",
		"NULL",
		"'<'",
		"'>'",
		"COMP_GE",
		"COMP_LE",
		"COMP_NE",
		"ASC",
		"CompareOperate",
		"DESC",
		"NOT",
		"SelectWhere",
		"SelectWhereList",
		"Ascend",
		"CreateField",
		"CreateIndex",
		"CreatePrimary",
		"CreateStmt",
		"DeleteStmt",
		"INDEX",
		"InsertStmt",
		"PRIMARY",
		"SelectLimit",
		"SelectStmt",
		"SET",
		"Stmt",
		"UpdateStmt",
		"VALUE",
		"VaribleList",
		"AllowNull",
		"BY",
		"CreateTable",
		"CreateTableOption",
		"DefaultVal",
		"FieldType",
		"InsertField",
		"InsertFieldList",
		"InsertValue",
		"InsertValueList",
		"INTO",
		"KEY",
		"OFFSET",
		"SelectFieldList",
		"SelectOrder",
		"SelectOrderList",
		"start",
		"StmtList",
		"TABLE",
		"UpdateValue",
		"$default",
		"'*'",
		"'+'",
		"'-'",
		"'/'",
		"error",
	}

	yyTokenLiteralStrings = map[int]string{
		57369: "LIMIT",
		57346: "CREATE",
		57359: "DELETE",
		57354: "INSERT",
		57360: "SELECT",
		57357: "UPDATE",
		57364: "AND",
		57363: "OR",
		57365: "ORDER",
		57352: "DEFAULT",
		57362: "WHERE",
		57361: "FROM",
		57350: "NULL",
		57373: ">=",
		57372: "<=",
		57371: "!=",
		57367: "ASC",
		57368: "DESC",
		57349: "NOT",
		57351: "INDEX",
		57353: "PRIMARY",
		57358: "SET",
		57356: "VALUE",
		57366: "BY",
		57355: "INTO",
		57348: "KEY",
		57370: "OFFSET",
		57347: "TABLE",
	}

	yyReductions = map[int]struct{ xsym, components int }{
		0:  {0, 1},
		1:  {64, 1},
		2:  {4, 1},
		3:  {47, 1},
		4:  {47, 3},
		5:  {44, 1},
		6:  {44, 1},
		7:  {44, 1},
		8:  {44, 1},
		9:  {44, 1},
		10: {65, 1},
		11: {65, 2},
		12: {48, 0},
		13: {48, 1},
		14: {48, 2},
		15: {52, 0},
		16: {52, 1},
		17: {52, 2},
		18: {52, 2},
		19: {53, 1},
		20: {36, 8},
		21: {50, 1},
		22: {50, 1},
		23: {50, 1},
		24: {50, 3},
		25: {50, 3},
		26: {50, 3},
		27: {33, 4},
		28: {34, 5},
		29: {35, 5},
		30: {51, 0},
		31: {39, 6},
		32: {54, 3},
		33: {55, 0},
		34: {55, 1},
		35: {56, 4},
		36: {57, 0},
		37: {57, 1},
		38: {45, 6},
		39: {67, 3},
		40: {67, 5},
		41: {37, 5},
		42: {32, 0},
		43: {32, 1},
		44: {32, 1},
		45: {27, 1},
		46: {27, 1},
		47: {27, 1},
		48: {27, 1},
		49: {27, 1},
		50: {27, 1},
		51: {42, 4},
		52: {42, 8},
		53: {61, 1},
		54: {61, 3},
		55: {30, 0},
		56: {30, 2},
		57: {31, 3},
		58: {31, 5},
		59: {31, 5},
		60: {31, 5},
		61: {31, 5},
		62: {62, 0},
		63: {62, 3},
		64: {63, 2},
		65: {63, 4},
		66: {41, 0},
		67: {41, 2},
		68: {41, 4},
		69: {41, 4},
	}

	yyXErrors = map[yyXError]string{}

	yyParseTab = [137][]uint16{
		// 0
		{7: 79, 82, 80, 83, 81, 36: 73, 77, 39: 75, 42: 74, 44: 78, 76, 64: 71, 72},
		{6: 70},
		{6: 69, 79, 82, 80, 83, 81, 36: 73, 77, 39: 75, 42: 74, 44: 206, 76},
		{6: 65, 65, 65, 65, 65, 65},
		{6: 64, 64, 64, 64, 64, 64},
		// 5
		{6: 63, 63, 63, 63, 63, 63},
		{6: 62, 62, 62, 62, 62, 62},
		{6: 61, 61, 61, 61, 61, 61},
		{6: 60, 60, 60, 60, 60, 60},
		{66: 171},
		// 10
		{58: 154},
		{84, 4: 142},
		{19: 138},
		{84, 4: 86, 61: 85},
		{68, 68, 68, 68, 5: 68, 12: 68, 68, 68, 68, 68, 68, 68, 68, 68, 68, 68, 68, 68, 68, 68, 28: 68, 68, 43: 68},
		// 15
		{1: 89, 3: 4, 5: 90, 19: 88, 41: 87},
		{1: 17, 3: 17, 5: 17, 19: 17},
		{3: 137},
		{84, 4: 97},
		{84, 4: 96},
		// 20
		{91},
		{1: 92, 3: 3, 60: 93},
		{95},
		{94},
		{3: 1},
		// 25
		{3: 2},
		{1: 16, 3: 16, 5: 16, 19: 16},
		{3: 15, 5: 15, 14: 15, 18: 99, 30: 98},
		{3: 8, 5: 8, 14: 125, 62: 124},
		{84, 4: 101, 31: 100},
		// 30
		{3: 14, 5: 14, 12: 111, 110, 14},
		{16: 102, 21: 103, 104, 106, 105, 107, 27: 108},
		{25},
		{24},
		{23},
		// 35
		{22},
		{21},
		{20},
		{84, 4: 109},
		{2: 13, 13, 5: 13, 12: 13, 13, 13},
		// 40
		{84, 4: 118, 15: 119},
		{84, 4: 112, 15: 113},
		{16: 102, 21: 103, 104, 106, 105, 107, 27: 116},
		{84, 4: 101, 31: 114},
		{2: 115, 12: 111, 110},
		// 45
		{2: 9, 9, 5: 9, 12: 9, 9, 9},
		{84, 4: 117},
		{2: 11, 11, 5: 11, 12: 11, 11, 11},
		{16: 102, 21: 103, 104, 106, 105, 107, 27: 122},
		{84, 4: 101, 31: 120},
		// 50
		{2: 121, 12: 111, 110},
		{2: 10, 10, 5: 10, 12: 10, 10, 10},
		{84, 4: 123},
		{2: 12, 12, 5: 12, 12: 12, 12, 12},
		{3: 4, 5: 90, 41: 135},
		// 55
		{49: 126},
		{84, 4: 128, 63: 127},
		{1: 132, 3: 7, 5: 7},
		{1: 28, 3: 28, 5: 28, 26: 129, 28: 130, 32: 131},
		{1: 27, 3: 27, 5: 27},
		// 60
		{1: 26, 3: 26, 5: 26},
		{1: 6, 3: 6, 5: 6},
		{84, 4: 133},
		{1: 28, 3: 28, 5: 28, 26: 129, 28: 130, 32: 134},
		{1: 5, 3: 5, 5: 5},
		// 65
		{3: 136},
		{6: 18, 18, 18, 18, 18, 18},
		{6: 19, 19, 19, 19, 19, 19},
		{84, 4: 139},
		{3: 15, 18: 99, 30: 140},
		// 70
		{3: 141},
		{6: 29, 29, 29, 29, 29, 29},
		{43: 143},
		{84, 4: 145, 67: 144},
		{1: 149, 3: 15, 18: 99, 30: 148},
		// 75
		{16: 146},
		{84, 4: 147},
		{1: 31, 3: 31, 18: 31},
		{3: 153},
		{84, 4: 150},
		// 80
		{16: 151},
		{84, 4: 152},
		{1: 30, 3: 30, 18: 30},
		{6: 32, 32, 32, 32, 32, 32},
		{84, 4: 155},
		// 85
		{15: 157, 54: 156},
		{46: 165, 56: 164},
		{84, 2: 37, 4: 158, 47: 159, 55: 160},
		{1: 67, 67},
		{1: 162, 36},
		// 90
		{2: 161},
		{46: 38},
		{84, 4: 163},
		{1: 66, 66},
		{3: 170},
		// 95
		{15: 166},
		{84, 2: 34, 4: 158, 47: 167, 57: 168},
		{1: 162, 33},
		{2: 169},
		{3: 35},
		// 100
		{6: 39, 39, 39, 39, 39, 39},
		{84, 4: 172},
		{15: 173},
		{84, 4: 178, 33: 175, 176, 177, 38: 179, 40: 180, 50: 174},
		{1: 200, 199},
		// 105
		{1: 49, 49},
		{1: 48, 48},
		{1: 47, 47},
		{84, 4: 189, 53: 190},
		{84, 4: 185},
		// 110
		{59: 181},
		{15: 182},
		{84, 4: 183},
		{2: 184},
		{1: 41, 41},
		// 115
		{15: 186},
		{84, 4: 187},
		{2: 188},
		{1: 42, 42},
		{1: 51, 51, 17: 51, 20: 51, 29: 51},
		// 120
		{1: 58, 58, 17: 58, 20: 191, 29: 192, 48: 193},
		{1: 57, 57, 17: 57},
		{20: 198},
		{1: 55, 55, 17: 194, 52: 195},
		{84, 54, 54, 4: 197, 20: 196},
		// 125
		{1: 43, 43},
		{1: 53, 53},
		{1: 52, 52},
		{1: 56, 56, 17: 56},
		{3: 40, 51: 204},
		// 130
		{84, 4: 178, 33: 201, 202, 203, 38: 179, 40: 180},
		{1: 46, 46},
		{1: 45, 45},
		{1: 44, 44},
		{3: 205},
		// 135
		{6: 50, 50, 50, 50, 50, 50},
		{6: 59, 59, 59, 59, 59, 59},
	}
)

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyLexerEx interface {
	yyLexer
	Reduced(rule, state int, lval *yySymType) bool
}

func yySymName(c int) (s string) {
	x, ok := yyXLAT[c]
	if ok {
		return yySymNames[x]
	}

	if c < 0x7f {
		return __yyfmt__.Sprintf("%q", c)
	}

	return __yyfmt__.Sprintf("%d", c)
}

func yylex1(yylex yyLexer, lval *yySymType) (n int) {
	n = yylex.Lex(lval)
	if n <= 0 {
		n = yyEofCode
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("\nlex %s(%#x %d), lval: %+v\n", yySymName(n), n, n, lval)
	}
	return n
}

func yyParse(yylex yyLexer) int {
	const yyError = 73

	yyEx, _ := yylex.(yyLexerEx)
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, 200)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yyerrok := func() {
		if yyDebug >= 2 {
			__yyfmt__.Printf("yyerrok()\n")
		}
		Errflag = 0
	}
	_ = yyerrok
	yystate := 0
	yychar := -1
	var yyxchar int
	var yyshift int
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	if yychar < 0 {
		yylval.yys = yystate
		yychar = yylex1(yylex, &yylval)
		var ok bool
		if yyxchar, ok = yyXLAT[yychar]; !ok {
			yyxchar = len(yySymNames) // > tab width
		}
	}
	if yyDebug >= 4 {
		var a []int
		for _, v := range yyS[:yyp+1] {
			a = append(a, v.yys)
		}
		__yyfmt__.Printf("state stack %v\n", a)
	}
	row := yyParseTab[yystate]
	yyn = 0
	if yyxchar < len(row) {
		if yyn = int(row[yyxchar]); yyn != 0 {
			yyn += yyTabOfs
		}
	}
	switch {
	case yyn > 0: // shift
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		yyshift = yyn
		if yyDebug >= 2 {
			__yyfmt__.Printf("shift, and goto state %d\n", yystate)
		}
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	case yyn < 0: // reduce
	case yystate == 1: // accept
		if yyDebug >= 2 {
			__yyfmt__.Println("accept")
		}
		goto ret0
	}

	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			if yyDebug >= 1 {
				__yyfmt__.Printf("no action for %s in state %d\n", yySymName(yychar), yystate)
			}
			msg, ok := yyXErrors[yyXError{yystate, yyxchar}]
			if !ok {
				msg, ok = yyXErrors[yyXError{yystate, -1}]
			}
			if !ok && yyshift != 0 {
				msg, ok = yyXErrors[yyXError{yyshift, yyxchar}]
			}
			if !ok {
				msg, ok = yyXErrors[yyXError{yyshift, -1}]
			}
			if yychar > 0 {
				ls := yyTokenLiteralStrings[yychar]
				if ls == "" {
					ls = yySymName(yychar)
				}
				if ls != "" {
					switch {
					case msg == "":
						msg = __yyfmt__.Sprintf("unexpected %s", ls)
					default:
						msg = __yyfmt__.Sprintf("unexpected %s, %s", ls, msg)
					}
				}
			}
			if msg == "" {
				msg = "syntax error"
			}
			yylex.Error(msg)
			Nerrs++
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				row := yyParseTab[yyS[yyp].yys]
				if yyError < len(row) {
					yyn = int(row[yyError]) + yyTabOfs
					if yyn > 0 { // hit
						if yyDebug >= 2 {
							__yyfmt__.Printf("error recovery found error shift in state %d\n", yyS[yyp].yys)
						}
						yystate = yyn /* simulate a shift of "error" */
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery failed\n")
			}
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yySymName(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}

			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	r := -yyn
	x0 := yyReductions[r]
	x, n := x0.xsym, x0.components
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= n
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	exState := yystate
	yystate = int(yyParseTab[yyS[yyp].yys][x]) + yyTabOfs
	/* reduction by production r */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce using rule %v (%s), and goto state %d\n", r, yySymNames[x], yystate)
	}

	switch r {
	case 2:
		{
			str, err := TrimQuote(yyS[yypt-0].str)
			if err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
			yyVAL.str = str
		}
	case 3:
		{
			yyVAL.strList = []string{yyS[yypt-0].str}
		}
	case 4:
		{
			yyVAL.strList = append(yyS[yypt-2].strList, yyS[yypt-0].str)
		}
	case 5:
		{
			yyVAL.stmt = Statement(yyS[yypt-0].createStmt)
		}
	case 6:
		{
			yyVAL.stmt = Statement(yyS[yypt-0].selectStmt)
		}
	case 7:
		{
			yyVAL.stmt = Statement(yyS[yypt-0].insertStmt)
		}
	case 8:
		{
			yyVAL.stmt = Statement(yyS[yypt-0].updateStmt)
		}
	case 9:
		{
			yyVAL.stmt = Statement(yyS[yypt-0].deleteStmt)
		}
	case 10:
		{
			yyVAL.stmtList = append(yyVAL.stmtList, yyS[yypt-0].stmt)
		}
	case 11:
		{
			yyVAL.stmtList = append(yyVAL.stmtList, yyS[yypt-0].stmt)
		}
	case 12:
		{
			yyVAL.boolean = true
		}
	case 13:
		{
			yyVAL.boolean = true
		}
	case 14:
		{
			yyVAL.boolean = false
		}
	case 15:
		{
			yyVAL.str = ""
		}
	case 16:
		{
			yyVAL.str = ""
		}
	case 17:
		{
			yyVAL.str = ""
		}
	case 18:
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 19:
		{
			t, ok := typeMapping[yyS[yypt-0].str]
			if ok {
				yyVAL.fieldType = t
			} else {
				__yyfmt__.Printf("不支持的数据类型 %s", yyS[yypt-0].str)
				goto ret1
			}
		}
	case 20:
		{
			yyVAL.createStmt = &CreateStmt{
				Name:   yyS[yypt-5].str,
				Table:  yyS[yypt-3].createTable,
				Option: yyS[yypt-1].createTableOption,
			}
		}
	case 21:
		{
			yyVAL.createTable = &CreateTable{
				Field: []*CreateField{yyS[yypt-0].createField},
				Index: []*CreateIndex{},
			}
		}
	case 22:
		{
			yyVAL.createTable = &CreateTable{
				Field: []*CreateField{},
				Index: []*CreateIndex{yyS[yypt-0].createIndex},
			}
		}
	case 23:
		{
			yyVAL.createTable = &CreateTable{
				Pk:    yyS[yypt-0].createIndex,
				Field: []*CreateField{},
				Index: []*CreateIndex{},
			}
		}
	case 24:
		{
			yyVAL.createTable.Field = append(yyVAL.createTable.Field, yyS[yypt-0].createField)
		}
	case 25:
		{
			yyVAL.createTable.Index = append(yyVAL.createTable.Index, yyS[yypt-0].createIndex)
		}
	case 26:
		{
			if yyVAL.createTable.Pk == nil {
				yyVAL.createTable.Pk = yyS[yypt-0].createIndex
			} else {
				__yyfmt__.Printf("重复定义主键 %v %v ", yyVAL.createTable.Pk, yyS[yypt-0].createIndex)
				goto ret1
			}
		}
	case 27:
		{
			yyVAL.createField = &CreateField{
				Name:       yyS[yypt-3].str,
				Type:       yyS[yypt-2].fieldType,
				AllowNull:  yyS[yypt-1].boolean,
				DefaultVal: yyS[yypt-0].str,
			}
		}
	case 28:
		{
			yyVAL.createIndex = &CreateIndex{
				Name:  yyS[yypt-3].str,
				Field: yyS[yypt-1].str,
			}
		}
	case 29:
		{
			yyVAL.createIndex = &CreateIndex{
				Pk:    true,
				Field: yyS[yypt-1].str,
			}
		}
	case 30:
		{
			yyVAL.createTableOption = nil
		}
	case 31:
		{
			yyVAL.insertStmt = &InsertStmt{
				Table: yyS[yypt-3].str,
				Field: yyS[yypt-2].strList,
				Value: yyS[yypt-1].strList,
			}
		}
	case 32:
		{
			yyVAL.strList = yyS[yypt-1].strList
		}
	case 33:
		{
			yyVAL.strList = nil
		}
	case 35:
		{
			yyVAL.strList = yyS[yypt-1].strList
		}
	case 36:
		{
			yyVAL.strList = nil
		}
	case 38:
		{
			yyVAL.updateStmt = &UpdateStmt{
				Table: yyS[yypt-4].str,
				Value: yyS[yypt-2].updateValue,
				Where: yyS[yypt-1].selectWhereList,
			}
		}
	case 39:
		{
			yyVAL.updateValue = map[string]string{
				yyS[yypt-2].str: yyS[yypt-0].str,
			}
		}
	case 40:
		{
			yyVAL.updateValue[yyS[yypt-2].str] = yyS[yypt-0].str
		}
	case 41:
		{
			yyVAL.deleteStmt = &DeleteStmt{
				Table: yyS[yypt-2].str,
				Where: yyS[yypt-1].selectWhereList,
			}
		}
	case 42:
		{
			yyVAL.boolean = true
		}
	case 43:
		{
			yyVAL.boolean = true
		}
	case 44:
		{
			yyVAL.boolean = false
		}
	case 45:
		{
			yyVAL.compareOperate = EQ
		}
	case 46:
		{
			yyVAL.compareOperate = LT
		}
	case 47:
		{
			yyVAL.compareOperate = GT
		}
	case 48:
		{
			yyVAL.compareOperate = LE
		}
	case 49:
		{
			yyVAL.compareOperate = GE
		}
	case 50:
		{
			yyVAL.compareOperate = NE
		}
	case 51:
		{
			yyVAL.selectStmt = &SelectStmt{
				Field: yyS[yypt-2].selectFieldList,
				Limit: yyS[yypt-1].selectLimit,
			}
		}
	case 52:
		{
			yyVAL.selectStmt = &SelectStmt{
				Table: yyS[yypt-4].str,
				Field: yyS[yypt-6].selectFieldList,
				Where: yyS[yypt-3].selectWhereList,
				Order: yyS[yypt-2].selectOrderList,
				Limit: yyS[yypt-1].selectLimit,
			}
		}
	case 53:
		{
			yyVAL.selectFieldList = []*SelectField{
				&SelectField{
					Name: yyS[yypt-0].str,
				},
			}
		}
	case 54:
		{
			yyVAL.selectFieldList = append(yyS[yypt-2].selectFieldList, &SelectField{
				Name: yyS[yypt-0].str,
			})
		}
	case 55:
		{
			yyVAL.selectWhereList = nil
		}
	case 56:
		{
			yyVAL.selectWhereList = yyS[yypt-0].selectWhereList
		}
	case 57:
		{
			yyVAL.selectWhereList = []SelectWhere{
				&SelectWhereField{
					Field:   yyS[yypt-2].str,
					Value:   yyS[yypt-0].str,
					Operate: yyS[yypt-1].compareOperate,
				},
			}
		}
	case 58:
		{
			yyS[yypt-1].compareOperate.Negate()
			field := &SelectWhereField{
				Field:   yyS[yypt-2].str,
				Value:   yyS[yypt-0].str,
				Operate: yyS[yypt-1].compareOperate,
			}
			if len(yyVAL.selectWhereList) == 1 {
				yyVAL.selectWhereList[0].Negate()
				yyVAL.selectWhereList = append(yyVAL.selectWhereList, field)
				yyVAL.selectWhereList = []SelectWhere{
					&SelectWhereExpr{
						Negation: true,
						Cnf:      yyVAL.selectWhereList,
					},
				}
			} else {
				yyVAL.selectWhereList = []SelectWhere{
					&SelectWhereExpr{
						Negation: true,
						Cnf: []SelectWhere{
							&SelectWhereExpr{
								Negation: true,
								Cnf:      yyVAL.selectWhereList,
							},
							field,
						},
					},
				}
			}
		}
	case 59:
		{
			yyVAL.selectWhereList = append(yyVAL.selectWhereList, &SelectWhereField{
				Field:   yyS[yypt-2].str,
				Value:   yyS[yypt-0].str,
				Operate: yyS[yypt-1].compareOperate,
			})
		}
	case 60:
		{
			expr := &SelectWhereExpr{
				Negation: true,
				Cnf:      yyS[yypt-1].selectWhereList,
			}
			if len(yyVAL.selectWhereList) == 1 {
				yyVAL.selectWhereList[0].Negate()
				yyVAL.selectWhereList = append(yyVAL.selectWhereList, expr)
				yyVAL.selectWhereList = []SelectWhere{
					&SelectWhereExpr{
						Negation: true,
						Cnf:      yyVAL.selectWhereList,
					},
				}
			} else {
				yyVAL.selectWhereList = []SelectWhere{
					&SelectWhereExpr{
						Negation: true,
						Cnf: []SelectWhere{
							&SelectWhereExpr{
								Negation: true,
								Cnf:      yyVAL.selectWhereList,
							},
							expr,
						},
					},
				}
			}
		}
	case 61:
		{
			yyVAL.selectWhereList = append(yyVAL.selectWhereList, yyS[yypt-1].selectWhereList...)
		}
	case 62:
		{
			yyVAL.selectOrderList = nil
		}
	case 63:
		{
			yyVAL.selectOrderList = yyS[yypt-0].selectOrderList
		}
	case 64:
		{
			yyVAL.selectOrderList = []*SelectOrder{
				&SelectOrder{
					Asc:   yyS[yypt-0].boolean,
					Field: yyS[yypt-1].str,
				},
			}
		}
	case 65:
		{
			yyVAL.selectOrderList = append(yyS[yypt-3].selectOrderList, &SelectOrder{
				Asc:   yyS[yypt-0].boolean,
				Field: yyS[yypt-1].str,
			})
		}
	case 66:
		{
			yyVAL.selectLimit = nil
		}
	case 67:
		{
			limit, err := strconv.Atoi(yyS[yypt-0].str)
			if err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
			yyVAL.selectLimit = &SelectLimit{
				Limit: limit,
			}
		}
	case 68:
		{
			limit, err := strconv.Atoi(yyS[yypt-2].str)
			if err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
			offset, err := strconv.Atoi(yyS[yypt-0].str)
			if err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
			yyVAL.selectLimit = &SelectLimit{
				Limit:  limit,
				Offset: offset,
			}
		}
	case 69:
		{
			limit, err := strconv.Atoi(yyS[yypt-2].str)
			if err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
			offset, err := strconv.Atoi(yyS[yypt-0].str)
			if err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
			yyVAL.selectLimit = &SelectLimit{
				Limit:  limit,
				Offset: offset,
			}
		}

	}

	if yyEx != nil && yyEx.Reduced(r, exState, &yyVAL) {
		return -1
	}
	goto yystack /* stack new state and value */
}
