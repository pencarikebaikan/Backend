package main

import (
    "net/http"
    "strconv"
    "strings"
	"github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// Model sesuai tabel posts
type Post struct {
    ID          uint   `gorm:"primaryKey"`
    Title       string `gorm:"size:200;not null"`
    Content     string `gorm:"type:text;not null"`
    Category    string `gorm:"size:100;not null"`
    Status      string `gorm:"size:100;not null"`
    CreatedDate string
    UpdatedDate string
}

func main() {
    // Koneksi ke MySQL (pakai user & password kamu)
    dsn := "alwi:alwi2007@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Gagal koneksi database")
    }

    // Auto migrate tabel posts
    db.AutoMigrate(&Post{})

    r := gin.Default()

	r.Use(cors.Default()) // 🔥 WAJIB biar React bisa akses

    // CREATE article
    r.POST("/article", func(c *gin.Context) {
        var post Post
        if err := c.ShouldBindJSON(&post); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        // Validasi
        if len(post.Title) < 20 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Title minimal 20 karakter"})
            return
        }
        if len(post.Content) < 200 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Content minimal 200 karakter"})
            return
        }
        if len(post.Category) < 3 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Category minimal 3 karakter"})
            return
        }
        status := strings.ToLower(post.Status)
     //   if status != "publish" && status != "draft" && status != "thrash" {
     //       c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus publish/draft/thrash"})
     //       return
      //  }
		if status != "publish" && status != "draft" && status != "trash" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus publish/draft/trash"})
			return
		}
		
        db.Create(&post)
        c.JSON(http.StatusOK, gin.H{"message": "Article created"})
    })

    // READ all with paging
    r.GET("/article/:limit/:offset", func(c *gin.Context) {
        var posts []Post

        limitStr := c.Param("limit")
        offsetStr := c.Param("offset")

        limit, err1 := strconv.Atoi(limitStr)
        offset, err2 := strconv.Atoi(offsetStr)
        if err1 != nil || err2 != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Limit dan offset harus angka"})
            return
        }

        db.Limit(limit).Offset(offset).Find(&posts)
        c.JSON(http.StatusOK, posts)
    })

    // READ by ID
    r.GET("/article/id/:id", func(c *gin.Context) {
        var post Post
        id := c.Param("id")
        if err := db.First(&post, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
            return
        }
        c.JSON(http.StatusOK, post)
    })

    // UPDATE by ID
    r.PUT("/article/id/:id", func(c *gin.Context) {
        var post Post
        id := c.Param("id")
        if err := db.First(&post, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
            return
        }

        var input Post
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Validasi sama seperti create
        if len(input.Title) < 20 || len(input.Content) < 200 || len(input.Category) < 3 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Validasi gagal"})
            return
        }
        status := strings.ToLower(input.Status)
       // if status != "publish" && status != "draft" && status != "thrash" {
        //    c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus publish/draft/thrash"})
        //    return
       // }
		if status != "publish" && status != "draft" && status != "trash" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus publish/draft/trash"})
			return
		}		
		
        db.Model(&post).Updates(input)
        c.JSON(http.StatusOK, gin.H{"message": "Article updated"})
    })

    // DELETE by ID
    r.DELETE("/article/id/:id", func(c *gin.Context) {
        id := c.Param("id")
        if err := db.Delete(&Post{}, id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
    })

    // Jalankan server di port 3000
    r.Run(":3000")
}
