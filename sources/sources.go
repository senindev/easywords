package sources

import (
	"server/types"
	"strings"
	"unicode"

	gosrt "github.com/konifar/go-srt"
)

type SrtSource struct {
	fileName string
}

func NewSrtSource(fileName string) *SrtSource {
	return &SrtSource{fileName: fileName}
}
func (s *SrtSource) Source() (types.Storage, error) {
	subtitles, err := gosrt.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}
	text := ""
	for _, subtitle := range subtitles {
		text += subtitle.Text + " "
	}
	return countWords(extractWords(text)), nil
}
func extractWords(s string) []string {
	replaceFunc := func(r rune) rune {
		if r == '?' || r == '.' || r == ',' || r == '!' || r == '\n' {
			return ' '
		}
		return unicode.ToLower(r)
	}
	result := strings.Map(replaceFunc, s)
	words := []string{}
	buf := ""
	for _, r := range result {
		if r == ' ' && buf != "" {
			words = append(words, buf)
			buf = ""
		} else if r != ' ' {
			buf += string(r)
		}
	}
	return words
}
func countWords(words []string) types.Storage {
	wordCount := make(map[string]*types.Unit)
	for _, word := range words {
		if unit, ok := wordCount[word]; ok {
			unit.Frequency++
		} else {
			wordCount[word] = &types.Unit{Expression: word, Frequency: 1}
		}
	}
	ret := types.Storage{}
	for _, unit := range wordCount {
		ret = append(ret, unit)
	}
	return ret
}
