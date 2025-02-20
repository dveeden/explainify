package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"unicode/utf8"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"golang.org/x/text/width"
)

var tableFormat = regexp.MustCompile(`^\+-+[-\+]*\+$`)
var isExplainTable = regexp.MustCompile(`^\| (EXPLAIN|TiDB_JSON)\s*\|$`)

func unicodeTable(output [][]byte) [][]byte {
	if len(output) == 0 {
		return output
	}

	if !tableFormat.Match(output[0]) {
		return output
	}

	var columns []int
	for i, c := range output[0] {
		if c == '+' {
			columns = append(columns, i)
		}
	}

	for i, line := range output {
		if tableFormat.Match(line) {
			output[i] = bytes.ReplaceAll(line, []byte("-"), []byte("─"))
			if i == 0 {
				output[i] = slices.Concat([]byte("╭"), output[i][1:len(output[i])-1], []byte("╮"))
				output[i] = bytes.ReplaceAll(output[i], []byte("+"), []byte("┬"))
			} else if i == len(output)-1 {
				output[i] = slices.Concat([]byte("╰"), output[i][1:len(output[i])-1], []byte("╯"))
				output[i] = bytes.ReplaceAll(output[i], []byte("+"), []byte("┴"))
			} else {
				output[i] = slices.Concat([]byte("├"), output[i][1:len(output[i])-1], []byte("┤"))
				output[i] = bytes.ReplaceAll(output[i], []byte("+"), []byte("┼"))
			}
		} else {
			var b bytes.Buffer
			tmpBytes := make([]byte, 4)
			skip := 0
			skipcnt := 0
			for j, c := range bytes.Runes(output[i]) {
				n := utf8.EncodeRune(tmpBytes, c)
				if slices.Contains(columns, j-skipcnt) {
					b.WriteRune('│')
				} else {
					if skip > 0 && c == ' ' {
						skip--
						if skip%2 == 1 {
							skipcnt++
						}
						continue
					}
					b.WriteRune(c)
				}
				if n > 3 {
					skip++
				}
				if width.LookupRune(c).Kind() == width.EastAsianWide {
					skip++
				}
			}
			output[i] = b.Bytes()
		}
	}

	// make the header bold
	output[1] = slices.Concat([]byte("\033[1m"), output[1], []byte("\033[0m"))

	return output
}

func mysqlExplain(output [][]byte, theme string) [][]byte {
	header := [][]byte{
		[]byte("          EXPLAIN"),
		bytes.Repeat([]byte("-"), 80),
	}

	if !isExplainTable.Match(output[1]) {
		return output
	}

	var newOutput [][]byte
	for i := 3; i < len(output)-1; i++ {
		if output[i][0] == '|' {
			newOutput = append(newOutput, bytes.TrimRight(output[i][2:], "|"))
		} else {
			newOutput = append(newOutput, bytes.TrimRight(output[i], "|"))
		}
	}

	if json.Valid(bytes.Join(newOutput, []byte("\n"))) {
		var b bytes.Buffer
		fmtOutput := bufio.NewWriter(&b)
		lexer := lexers.Get("json")
		style := styles.Get(theme)
		if style == nil {
			style = styles.Fallback
		}
		formatter := formatters.Get("terminal")
		if formatter == nil {
			formatter = formatters.Fallback
		}
		iterator, _ := lexer.Tokenise(nil, string(bytes.Join(newOutput, []byte("\n"))))
		_ = formatter.Format(fmtOutput, style, iterator)
		fmtOutput.Flush()
		return slices.Concat(header, bytes.Split(b.Bytes(), []byte("\n")))
	}

	return slices.Concat(header, newOutput)
}

func markdownTable(output [][]byte) [][]byte {
	if !tableFormat.Match(output[0]) {
		return output
	}

	output[2] = bytes.ReplaceAll(output[2], []byte("+"), []byte("|"))

	return output[1 : len(output)-1]
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var outfmt = flag.String("format","unicode","output format: plain, markdown or unicode")
	var theme = flag.String("theme", "monokailight", "chroma syntax highlighting theme")
	flag.Parse()

	parts := bytes.Split(bytes.TrimRight(data, "\n"), []byte("\n"))

	parts = mysqlExplain(parts, *theme)

	switch (*outfmt) {
	case "unicode":
		parts = unicodeTable(parts)
	case "markdown":
		parts = markdownTable(parts)
	case "plain":
	default:
		panic("unknown output format")
	}

	for _, part := range parts {
		fmt.Printf("%s\n", part)
	}
}
