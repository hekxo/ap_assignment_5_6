package main

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Barber struct {
	ID         int
	Name       string
	BasicInfo  string
	Price      int
	Experience string
	Status     string
	ImagePath  string
}

var db *sql.DB
var tpl *template.Template

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://hekxo:123456@localhost/barbershop?sslmode=disable")
	if err != nil {
		panic(err)
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func rateLimiter(duration time.Duration) gin.HandlerFunc {
	ticker := time.NewTicker(duration)
	return func(c *gin.Context) {
		select {
		case <-ticker.C:
			c.Next()
		default:
			log.Warn("Rate limit hit")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		}
	}
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")
	router.Use(rateLimiter(time.Second))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/barbers", func(c *gin.Context) {
		barbers, err := getBarbersFromDB()
		if err != nil {
			log.WithFields(log.Fields{
				"action":    "fetch_barbers",
				"timestamp": time.Now().Format(time.RFC3339),
				"error":     err,
			}).Error("Error occurred while fetching barbers from database")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.HTML(http.StatusOK, "barbers.html", gin.H{
			"Barbers": barbers,
		})
	})

	router.GET("/filtered-barbers", func(c *gin.Context) {
		statusFilter := c.Query("status")
		experienceFilter := c.Query("experience")
		sortBy := c.Query("sort")
		pageStr := c.Query("page")
		itemsPerPage := 3

		log.WithFields(log.Fields{
			"action":           "filter_barbers",
			"timestamp":        time.Now().Format(time.RFC3339),
			"statusFilter":     statusFilter,
			"experienceFilter": experienceFilter,
			"sortBy":           sortBy,
			"page":             pageStr,
		}).Info("Filtering and sorting barbers")

		barbers, err := getFilteredBarbersFromDB(statusFilter, experienceFilter, sortBy, pageStr, itemsPerPage)
		if err != nil {
			log.WithFields(log.Fields{
				"action":    "filter_barbers",
				"timestamp": time.Now().Format(time.RFC3339),
				"error":     err,
			}).Error("Error occurred while fetching filtered barbers from database")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.HTML(http.StatusOK, "barbers.html", gin.H{
			"Barbers": barbers,
		})
	})

	router.Run(":8080")
}

func getBarbersFromDB() ([]Barber, error) {
	rows, err := db.Query("SELECT id, name, basic_info, price, experience, status, image_path FROM barbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barbers []Barber
	for rows.Next() {
		var b Barber
		if err := rows.Scan(&b.ID, &b.Name, &b.BasicInfo, &b.Price, &b.Experience, &b.Status, &b.ImagePath); err != nil {
			return nil, err
		}
		barbers = append(barbers, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return barbers, nil
}

func getFilteredBarbersFromDB(statusFilter, experienceFilter, sortBy, pageStr string, itemsPerPage int) ([]Barber, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	query := "SELECT id, name, basic_info, price, experience, status, image_path FROM barbers WHERE true"
	if statusFilter != "" {
		query += " AND status = '" + statusFilter + "'"
	}
	if experienceFilter != "" {
		query += " AND experience = '" + experienceFilter + "'"
	}
	switch sortBy {
	case "name":
		query += " ORDER BY name"
	case "price":
		query += " ORDER BY price"
	}
	query += " LIMIT " + strconv.Itoa(itemsPerPage) + " OFFSET " + strconv.Itoa((page-1)*itemsPerPage)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barbers []Barber
	for rows.Next() {
		var b Barber
		if err := rows.Scan(&b.ID, &b.Name, &b.BasicInfo, &b.Price, &b.Experience, &b.Status, &b.ImagePath); err != nil {
			return nil, err
		}
		barbers = append(barbers, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return barbers, nil
}
