package logparser

import "regexp"

// Regular expressions to identify event types
var RegexPatterns = map[string]*regexp.Regexp{
	// ^ matches the start of the line
	// .* matches any character zero or more times
	// InitGame: matches the string "InitGame:"
	// $ matches the end of the line
	"InitGame": regexp.MustCompile(`^.*InitGame: (.*)$`),

	// ^ matches the start of the line
	// .* matches any character zero or more times
	// Kill: matches the string "Kill:"
	// (\d+) (\d+) (\d+): matches one or more digits (3 times) followed by a colon
	// (.+) killed (.+) by (.+) matches any character one or more times followed by "killed" and any character one or more times followed by "by" and any character one or more times
	// $ matches the end of the line
	"Kill": regexp.MustCompile(`^.*Kill: (\d+) (\d+) (\d+): (.+) killed (.+) by (.+)$`),

	// ^ matches the start of the line
	// .* matches any character zero or more times
	// ClientUserinfoChanged: matches the string "ClientUserinfoChanged:"
	// (\d+) matches one or more digits
	// n\\([^\\]+)\\t = n\\ matches "n\", ([^\\]+) captures a group of characters (username) excluding the backslash, \\t ends with "\t" after the username
	"ClientUserinfoChanged": regexp.MustCompile(`^.*ClientUserinfoChanged: (\d+) n\\([^\\]+)\\t`),

	// ^ matches the start of the line
	// .* matches any character zero or more times
	// ShutdownGame: matches the string "ShutdownGame:"
	// $ matches the end of the line
	"ShutdownGame": regexp.MustCompile(`^.*ShutdownGame:(.*)$`),
}
