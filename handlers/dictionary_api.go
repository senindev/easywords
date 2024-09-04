package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"server/types"
	"strings"
	"sync"
	"time"
)

type definition struct {
	Definition string   `json:"definition"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
	Example    string   `json:"example"`
}
type meaning struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Definitions  []definition
}

type successResponse struct {
	Word     string    `json:"word"`
	Phonetic string    `json:"phonetic"`
	Meanings []meaning `json:"meanings"`
}
type errorResponse struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	Resolution string `json:"resolution"`
}
type DictionaryAPIHandler struct {
}

func NewDictionaryAPIHandler() *DictionaryAPIHandler {
	return &DictionaryAPIHandler{}
}
func safeSendRequest(url string, retry int) (*http.Response, error) {
	const sleepMin = 500
	const sleepMax = 2500
	var resp *http.Response
	for {
		retry--
		r, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if r.StatusCode == http.StatusTooManyRequests {
			if retry == 0 {
				return nil, fmt.Errorf("Failed to get response after many retries")
			}
			fmt.Println("sleep!")
			time.Sleep(time.Duration(rand.IntN(sleepMax-sleepMin)+sleepMin) * time.Millisecond)
			continue
		}
		resp = r
		break
	}
	return resp, nil
}
func getTranslation(unit *types.Unit, retry int, wg *sync.WaitGroup) error {
	defer wg.Done()

	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", unit.Expression)
	resp, err := safeSendRequest(url, retry)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Word not found in dictionary")
	}
	decoder := json.NewDecoder(resp.Body)

	var result []successResponse
	if err := decoder.Decode(&result); err != nil {
		if err == io.EOF {
			log.Println("Empty response body")
		} else {
			log.Fatalf("Error decoding JSON: %v", err)
		}
	}
	for _, r := range result {
		unit.Translates = append(unit.Translates, r.Meanings[0].Definitions[0].Definition)
	}
	return nil
}
func (d *DictionaryAPIHandler) Handle(s types.Storage) error {
	const blockSize = 10

	var wg sync.WaitGroup
	for i := 0; i < len(s); i += blockSize {
		for j := i; j < i+blockSize && j < len(s); j++ {
			unit := s[j]
			if strings.Contains(unit.Expression, "'") {
				continue
			}
			wg.Add(1)
			fmt.Println(unit)
			go getTranslation(unit, 3, &wg)
		}
		wg.Wait()
	}
	return nil
}
