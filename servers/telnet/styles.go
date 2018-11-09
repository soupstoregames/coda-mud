package telnet

import (
	"bytes"
	"github.com/aybabtme/rgbterm"
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
	buffer.Write(rgbterm.FgBytes([]byte(text), 255, 255, 255))
	buffer.WriteString("\033[0m")
	return buffer.Bytes()
}

func styleSpeech(text string) []byte {
	buffer := bytes.Buffer{}
	buffer.WriteByte('"')
	buffer.Write(rgbterm.FgBytes([]byte(text), 255, 255, 255))
	buffer.WriteByte('"')
	return buffer.Bytes()
}

func styleDescription(text string) []byte {
	return rgbterm.FgBytes([]byte(text), 200, 200, 200)
}

func styleHint(text string) []byte {
	return rgbterm.FgBytes([]byte(text), 255, 207, 158)
}
