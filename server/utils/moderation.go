package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ModerationResult represents the result of content moderation
type ModerationResult struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

// Keyword-based moderation
var bannedKeywords = []string{
	// Common slurs and hateful terms
	"nigger", "nigga", "coon", "chink", "gook", "spic", "wetback", "beaner",
	"kike", "heeb", "raghead", "sandnigger", "towelhead", "paki", "currymuncher",
	"jap", "slope", "zipperhead", "faggot", "fag", "homo", "queer", "tranny",
	"cunt", "whore", "slut", "bitch", "pussy", "dick", "cock", "asshole",
	"fuck", "shit", "bastard", "motherfucker", "cocksucker",
	// Hateful terms
	"racist", "nazi", "hitler", "supremacist", "terrorist", "pedophile",
	"rapist", "murderer", "killer", "suicide", "kill yourself", "die",
	// Common bypass attempts
	"n1gger", "n1gga", "nigg3r", "nigg3r", "f4ggot", "f4g", "c0ck", "d1ck",
	"f*ck", "sh*t", "b*tch", "c*nt", "p*ssy", "a**hole", "m*therf*cker",
	// Leetspeak variations
	"phaggot", "phag", "cocksuka", "cocksucker", "motherfuka", "motherfucker",
}

// CheckKeywordModeration performs keyword-based content moderation
func CheckKeywordModeration(content string) *ModerationResult {
	contentLower := strings.ToLower(content)

	// Check for banned keywords (substring matching to catch variations)
	for _, keyword := range bannedKeywords {
		if strings.Contains(contentLower, keyword) {
			return &ModerationResult{
				Allowed: false,
			}
		}
	}

	// Check for repeated characters (potential bypass)
	if hasRepeatedChars(content) {
		return &ModerationResult{
			Allowed: false,
		}
	}

	// Check for excessive caps
	if hasExcessiveCaps(content) {
		return &ModerationResult{
			Allowed: false,
		}
	}

	return &ModerationResult{Allowed: true}
}

// hasRepeatedChars checks for suspicious repeated characters
func hasRepeatedChars(content string) bool {
	// Check for 5 or more consecutive identical characters
	for i := 0; i < len(content)-4; i++ {
		if content[i] == content[i+1] && content[i] == content[i+2] &&
			content[i] == content[i+3] && content[i] == content[i+4] {
			return true
		}
	}
	return false
}

// hasExcessiveCaps checks for excessive capitalization
func hasExcessiveCaps(content string) bool {
	if len(content) < 10 {
		return false
	}

	capsCount := 0
	totalLetters := 0

	for _, char := range content {
		if char >= 'A' && char <= 'Z' {
			capsCount++
			totalLetters++
		} else if char >= 'a' && char <= 'z' {
			totalLetters++
		}
	}

	if totalLetters == 0 {
		return false
	}

	return float64(capsCount)/float64(totalLetters) > 0.8 // 80% caps
}

// LLMModerationRequest represents a request to the LLM moderation service
type LLMModerationRequest struct {
	Content string `json:"content"`
}

// LLMModerationResponse represents the response from the LLM moderation service
type LLMModerationResponse struct {
	Toxic      bool    `json:"toxic"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason,omitempty"`
}

// CheckLLMModeration performs LLM-based content moderation
func CheckLLMModeration(content string) *ModerationResult {
	llmEndpoint := os.Getenv("LLM_MODERATION_ENDPOINT")
	if llmEndpoint == "" {
		// If no LLM endpoint configured, skip LLM check
		return &ModerationResult{Allowed: true}
	}

	reqBody := LLMModerationRequest{Content: content}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		Debugf("Failed to marshal LLM request: %v", err)
		return &ModerationResult{Allowed: true} // Allow on error
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(llmEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		Debugf("LLM moderation request failed: %v", err)
		return &ModerationResult{Allowed: true} // Allow on error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Debugf("LLM moderation returned status %d", resp.StatusCode)
		return &ModerationResult{Allowed: true} // Allow on error
	}

	var llmResp LLMModerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		Debugf("Failed to decode LLM response: %v", err)
		return &ModerationResult{Allowed: true} // Allow on error
	}

	if llmResp.Toxic && llmResp.Confidence > 0.7 { // 70% confidence threshold
		return &ModerationResult{
			Allowed: false,
			Reason:  fmt.Sprintf("LLM detected toxic content: %s", llmResp.Reason),
		}
	}

	return &ModerationResult{Allowed: true}
}

// CheckMarkdownModeration filters out non-inline markdown content
func CheckMarkdownModeration(content string) *ModerationResult {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Reject headers (# ## ### etc.)
		if strings.HasPrefix(line, "#") {
			return &ModerationResult{
				Allowed: false,
				Reason:  "Headers are not allowed",
			}
		}

		// Reject blockquotes (> at start of line)
		if strings.HasPrefix(line, ">") {
			return &ModerationResult{
				Allowed: false,
				Reason:  "Blockquotes are not allowed",
			}
		}

		// Reject unordered lists (- or * at start of line)
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			return &ModerationResult{
				Allowed: false,
				Reason:  "Lists are not allowed",
			}
		}

		// Reject ordered lists (1. 2. etc.)
		orderedListPattern := regexp.MustCompile(`^\d+\.\s`)
		if orderedListPattern.MatchString(line) {
			return &ModerationResult{
				Allowed: false,
				Reason:  "Ordered lists are not allowed",
			}
		}

		// Reject horizontal rules (--- or *** or ___)
		if line == "---" || line == "***" || line == "___" {
			return &ModerationResult{
				Allowed: false,
				Reason:  "Horizontal rules are not allowed",
			}
		}

		// Reject table rows (containing | and potentially - for separators)
		if strings.Contains(line, "|") && (strings.Contains(line, "---") || strings.Contains(line, ":-") || strings.Contains(line, "-:")) {
			return &ModerationResult{
				Allowed: false,
				Reason:  "Tables are not allowed",
			}
		}
	}

	// Reject code blocks (```)
	if strings.Contains(content, "```") {
		return &ModerationResult{
			Allowed: false,
			Reason:  "Code blocks are not allowed",
		}
	}

	// Reject images (![alt](url) or ![alt][ref])
	imagePattern := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	if imagePattern.MatchString(content) {
		return &ModerationResult{
			Allowed: false,
			Reason:  "Images are not allowed",
		}
	}

	imageRefPattern := regexp.MustCompile(`!\[.*?\]\[.*?\]`)
	if imageRefPattern.MatchString(content) {
		return &ModerationResult{
			Allowed: false,
			Reason:  "Images are not allowed",
		}
	}

	// Allow links - [text](url) and [text][ref] are OK
	// Allow bare URLs
	// Allow inline formatting: *italic*, **bold**, ***bold italic***, `code`, ~~strikethrough~~

	return &ModerationResult{Allowed: true}
}

// ModerateContent performs comprehensive content moderation
func ModerateContent(content string) *ModerationResult {
	// First, check markdown formatting
	markdownResult := CheckMarkdownModeration(content)
	if !markdownResult.Allowed {
		return markdownResult
	}

	// Then, check keywords
	keywordResult := CheckKeywordModeration(content)
	if !keywordResult.Allowed {
		return &ModerationResult{
			Allowed: false,
			Reason:  "Content contains banned keywords or patterns",
		}
	}

	// Only check with LLM if endpoint is configured
	llmEndpoint := os.Getenv("LLM_MODERATION_ENDPOINT")
	if llmEndpoint != "" {
		llmResult := CheckLLMModeration(content)
		if !llmResult.Allowed {
			return llmResult
		}
	}

	return &ModerationResult{Allowed: true}
}
