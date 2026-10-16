package auth

import (
	"strings"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Rank         string    `json:"rank"`
	RegisteredAt time.Time `json:"registeredAt"`
	RegisteredBy string    `json:"registeredBy"`
}

type Rank struct {
	Name        string   `json:"name"`
	Level       int      `json:"level"`
	Commands    []string `json:"commands"`
	Description string   `json:"description"`
}

type Group struct {
	ID           string    `json:"id"`
	RegisteredAt time.Time `json:"registeredAt"`
	RegisteredBy string    `json:"registeredBy"`
}

func (r *Rank) HasCommand(command string) bool {
	for _, cmd := range r.Commands {
		if cmd == "*" || cmd == command {
			return true
		}
	}

	return false
}

func ParseCommands(raw string) []string {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")

	commands := make([]string, 0, len(parts))
	for _, part := range parts {
		if cmd := strings.TrimSpace(part); cmd != "" {
			commands = append(commands, cmd)
		}
	}

	return commands
}

func JoinCommands(commands []string) string {
	if len(commands) == 0 {
		return ""
	}

	filtered := make([]string, 0, len(commands))
	for _, cmd := range commands {
		if cmd := strings.TrimSpace(cmd); cmd != "" {
			filtered = append(filtered, cmd)
		}
	}

	return strings.Join(filtered, ",")
}
