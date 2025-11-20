package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"github.com/joho/godotenv"
	"github.com/user/egm/internal/database"
	"github.com/user/egm/internal/handlers"
	"github.com/user/egm/internal/middleware"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Initialize Database
	database.Connect()

	// Initialize Gin
	r := gin.Default()

	// Add custom template functions
	r.SetFuncMap(map[string]any{
		"title": func(s string) string {
			return cases.Title(language.English).String(s)
		},
	})

	// Load HTML Templates
	r.LoadHTMLGlob("views/**/*")

	// Static Files
	r.Static("/static", "./static")

	// Routes
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":       "English Grammar Master",
			"SupabaseURL": supabaseURL,
			"SupabaseKey": supabaseKey,
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Title":       "Login",
			"SupabaseURL": supabaseURL,
			"SupabaseKey": supabaseKey,
		})
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{
			"Title":       "Signup",
			"SupabaseURL": supabaseURL,
			"SupabaseKey": supabaseKey,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Course Routes (Protected)
	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		authorized.GET("/learn/:level", handlers.ListCourses)
		authorized.GET("/learn/:level/:slug", handlers.CourseDetail)
		authorized.GET("/learn/:level/:slug/quiz", handlers.QuizPage)
		authorized.POST("/api/quiz/submit", handlers.SubmitQuiz)
	}

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
