package models

import "time"

type Log struct {
	ID        string    `json:"id" bson:"_id"`
	Content   string    `json:"content,omitempty" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}

type AnalysisResult struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Version     string       `json:"version"`
	Information []Info       `json:"information"`
	Problems    []Problem    `json:"problems"`
}

type Info struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Problem struct {
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	Solutions []Solution `json:"solutions"`
}

type Solution struct {
	Message string `json:"message"`
}
