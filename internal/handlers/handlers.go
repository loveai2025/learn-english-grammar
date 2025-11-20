package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/user/egm/internal/database"
	"github.com/user/egm/internal/models"
)

func ListCourses(c *gin.Context) {
	level := c.Param("level")
	var courses []models.Course

	if err := database.DB.Where("level = ?", level).Order("order_index asc").Find(&courses).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Message": "Error fetching courses"})
		return
	}

	c.HTML(http.StatusOK, "course_list.html", gin.H{
		"Title":   "Courses - " + level,
		"Level":   level,
		"Courses": courses,
	})
}

func CourseDetail(c *gin.Context) {
	// level := c.Param("level")
	slug := c.Param("slug")
	var course models.Course

	if err := database.DB.Where("slug = ?", slug).First(&course).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"Message": "Course not found"})
		return
	}

	// TODO: Render Markdown content to HTML
	// For now, we pass raw markdown.
	// We should probably use a library like "github.com/russross/blackfriday" or similar, or do it in JS.
	// Let's assume we pass it as is and the frontend might handle it or we use a go library.
	// The requirement says "Render Markdown from DB".

	c.HTML(http.StatusOK, "course_detail.html", gin.H{
		"Title":  course.Title,
		"Course": course,
	})
}

func QuizPage(c *gin.Context) {
	// level := c.Param("level")
	slug := c.Param("slug")
	var course models.Course

	if err := database.DB.Where("slug = ?", slug).First(&course).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"Message": "Course not found"})
		return
	}

	var exercises []models.Exercise
	if err := database.DB.Where("course_id = ?", course.ID).Find(&exercises).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Message": "Error fetching exercises"})
		return
	}

	c.HTML(http.StatusOK, "quiz.html", gin.H{
		"Title":     "Quiz - " + course.Title,
		"Course":    course,
		"Exercises": exercises,
	})
}

type QuizSubmission struct {
	CourseID uint              `json:"course_id"`
	Answers  map[uint]string `json:"answers"` // ExerciseID -> Answer
}

type QuizResult struct {
	Score       int                    `json:"score"`
	Total       int                    `json:"total"`
	Passed      bool                   `json:"passed"`
	Explanation map[uint]string        `json:"explanation"` // ExerciseID -> Explanation
	Results     map[uint]bool          `json:"results"`     // ExerciseID -> Correct/Incorrect
}

func SubmitQuiz(c *gin.Context) {
	var submission QuizSubmission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exercises []models.Exercise
	if err := database.DB.Where("course_id = ?", submission.CourseID).Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching exercises"})
		return
	}

	score := 0
	total := len(exercises)
	results := make(map[uint]bool)
	explanations := make(map[uint]string)

	for _, ex := range exercises {
		userAns := submission.Answers[ex.ID]
		isCorrect := userAns == ex.CorrectAnswer
		results[ex.ID] = isCorrect
		explanations[ex.ID] = ex.Explanation
		if isCorrect {
			score++
		}
	}

	passed := float64(score)/float64(total) >= 0.6 // 60% pass rate

	// Save progress
	userID, exists := c.Get("user_id")
	if exists {
		progress := models.UserProgress{
			UserID:      userID.(string),
			CourseID:    submission.CourseID,
			IsCompleted: passed,
			Score:       score,
		}
		// Verify if record exists to update or create
		var existingProgress models.UserProgress
		if err := database.DB.Where("user_id = ? AND course_id = ?", userID.(string), submission.CourseID).First(&existingProgress).Error; err == nil {
			// Update
			existingProgress.IsCompleted = passed || existingProgress.IsCompleted // Once completed, stays completed
			existingProgress.Score = score // Update score to latest
			database.DB.Save(&existingProgress)
		} else {
			// Create
			database.DB.Create(&progress)
		}
	}

	c.JSON(http.StatusOK, QuizResult{
		Score:       score,
		Total:       total,
		Passed:      passed,
		Explanation: explanations,
		Results:     results,
	})
}
