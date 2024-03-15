package main

import (
	"net/http"

	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "food-orders/docs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Calls struct {
	gorm.Model    `json:"-"`
	Id            string `gorm:"primaryKey"  json:"id"`
	NumberOfCalls uint64 `json:"numberOfCalls"`
}

func isTrue(value string) bool {
	return value == "1" || strings.EqualFold(value, "yes") || strings.EqualFold(value, "true")
}

func connectDB() *gorm.DB {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB_NAME")
	port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

// getNumCalls godoc
// @Summary Get number of calls
// @Scheme http
// @Description Get number of calls
// @Tags calls
// @Produce json
// @Success 200 {object} []Calls
// @Router /getNumCalls [get]
func getNumCalls(context *gin.Context, db *gorm.DB) {
	var calls_arr []Calls
	db.Find(&calls_arr)
	context.JSON(http.StatusOK, calls_arr)
}

// incrementNumCalls godoc
// @Summary Increment number of calls
// @Scheme http
// @Description Increment number of calls
// @Tags calls
// @Param company query string true "Company id"
// @Produce json
// @Success 200 {object} Calls
// @Router /incrementNumCalls [post]
func incrementNumCalls(context *gin.Context, db *gorm.DB) {
	id := context.Query("company")
	calls := Calls{Id: id, NumberOfCalls: 0}
	db.FirstOrCreate(&calls, "id = ?", id)
	calls.NumberOfCalls++
	db.Save(&calls)
	context.JSON(http.StatusOK, calls)
}

// setNumCalls godoc
// @Summary Set number of calls
// @Scheme http
// @Description Set number of calls
// @Tags calls
// @Param company query string true "Company id"
// @Param numberOfCalls query uint64 true "Number of calls"
// @Produce json
// @Success 200 {object} Calls
// @Router /setNumCalls [post]
func setNumCalls(context *gin.Context, db *gorm.DB) {
	id := context.Query("company")
	numberOfCalls, err := strconv.ParseUint(context.Query("numberOfCalls"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of calls, must be a non-negative integer"})
		return
	}
	calls := Calls{Id: id, NumberOfCalls: numberOfCalls}
	db.Save(&calls)
	context.JSON(http.StatusOK, calls)
}

func main() {
	godotenv.Load()

	db := connectDB()

	db.AutoMigrate(&Calls{})

	release_mode := os.Getenv("GIN_RELEASE_MODE")
	if isTrue(release_mode) {
		gin.SetMode(gin.ReleaseMode)
	}
	// route.SetTrustedProxies([]string{"0.0.0.0"}) //to trust only a specific value
	route := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	route.GET("/getNumCalls", func(context *gin.Context) { getNumCalls(context, db) })

	route.POST("/incrementNumCalls", func(context *gin.Context) { incrementNumCalls(context, db) })

	route.POST("/setNumCalls", func(context *gin.Context) { setNumCalls(context, db) })

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}

}
