package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

var destCode = map[string]string{
	"": "000", "M": "001", "D": "010",
	"MD": "011", "A": "100", "AM": "101",
	"AD": "110", "AMD": "111",
}

var compCode = map[string]string{
	"0": "0101010", "1": "0111111",
	"-1": "0111010", "D": "0001100",
	"A": "0110000", "M": "1110000",
	"!D": "0001101", "!A": "0110001",
	"!M": "1110001", "-D": "0001111",
	"-A": "0110011", "-M": "1110011",
	"D+1": "0011111", "A+1": "0110111",
	"M+1": "1110111", "D-1": "0001110",
	"A-1": "0110010", "M-1": "1110010",
	"D+A": "0000010", "D+M": "1000010",
	"D-A": "0010011", "D-M": "1010011",
	"A-D": "0000111", "M-D": "1000111",
	"D&A": "0000000", "D&M": "1000000",
	"D|A": "0010101", "D|M": "1010101",
}

var jumpCode = map[string]string{
	"": "000", "JGT": "001",
	"JEQ": "010", "JGE": "011",
	"JLT": "100", "JNE": "101",
	"JLE": "110", "JMP": "111",
}

// Command Types
var A_COMMAND int = 0
var C_COMMAND int = 1
var L_COMMAND int = 2
var N_COMMAND int = 2 // not a command

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
		} else if !strings.HasPrefix(curLine, "//") && curLine != "" {
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
	cmd = removeInlineComment(cmd)
	if cmd == "" || strings.HasPrefix(cmd, "//") {
		return parseMap, N_COMMAND
	}
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

func removeInlineComment(str string) string {
	if strings.Contains(str, "//") {
		return strings.Split(str, "//")[0]
	}
	return str
}

func removeWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func strToBin(n int) string {
	return fmt.Sprintf("%015b", n)
}

func writeHackFile(fileName string, hackCode string) {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		fileName = fileName[:pos]
	}
	f, err := os.Create(fileName + ".hack")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(hackCode)
}

func main() {
	hackCode := ""
	// variable counter starts from 16
	varCount := 16
	// read asm file
	file := openSourceFile(os.Args[1])

	// first pass
	resolveLabels(file)

	// go line by line and parse
	file = openSourceFile(os.Args[1])
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		instruction := ""
		curLine := scanner.Text()
		parseMap, cmdType := parse(curLine)
		if cmdType == A_COMMAND {
			instruction += "0"
			symb := parseMap["symb"]
			n, err := strconv.Atoi(symb)
			if err != nil {
				// symb is not a number
				// check if symb is in symbolTable
				if _, ok := symbolTable[symb]; !ok {
					// symbol doesn't exist in symbolTable
					symbolTable[symb] = varCount
					varCount++
				}
				n = symbolTable[symb]
			}
			instruction += strToBin(n)
			hackCode += instruction + "\n"
		} else if cmdType == C_COMMAND {
			instruction += "111"
			instruction += compCode[parseMap["comp"]]
			instruction += destCode[parseMap["dest"]]
			instruction += jumpCode[parseMap["jump"]]
			hackCode += instruction + "\n"
		}
	}
	writeHackFile(os.Args[1], hackCode)
}
