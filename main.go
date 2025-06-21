package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Student struct {
	ID        int
	Name      string
	ClassRoom string
	Age       int
}

func main() {
	dsn := "truongbvn:Truong@123@tcp(14.225.217.120:3306)/student_db?parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Student{})

	r := gin.Default()

	r.GET("/students", func(c *gin.Context) {
		var students []Student

		db.Find(&students)
		c.JSON(http.StatusOK, gin.H{
			"students": students,
		})
	})

	r.GET("/student-detail/:id", func(c *gin.Context) {
		var student Student
		id := c.Param("id")

		if err := db.First(&student, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Student not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	})

	r.POST("/student", func(c *gin.Context) {
		var student Student

		err := c.BindJSON(&student)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		db.Create(&student)
		c.JSON(http.StatusOK, gin.H{
			"message": "Add student successfully",
			"student": student,
		})
	})

	r.PUT("/student/:id", func(c *gin.Context) {
		var student Student
		id := c.Param("id")

		if err := db.First(&student, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Student not found",
			})
			return
		}

		var updateData Student
		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		student.Name = updateData.Name
		db.Save(&student)

		c.JSON(http.StatusOK, gin.H{
			"message": "Update student successfully",
			"student": student,
		})
	})

	r.DELETE("/student/:id", func(c *gin.Context) {
		var student Student
		id := c.Param("id")

		if err := db.First(&student, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Student not found",
			})
			return
		}

		db.Delete(&student)
		c.JSON(http.StatusOK, gin.H{
			"message": "Delete student successfully",
		})
	})

	r.Run()
}
