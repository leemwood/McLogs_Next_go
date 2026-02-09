package filter

type Filter interface {
	Filter(content string) (string, error)
}

type LengthFilter struct {
	MaxLength int
}

func (f *LengthFilter) Filter(content string) (string, error) {
	if len(content) > f.MaxLength {
		return content[:f.MaxLength], nil
	}
	return content, nil
}

type LinesFilter struct {
	MaxLines int
}

func (f *LinesFilter) Filter(content string) (string, error) {
	// Simple implementation, could be optimized for very large files
	lines := 0
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			lines++
			if lines >= f.MaxLines {
				return content[:i], nil
			}
		}
	}
	return content, nil
}
