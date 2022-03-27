package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var symbolTable = map[string]int{
	"SP": 0, "LCL": 1, "ARG": 2, "THIS": 3, "THAT": 4,
	"R0": 0, "R1": 1, "R2": 2, "R3": 3,
	"R4": 4, "R5": 5, "R6": 6, "R7": 7,
	"R8": 8, "R9": 9, "R10": 10, "R11": 11,
	"R12": 12, "R13": 13, "R14": 14, "R15": 15,
	"SCREEN": 16384, "KBD": 24576,
}

// Command Types
var A_COMMAND int = 0
var C_COMMAND int = 1
var L_COMMAND int = 2

func openSourceFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	return file
}

func resolveLabels(file *os.File) {
	defer file.Close()
	curInstr := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		curLine := removeWhitespace(scanner.Text())
		if strings.HasPrefix(curLine, "(") && strings.HasSuffix(curLine, ")") {
			symbolTable[getSymbol(curLine)] = curInstr
		} else if !strings.HasPrefix(curLine, "//") {
			curInstr++
		}
	}
}

func parse(cmd string) (map[string]string, int) {
	var parseMap = map[string]string{
		"symb": "",
		"dest": "",
		"comp": "",
		"jump": "",
	}
	cmd = removeWhitespace(cmd)
	cmdType := commandType(cmd)
	switch cmdType {
	case A_COMMAND, L_COMMAND:
		parseMap["symb"] = getSymbol(cmd)
	case C_COMMAND:
		parseMap["dest"] = getDest(cmd)
		parseMap["comp"] = getComp(cmd)
		parseMap["jump"] = getJump(cmd)
	}
	return parseMap, cmdType
}

func commandType(cmd string) int {
	switch cmd[0] {
	case '@':
		return A_COMMAND
	case '(':
		return L_COMMAND
	default:
		return C_COMMAND
	}
}

func getSymbol(cmd string) string {
	if strings.HasPrefix(cmd, "@") {
		return cmd[1:]
	} else if strings.HasPrefix(cmd, "(") {
		return cmd[1 : len(cmd)-1]
	} else {
		return ""
	}
}

// e.g. MD=D+1;JEQ
// return "MD"
func getDest(cmd string) string {
	if strings.Contains(cmd, "=") {
		return strings.Split(cmd, "=")[0]
	} else {
		return ""
	}
}

func getComp(cmd string) string {
	if strings.Contains(cmd, "=") && strings.Contains(cmd, ";") {
		return strings.Split(strings.Split(cmd, ";")[0], "=")[1]
	} else if strings.Contains(cmd, "=") {
		return strings.Split(cmd, "=")[1]
	} else {
		return strings.Split(cmd, ";")[0]
	}
}

func getJump(cmd string) string {
	if strings.Contains(cmd, ";") {
		return strings.Split(cmd, ";")[1]
	} else {
		return ""
	}
}

func removeWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func main() {
	// read asm file
	file := openSourceFile(os.Args[1])

	// first pass
	resolveLabels(file)

	// go line by line and parse
	file = openSourceFile(os.Args[1])
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		curLine := scanner.Text()
		fmt.Println(curLine)
		parseMap, cmdType := parse(curLine)
		for sym, val := range parseMap {
			fmt.Print(sym, ": ", val, "\t")
		}
		fmt.Println("Instruction type was ", cmdType)
		if cmdType == A_COMMAND {

		} else if cmdType == C_COMMAND {

		}
	}
}
