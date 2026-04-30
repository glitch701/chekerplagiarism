package chunker

import "strings"

type Chunker struct {
	size   int
	minLen int
}

func New(size, _ /*overlap*/, minLen int) *Chunker {
	return &Chunker{size: size, minLen: minLen}
}

// Chunk splits text into sentence-based chunks grouped up to size bytes.
func (c *Chunker) Chunk(text string) []string {
	if len(text) == 0 {
		return nil
	}

	sentences := splitSentences(text)

	var chunks []string
	var buf strings.Builder

	for _, s := range sentences {
		s = strings.Join(strings.Fields(s), " ")
		if len(s) < 2 {
			continue
		}

		if buf.Len() > 0 && buf.Len()+1+len(s) > c.size {
			if chunk := buf.String(); len(chunk) >= c.minLen {
				chunks = append(chunks, chunk)
			}
			buf.Reset()
		}

		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(s)
	}

	if chunk := buf.String(); len(chunk) >= c.minLen {
		chunks = append(chunks, chunk)
	}

	return chunks
}

// splitSentences splits text on ". ", "! ", "? " and newlines.
func splitSentences(text string) []string {
	var result []string
	var buf strings.Builder
	runes := []rune(text)
	n := len(runes)

	for i := 0; i < n; i++ {
		r := runes[i]
		buf.WriteRune(r)

		end := r == '\n'
		if !end && (r == '.' || r == '!' || r == '?') && i+1 < n {
			next := runes[i+1]
			end = next == ' ' || next == '\n' || next == '"' || next == ')'
		}

		if end {
			if s := strings.TrimSpace(buf.String()); s != "" {
				result = append(result, s)
			}
			buf.Reset()
		}
	}

	if s := strings.TrimSpace(buf.String()); s != "" {
		result = append(result, s)
	}

	return result
}
