package rake

import (
	"io/ioutil"
	"sort"
	"strings"
	"unicode"
)

// Score : (Word, Score) pair
type score struct {
	word  string
	score float64
}

type byScore []score

func (s byScore) Len() int {
	return len(s)
}

func (s byScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byScore) Less(i, j int) bool {
	return s[i].score > s[j].score
}

func getTextFromFile(filename string) string {
	content, _ := ioutil.ReadFile(filename)
	return string(content)
}

func getLinesFromFile(filename string) []string {
	content, _ := ioutil.ReadFile(filename)
	return strings.Split(string(content), "\n")
}

func splitIntoWords(text string) []string {
	words := []string{}
	splitWords := strings.Fields(text)
	for _, word := range splitWords {
		currentWord := strings.ToLower(strings.TrimSpace(word))
		if currentWord != "" {
			words = append(words, currentWord)
		}
	}
	return words
}

func getStopwords() map[string]bool {
	stopwords := stopwords
	dict := map[string]bool{}
	for _, word := range stopwords {
		dict[word] = true
	}
	return dict
}

func generateCandidatePhrases(text string) []string {
	stopwords := getStopwords()
	words := splitIntoWords(text)
	acceptedWords := []string{}
	for _, word := range words {
		if !stopwords[word] {
			acceptedWords = append(acceptedWords, word)
		} else {
			acceptedWords = append(acceptedWords, "|")
		}
	}

	phraseList := []string{}
	phrase := ""
	for _, word := range acceptedWords {
		if word == "|" {
			phraseList = append(phraseList, phrase)
			phrase = ""
		} else {
			phrase = phrase + " " + word
		}
	}
	return phraseList
}

func splitIntoSentences(text string) []string {
	splitFunc := func(c rune) bool {
		return unicode.IsPunct(c)
	}
	return strings.FieldsFunc(text, splitFunc)
}

func combineScores(phraseList []string, scores map[string]float64) map[string]float64 {
	candidateScores := map[string]float64{}
	for _, phrase := range phraseList {
		words := splitIntoWords(phrase)
		candidateScore := float64(0.0)

		for _, word := range words {
			candidateScore += scores[word]
		}
		candidateScores[phrase] = candidateScore
	}
	return candidateScores
}

func calculateWordScores(phraseList []string) map[string]float64 {
	frequencies := map[string]int{}
	degrees := map[string]int{}
	for _, phrase := range phraseList {
		words := splitIntoWords(phrase)
		length := len(words)
		degree := length - 1

		for _, word := range words {
			frequencies[word]++
			degrees[word] += degree
		}
	}
	for key := range frequencies {
		degrees[key] = degrees[key] + frequencies[key]
	}

	score := map[string]float64{}

	for key := range frequencies {
		score[key] += (float64(degrees[key]) / float64(frequencies[key]))
	}

	return score
}

func sortScores(scores map[string]float64, topN int) []score {
	rakeScores := []score{}
	for k, v := range scores {
		rakeScores = append(rakeScores, score{k, v})
	}
	sort.Sort(byScore(rakeScores))
	if topN < len(rakeScores) && topN > 0 {
		return rakeScores[0:topN]
	}
	return rakeScores
}

func rake(text string, topN int) map[string]float64 {
	sentences := splitIntoSentences(text)
	phraseList := []string{}
	for _, sentence := range sentences {
		phraseList = append(phraseList, generateCandidatePhrases(sentence)...)
	}
	wordScores := calculateWordScores(phraseList)
	candidateScores := combineScores(phraseList, wordScores)
	sortedScores := sortScores(candidateScores, topN)
	scoreDict := make(map[string]float64)
	for _, score := range sortedScores {
		scoreDict[strings.TrimSpace(score.word)] = score.score
	}
	return scoreDict
}

// WithFile : Run rake with text from a file
func WithFile(filename string) map[string]float64 {
	text := getTextFromFile(filename)
	return rake(text, 10)
}

// WithText : Run rake directly from text
func WithText(text string) map[string]float64 {
	return rake(text, 10)
}

// TopNWithText : Run rake directly from text and return the top N results
func TopNWithText(text string, topN int) map[string]float64 {
	return rake(text, topN)
}
