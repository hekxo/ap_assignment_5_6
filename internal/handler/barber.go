package handler

import (
	"adv_programming_3_4-main/internal/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type BarberHandler struct {
	repo *repository.BarberRepository
}

func NewBarberHandler(repo *repository.BarberRepository) *BarberHandler {
	return &BarberHandler{repo: repo}
}

// GetBarbers handles the "/barbers" route
func (h *BarberHandler) GetBarbers(c *gin.Context) {
	barbers, err := h.repo.GetBarbersFromDB() // Make sure this method is defined in your repository
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
}

// GetFilteredBarbers handles the "/filtered-barbers" route
func (h *BarberHandler) GetFilteredBarbers(c *gin.Context) {
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

	barbers, err := h.repo.GetFilteredBarbersFromDB(statusFilter, experienceFilter, sortBy, pageStr, itemsPerPage) // Make sure this method is defined in your repository
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
}
