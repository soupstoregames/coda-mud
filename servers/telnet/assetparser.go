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

type Parser struct {
	sectionBuffer    bytes.Buffer
	sectionTypeStack []SectionType
	reader           *bytes.Reader
}

func (p *Parser) Parse(text string) (AssetDescription, error) {
	output := AssetDescription{
		Sections: []Section{},
	}

	p.sectionTypeStack = []SectionType{SectionTypeDefault}
	p.sectionBuffer = bytes.Buffer{}
	p.reader = bytes.NewReader([]byte(text))

	for {
		i, err := p.reader.ReadByte()
		if err == io.EOF {
			if p.sectionBuffer.Len() > 0 {
				output.Sections = append(output.Sections, Section{
					Type: p.sectionTypeStack[len(p.sectionTypeStack)-1],
					Text: p.sectionBuffer.String(),
				})
				break
			}
		}

		switch i {
		case '^':
			output.Sections = append(output.Sections, p.flushSectionBuffer())
			p.transitionSection(SectionTypeCommand)
		case '"':
			output.Sections = append(output.Sections, p.flushSectionBuffer())
			p.transitionSection(SectionTypeSpeech)
		case '@':
			output.Sections = append(output.Sections, p.flushSectionBuffer())
			p.transitionSection(SectionTypeHint)
		default:
			p.sectionBuffer.WriteByte(i)
		}
	}

	return output, nil
}

func (p *Parser) flushSectionBuffer() Section {
	var result Section

	if p.sectionBuffer.Len() > 0 {
		result = Section{
			Type: p.sectionTypeStack[len(p.sectionTypeStack)-1],
			Text: p.sectionBuffer.String(),
		}
		p.sectionBuffer.Reset()
	}

	return result
}

func (p *Parser) transitionSection(sectionType SectionType) {
	if p.sectionTypeStack[len(p.sectionTypeStack)-1] != sectionType {
		p.sectionTypeStack = append(p.sectionTypeStack, sectionType)
	} else {
		p.sectionTypeStack = p.sectionTypeStack[:len(p.sectionTypeStack)-1]
	}
}
