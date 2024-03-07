package telnet

import (
	"bytes"
	"github.com/aybabtme/rgbterm"
	"regexp"
	"strings"
	"unicode"
)

func styleLocation(name, area string) []byte {
	buffer := bytes.Buffer{}
	buffer.WriteString(name)
	if area != "" {
		buffer.WriteString(", ")
		buffer.WriteString(area)
	}
	return rgbterm.FgBytes(buffer.Bytes(), 100, 255, 100)
}

func styleCommand(text string) []byte {
	buffer := bytes.Buffer{}
	buffer.WriteString("\033[1;4m")
	buffer.Write(rgbterm.FgBytes([]byte(strings.TrimSpace(text)), 255, 255, 255))
	buffer.WriteString("\033[0m")
	return buffer.Bytes()
}

func styleSpeech(text string) []byte {
	buffer := bytes.Buffer{}
	buffer.WriteByte('"')
	buffer.Write(rgbterm.FgBytes([]byte(strings.TrimSpace(text)), 255, 255, 255))
	buffer.WriteByte('"')
	return buffer.Bytes()
}

func styleDescription(text string) []byte {
	return rgbterm.FgBytes([]byte(text), 200, 200, 200)
}

func styleHint(text string) []byte {
	return rgbterm.FgBytes([]byte(text), 255, 207, 158)
}

func wrap(maxLength int, s string) string {
	var sb strings.Builder
	var currentWidth int

	var buf bytes.Buffer
	for _, c := range s {
		if unicode.IsSpace(c) {
			// add word to current line if fits

			wordLen := buf.Len() - countANSIEscapeCodes(buf.Bytes())

			if wordLen+currentWidth > maxLength-1 {
				sb.WriteByte(charCR)
				sb.WriteByte(charLF)
				currentWidth = 0
			}

			sb.Write(buf.Bytes())
			currentWidth += wordLen
			buf.Reset()

			if c == ' ' {
				sb.WriteRune(c)
				currentWidth++
			}
			if c == charCR || c == charLF {
				sb.WriteRune(c)
				currentWidth = 0
			}
		} else {
			buf.WriteRune(c)
		}
	}
	sb.Write(buf.Bytes())

	return sb.String()
}

func countANSIEscapeCodes(input []byte) int {
	// Regular expression to match ANSI escape codes
	regex := `\x1b\[[0-9;]*[a-zA-Z]`
	re := regexp.MustCompile(regex)

	// Convert byte slice to string for regex matching
	inputStr := string(input)

	// Find all matches of ANSI escape codes in the input string
	matches := re.FindAllString(inputStr, -1)

	// Calculate the total length of ANSI escape codes
	totalLength := 0
	for _, match := range matches {
		totalLength += len(match)
	}

	return totalLength
}
