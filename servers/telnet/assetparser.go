package telnet

import (
	"bytes"
	"io"
)

type AssetDescription struct {
	Sections []Section
}

type SectionType byte

const (
	SectionTypeDefault SectionType = iota
	SectionTypeCommand
	SectionTypeSpeech
	SectionTypeHint
)

type Section struct {
	Type SectionType
	Text string
}

func ParseAsset(text string) (AssetDescription, error) {
	var (
		sectionBuffer    bytes.Buffer
		sectionTypeStack []SectionType
		currentType      SectionType
		output           AssetDescription
		reader           *bytes.Reader
	)

	output = AssetDescription{
		Sections: []Section{},
	}

	sectionTypeStack = []SectionType{SectionTypeDefault}

	reader = bytes.NewReader([]byte(text))

	for {
		i, err := reader.ReadByte()
		if err == io.EOF {
			if sectionBuffer.Len() > 0 {
				currentType = sectionTypeStack[len(sectionTypeStack)-1]
				output.Sections = append(output.Sections, Section{
					Type: currentType,
					Text: sectionBuffer.String(),
				})
				break
			}
		}

		switch i {
		case '^':
			if sectionBuffer.Len() > 0 {
				output.Sections = append(output.Sections, Section{
					Type: sectionTypeStack[len(sectionTypeStack)-1],
					Text: sectionBuffer.String(),
				})
				sectionBuffer.Reset()
			}
			if sectionTypeStack[len(sectionTypeStack)-1] != SectionTypeCommand {
				sectionTypeStack = append(sectionTypeStack, SectionTypeCommand)
			} else {
				sectionTypeStack = sectionTypeStack[:len(sectionTypeStack)-1]
			}
		case '"':
			if sectionBuffer.Len() > 0 {
				output.Sections = append(output.Sections, Section{
					Type: sectionTypeStack[len(sectionTypeStack)-1],
					Text: sectionBuffer.String(),
				})
				sectionBuffer.Reset()
			}
			if sectionTypeStack[len(sectionTypeStack)-1] != SectionTypeSpeech {
				sectionTypeStack = append(sectionTypeStack, SectionTypeSpeech)
			} else {
				sectionTypeStack = sectionTypeStack[:len(sectionTypeStack)-1]
			}
		case '@':
			if sectionBuffer.Len() > 0 {
				output.Sections = append(output.Sections, Section{
					Type: sectionTypeStack[len(sectionTypeStack)-1],
					Text: sectionBuffer.String(),
				})
				sectionBuffer.Reset()
			}
			if sectionTypeStack[len(sectionTypeStack)-1] != SectionTypeHint {
				sectionTypeStack = append(sectionTypeStack, SectionTypeHint)
			} else {
				sectionTypeStack = sectionTypeStack[:len(sectionTypeStack)-1]
			}
		default:
			sectionBuffer.WriteByte(i)
		}
	}

	return output, nil
}
