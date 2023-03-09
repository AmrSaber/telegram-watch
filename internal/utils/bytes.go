package utils

import (
	"bytes"
)

func MeaningfullySplit(input []byte, size int) [][]byte {
	if size < 1 {
		return nil
	}

	input = bytes.TrimSpace(input)

	start := 0
	chunks := make([][]byte, 0, len(input)/size)

	for start < len(input) {
		end := start + size
		if end >= len(input) {
			chunks = append(chunks, input[start:])
			break
		}

		// If the next character is a delimiter, take it
		// this is to avoid adding "-" at the end when the next character is a delimiter anyway
		next := input[end]
		if next == ' ' || next == '\n' {
			end++
		}

		chunk := input[start:end]
		cutWord := false

		// Try to find a new line within the limit
		length := bytes.LastIndex(chunk, []byte{'\n'})

		// If no new line found, try to find a space
		if length == -1 {
			length = bytes.LastIndex(chunk, []byte{' '})

			// If no space found, then just split the text
			if length == -1 {
				length = len(chunk) - 1 // leave space for "-" character that will be appended
				cutWord = true
			}
		}

		chunk = chunk[:length]
		start += length
		if cutWord {
			chunk = append(chunk, '-')
		} else {
			// Ignore the space that we stopped at
			start++
		}

		chunks = append(chunks, chunk)
	}

	return chunks
}
