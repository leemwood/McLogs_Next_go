package parser

import (
	"mclogs-go/internal/models"
)

type Parser struct {
	engine *Engine
}

func NewParser(patternsPath string) (*Parser, error) {
	engine, err := NewEngine(patternsPath)
	if err != nil {
		return nil, err
	}
	return &Parser{engine: engine}, nil
}

func (p *Parser) Parse(content string) *models.AnalysisResult {
	return p.engine.Analyze(content)
}
