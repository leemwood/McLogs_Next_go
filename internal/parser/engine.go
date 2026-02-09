package parser

import (
	"fmt"
	"mclogs-go/internal/models"
	"os"
	"regexp"
	"sync"

	"gopkg.in/yaml.v3"
)

type Engine struct {
	config models.PatternConfig
}

func NewEngine(patternsPath string) (*Engine, error) {
	data, err := os.ReadFile(patternsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read patterns file: %w", err)
	}

	var config models.PatternConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal patterns: %w", err)
	}

	return &Engine{config: config}, nil
}

func (e *Engine) Analyze(content string) *models.AnalysisResult {
	result := &models.AnalysisResult{
		Name:    "Unknown Log",
		Type:    "unknown",
		Version: "unknown",
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// 1. Detect Log Type
	for _, det := range e.config.Detectors {
		re, err := regexp.Compile(det.Pattern)
		if err != nil {
			continue
		}
		if re.MatchString(content) {
			result.Name = det.Name
			result.Type = det.Type
			break
		}
	}

	// 2. Run Analyzers Concurrently
	wg.Add(len(e.config.Analyzers))
	for _, rule := range e.config.Analyzers {
		go func(r models.AnalyzerRule) {
			defer wg.Done()
			
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return
			}

			matches := re.FindAllStringSubmatch(content, -1)
			if len(matches) > 0 {
				mu.Lock()
				msg := r.Message
				// Dynamic replacement for capturing groups ($1, $2, $3...)
				if len(matches[0]) > 1 {
					for i := 1; i < len(matches[0]); i++ {
						placeholder := fmt.Sprintf("$%d", i)
						msg = regexp.MustCompile(regexp.QuoteMeta(placeholder)).ReplaceAllString(msg, matches[0][i])
					}
				}

				result.Problems = append(result.Problems, models.Problem{
					Severity:  r.Severity,
					Message:   msg,
					Solutions: r.Solutions,
				})
				mu.Unlock()
			}
		}(rule)
	}

	wg.Wait()
	return result
}
