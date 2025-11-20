package models

import (
	"time"
)

type DifficultyLevel string

const (
	DifficultyBeginner     DifficultyLevel = "beginner"
	DifficultyIntermediate DifficultyLevel = "intermediate"
	DifficultyAdvanced     DifficultyLevel = "advanced"
)

type QuestionType string

const (
	QuestionTypeChoice    QuestionType = "choice"
	QuestionTypeFillBlank QuestionType = "fill_blank"
)

type Profile struct {
	ID           string          `gorm:"primaryKey;type:uuid" json:"id"`
	Username     string          `json:"username"`
	AvatarURL    string          `json:"avatar_url"`
	CurrentLevel DifficultyLevel `gorm:"type:difficulty_level;default:beginner" json:"current_level"`
	CreatedAt    time.Time       `json:"created_at"`
}

type Course struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	Title      string          `json:"title"`
	Level      DifficultyLevel `gorm:"type:difficulty_level" json:"level"`
	ContentMD  string          `json:"content_md"`
	OrderIndex int             `json:"order_index"`
	Slug       string          `gorm:"unique" json:"slug"`
}

type Exercise struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	CourseID      uint         `json:"course_id"`
	Course        Course       `json:"course,omitempty"`
	Question      string       `json:"question"`
	Options       string       `gorm:"type:jsonb" json:"options"` // Storing JSON as string or []byte for GORM, or use specific JSON type
	CorrectAnswer string       `json:"correct_answer"`
	Explanation   string       `json:"explanation"`
	Type          QuestionType `gorm:"type:question_type" json:"type"`
}

type UserProgress struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      string    `gorm:"type:uuid" json:"user_id"`
	Profile     Profile   `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	CourseID    uint      `json:"course_id"`
	Course      Course    `json:"course,omitempty"`
	IsCompleted bool      `json:"is_completed"`
	Score       int       `json:"score"`
	CompletedAt time.Time `json:"completed_at"`
}
