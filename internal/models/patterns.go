package models

type PatternConfig struct {
	Detectors []DetectorPattern `yaml:"detectors"`
	Analyzers []AnalyzerRule    `yaml:"analyzers"`
}

type DetectorPattern struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Pattern string `yaml:"pattern"`
}

type AnalyzerRule struct {
	Name      string     `yaml:"name"`
	Pattern   string     `yaml:"pattern"`
	Severity  string     `yaml:"severity"`
	Message   string     `yaml:"message"`
	Solutions []Solution `yaml:"solutions"`
}
