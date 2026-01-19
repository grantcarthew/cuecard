package prompt

import (
	"regexp"
	"strings"
	"time"
)

// VariableResolver provides values for variable substitution
type VariableResolver struct {
	Input        string
	Clipboard    string
	FileSelector func() string // Called when ${FILE} is encountered
}

// Substitute replaces variables in the content with their values
func Substitute(content string, resolver *VariableResolver) string {
	if resolver == nil {
		resolver = &VariableResolver{}
	}

	// Define variable patterns and their replacements
	result := content

	// ${INPUT} - User-provided text
	result = strings.ReplaceAll(result, "${INPUT}", resolver.Input)

	// ${DATE} - Current date in ISO format
	result = strings.ReplaceAll(result, "${DATE}", time.Now().Format("2006-01-02"))

	// ${DATETIME} - Current date and time
	result = strings.ReplaceAll(result, "${DATETIME}", time.Now().Format("2006-01-02 15:04:05"))

	// ${CLIPBOARD} - Current clipboard content
	result = strings.ReplaceAll(result, "${CLIPBOARD}", resolver.Clipboard)

	// ${FILE} - File picker result
	if strings.Contains(result, "${FILE}") && resolver.FileSelector != nil {
		filePath := resolver.FileSelector()
		result = strings.ReplaceAll(result, "${FILE}", filePath)
	}

	return result
}

// SubstituteWithValues replaces variables using a map of values
func SubstituteWithValues(content string, values map[string]string) string {
	result := content
	for key, value := range values {
		placeholder := "${" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// ExtractVariables returns all unique variable names from content
func ExtractVariables(content string) []string {
	re := regexp.MustCompile(`\$\{([A-Z_]+)\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	var vars []string
	for _, match := range matches {
		if len(match) > 1 && !seen[match[1]] {
			seen[match[1]] = true
			vars = append(vars, match[1])
		}
	}
	return vars
}

// ContainsVariable checks if the content contains a specific variable
func ContainsVariable(content, variable string) bool {
	placeholder := "${" + variable + "}"
	return strings.Contains(content, placeholder)
}

// ContainsAnyVariable checks if the content contains any variables
func ContainsAnyVariable(content string) bool {
	re := regexp.MustCompile(`\$\{[A-Z_]+\}`)
	return re.MatchString(content)
}

// ValidateVariables checks if all variables in content are known
func ValidateVariables(content string) []string {
	knownVars := map[string]bool{
		"INPUT":     true,
		"DATE":      true,
		"DATETIME":  true,
		"CLIPBOARD": true,
		"FILE":      true,
	}

	vars := ExtractVariables(content)
	var unknown []string
	for _, v := range vars {
		if !knownVars[v] {
			unknown = append(unknown, v)
		}
	}
	return unknown
}
