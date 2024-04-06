package sql_generator

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func ParserDDL(ddl string) (result []*Column, tableName, tableComment string) {

	reTable := regexp.MustCompile(`CREATE TABLE +([^ ]+) +\(`)
	reTableComment := regexp.MustCompile(`.+COMMENT='(.+)'$`)
	reIndex := regexp.MustCompile(`(?i)(PRIMARY|UNIQUE)?\s*(INDEX|KEY)\s*(` + "`([^`]*)`" + `)?\s*\(([^)]+)\)`)

	var fieldmap map[string]string = make(map[string]string)
	indexMatches := reIndex.FindAllStringSubmatch(ddl, -1)
	for _, m := range indexMatches {
		idxAttr := strings.Trim(m[5], "`")
		PrefixName := strings.ToUpper(m[1])
		if PrefixName == "PRIMARY" {
			fieldmap[idxAttr] = "primary_key"
		} else if PrefixName == "UNIQUE" {
			fieldmap[idxAttr] = "unique_key"
		} else if PrefixName == "" {
			fieldmap[idxAttr] = "index"
		} else {
			log.Fatal(PrefixName)
		}
	}

	tableMatches := reTable.FindStringSubmatch(ddl)
	tableName = strings.Trim(tableMatches[1], "`")

	tableCommentMatches := reTableComment.FindStringSubmatch(ddl)
	if len(tableCommentMatches) > 0 {
		tableComment = strings.Trim(tableCommentMatches[1], "`")
		tableComment = strings.Trim(tableComment, "\n")
	}

	result = parseDDL(ddl)
	for _, col := range result {
		if v, ok := fieldmap[col.Name]; ok {
			col.IndexType = v
		}
	}

	return result, tableName, tableComment
}

// 自己解析的ddl

func parseDDL(ddl string) (result []*Column) {
	reField := regexp.MustCompile("\n +`[^`]+`.+")

	m := reField.FindAllString(ddl, -1)

	for _, s := range m {
		gcol := &Column{}
		s = strings.Trim(s, " \n,")
		getDDLBlock(s, gcol)
		result = append(result, gcol)
	}
	return
}

func getDDLBlock(_line string, gcol *Column) (result []string) {

	if strings.Contains(_line, "GENERATED") {
		gcol.IsUndefine = true
	}

	var cols [][]rune
	var col = []rune{}
	var curStatus = COLUMN_Name

	var isBackquote bool = true   // ``
	var isSingleQuote bool = true // ''

	line := []rune(_line)

	var i = 0
	for i < len(line) {
		c := line[i]

		if c == ' ' {
			if len(col) != 0 {

				if isBackquote && isSingleQuote {
					colstr := string(col)
					if colstr == "NOT" {
						curStatus = COLUMN_NotNull
						skipChar(&i, line, &col, ' ')
						continue
					} else if colstr == "DEFAULT" {
						curStatus = COLUMN_DefaultValue
						skipChar(&i, line, &col, ' ')
						continue
					} else if colstr == "COMMENT" {
						curStatus = COLUMN_Comment
						skipChar(&i, line, &col, ' ')
						continue
					} else if colstr == "AUTO_INCREMENT" {
						curStatus = COLUMN_AutoIncrement
					} else if colstr == "PRIMARY" || colstr == "UNIQUE" {
						curStatus = COLUMN_IndexType
					} else if colstr == "UNSIGNED" || colstr == "unsigned" {
						curStatus = COLUMN_Unsigned
					}

					switch curStatus {
					case COLUMN_NotNull:
						if colstr == "NOT NULL" {
							gcol.NotNull = true
						}
					case COLUMN_Comment:
						if strings.HasPrefix(colstr, "COMMENT") {
							gcol.Comment = strings.Trim(colstr[7:], " ")
						}
					case COLUMN_DefaultValue:
						if colstr == "DEFAULT NULL" {
							gcol.DefaultValue = nil
						} else {
							dv := fmt.Sprintf("'%s'", colstr[8:])
							gcol.DefaultValue = &dv
						}
					case COLUMN_Name:
						gcol.Name = colstr
					case COLUMN_Type:
						parseType(colstr, gcol)
					case COLUMN_AutoIncrement:
						gcol.AutoIncrement = true
					case COLUMN_IndexType:
						gcol.IndexType = colstr
						if gcol.IndexType == "PRIMARY" {
							gcol.IndexType = "primary_key"
						} else if gcol.IndexType == "UNIQUE" {
							gcol.IndexType = "unique_key"
						}
					case COLUMN_Unsigned:
						gcol.Unsigned = true
					default:
					}

					cols = append(cols, col)
					col = []rune{}
					isBackquote = true   // ``
					isSingleQuote = true // ''
					if curStatus == COLUMN_Name {
						curStatus = COLUMN_Type
					} else {
						curStatus = COLUMN_UNKNOWN
					}
				} else {
					col = append(col, c)
				}

			}

			i++
			continue
		} else {
			if c == '`' {
				if isSingleQuote {
					if curStatus != COLUMN_UNKNOWN {
						curStatus = COLUMN_Name
					}
					isBackquote = !isBackquote
				}
			} else if c == '\'' {
				if isBackquote {
					isSingleQuote = !isSingleQuote
				}
			} else {
				col = append(col, c)
			}
		}

		i++
	}

	if len(col) != 0 {
		if isBackquote && isSingleQuote {
			colstr := string(col)
			if strings.HasPrefix(colstr, "COMMENT") {
				gcol.Comment = strings.Trim(colstr[7:], " ")
			} else if colstr == "AUTO_INCREMENT" {
				gcol.AutoIncrement = true
			} else if curStatus == COLUMN_Type {
				gcol.Type = colstr
			}
			cols = append(cols, col)
		} else {
			panic("语法不合法")
		}
	}

	for _, c := range cols {
		result = append(result, string(c))
	}

	return result
}

func parseType(t string, gcol *Column) {
	line := []rune(t)
	var i = 0
	c := line[i]
	var typeStr []rune
	for c != '(' && i < len(line) {
		typeStr = append(typeStr, c)
		i++
		if i < len(line) {
			c = line[i]
		}

	}
	gcol.Type = string(typeStr)
	if gcol.Type == "json" {
		gcol.Type = "blob"
	}

	typeStr = []rune{}
	i++
	if i < len(line) {
		c = line[i]
	}
	for c != ')' && i < len(line) {
		typeStr = append(typeStr, c)
		i++
		if i < len(line) {
			c = line[i]
		}
	}
	if len(typeStr) != 0 {
		maylen := strings.Split(string(typeStr), ",")
		if len(maylen) >= 1 {
			clen, err := strconv.ParseInt(maylen[0], 10, 64)
			if err != nil {
				panic(err)
			}
			gcol.Length = int(clen)
		}
		if len(maylen) >= 2 {
			clen, err := strconv.ParseInt(maylen[1], 10, 64)
			if err != nil {
				panic(err)
			}
			gcol.Decimal = int(clen)
		}
	}
}

func skipChar(srci *int, line []rune, _col *[]rune, char rune) {
	i := *srci
	var c = line[i]
	var col = *_col
	for c == char && i < len(line) {
		col = append(col, c)
		i++
		if i < len(line) {
			c = line[i]
		}
	}
	*srci = i
	*_col = col
}
