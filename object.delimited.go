package tox

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type part struct {
	name        string
	dataType    string
	timeFormat  string
	arraySize   int
	arrayObjKey string
}

func parsePart(partString string, idx int) (part, error) {

	if len(partString) == 0 {
		return part{name: fmt.Sprintf("ignored.field%d", idx), dataType: "string"}, nil
	}

	lp := strings.IndexByte(partString, '(')
	rp := strings.IndexByte(partString, ')')

	if lp == -1 {
		return part{name: partString, dataType: "string"}, nil
	}

	if rp == -1 {
		return part{name: partString[0:lp]}, errors.New("missing parentheses")
	}

	if rp < lp {
		return part{name: partString[0:lp], dataType: "string"}, errors.New("bad format")
	}

	if rp-lp == 1 {
		return part{name: partString[0:lp], dataType: "string"}, nil
	}

	argStr := partString[lp+1 : rp]
	args := strings.Split(argStr, ",")

	ret := part{name: partString[0:lp]}

	if len(args) >= 1 {
		switch strings.ToLower(args[0]) {
		case "string", "int", "float", "bool", "timestamp", "array", "arrayobj":
			ret.dataType = strings.ToLower(args[0])
		default:
			return ret, errors.New("invalid data type")
		}
	} else {
		ret.dataType = "string"
	}

	if len(args) >= 2 {
		switch ret.dataType {
		case "timestamp":
			ret.timeFormat = args[1]
		case "array", "arrayobj":
			ret.arraySize = ToInt(args[1])
		}
	}
	if len(args) >= 3 {
		switch ret.dataType {
		case "arrayobj":
			ret.arrayObjKey = ToString(args[2])
		}
	}
	if ret.dataType == "arrayobj" && len(ret.arrayObjKey) == 0 {
		return ret, errors.New("arrayobj requires 2 parameters")
	}

	return ret, nil
}

func NewObjectFromDelimitedString(source string, pattern string, delim string) (Object, error) {
	patternParts := strings.Split(pattern, "|")
	patternIdx := 0

	data := Object{}

	var arrayBase string
	var arrayCount int
	var arrayIndex int
	var arrayPartIndex int
	var arrayPartCount int
	var arrayObjKey string
	var arrayData Object
	var arrayItems []Object

	sourceParts := strings.Split(source, delim)
	for _, sourcePart := range sourceParts {

		var partData Object
		var parsedPart part

		patternPart := ""

		if patternIdx < len(patternParts) {
			patternPart = patternParts[patternIdx]
		}
		parsedPart, err := parsePart(patternPart, patternIdx)
		if err != nil {
			return nil, err
		}
		patternIdx++

		if len(arrayBase) != 0 {
			partData = arrayData
		} else {
			partData = data
		}

		if len(sourcePart) != 0 {
			switch parsedPart.dataType {
			case "string":
				partData.Set(parsedPart.name, sourcePart)
			case "int":
				partData.Set(parsedPart.name, ToInt(sourcePart))
			case "float":
				partData.Set(parsedPart.name, ToFloat64(sourcePart))
			case "bool":
				partData.Set(parsedPart.name, ToBool(sourcePart))
			case "timestamp":
				tm, err := time.Parse(parsedPart.timeFormat, sourcePart)
				if err != nil {
					partData.Set(parsedPart.name, "err: "+err.Error())
				} else {
					partData.Set(parsedPart.name, tm.UTC().Format(time.RFC3339))
				}
			case "array", "arrayobj":
				arrayBase = parsedPart.name
				arrayCount = ToInt(sourcePart)
				arrayIndex = 0
				arrayPartCount = parsedPart.arraySize
				arrayPartIndex = -1
				arrayData = Object{}
				if parsedPart.dataType == "arrayobj" {
					arrayObjKey = parsedPart.arrayObjKey
				} else {
					arrayObjKey = ""
				}
			}
		}

		if len(arrayBase) != 0 {
			arrayPartIndex++
			if arrayPartIndex == arrayPartCount {
				// done with array item
				arrayItems = append(arrayItems, arrayData)
				arrayData = Object{}
				arrayPartIndex = 0

				arrayIndex++
				if arrayIndex == arrayCount {
					// done with array
					if len(arrayObjKey) != 0 {
						subObj := Object{}
						for _, item := range arrayItems {
							key := item.GetString(arrayObjKey, "")
							if len(key) == 0 {
								return nil, errors.New("bad array: missing object key")
							}
							delete(item, arrayObjKey)
							subObj.Set(key, item)
						}
						data.Set(arrayBase, subObj)
					} else {
						data.Set(arrayBase, arrayItems)
					}

					arrayBase = ""
				} else {
					patternIdx -= arrayPartCount
				}
			}
		}
	}
	return data, nil
}
