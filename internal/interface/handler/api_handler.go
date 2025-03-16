package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
	"github.com/masahiro0000/BizMatch/internal/usecase"
)

type ApiHandler struct {
	apiUsecase *usecase.ApiUsecase
}

func NewApiHandler(u *usecase.ApiUsecase) *ApiHandler{
	return &ApiHandler{
		apiUsecase: u,
	}
}

// GetPrefectures handles the HTTP request to retrieve prefecture data.
func (h *ApiHandler) GetPrefectures(c *gin.Context) {
	// Call the use case to get prefecture data.
	prefectures, err := h.apiUsecase.GetPrefectures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// Map the retrieved data to the output structure.
	var prefecture []domain.Prefectures
	for _, p := range prefectures {
		prefecture = append(prefecture, domain.Prefectures{
			ID:		p.ID,
			Name:	p.Name,
		})
	}

	c.JSON(http.StatusOK, prefecture)
}

// GetIndustries handles the HTTP request to retrieve industry data.
func (h *ApiHandler) GetIndustries(c *gin.Context) {
	// Call the use case to get industry data.
	industries, err := h.apiUsecase.GetIndustries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// Map the retrieved data to the output structure.
	var industry []domain.Industries
	for _, p := range industries {
		industry = append(industry, domain.Industries{
			ID:		p.ID,
			Name:	p.Name,
		})
	}
	c.JSON(http.StatusOK, industry)
}

// GetJobs handles the HTTP request to retrieve job data.
func (h *ApiHandler) GetJobs(c *gin.Context) {
	// Call the use case to get job data.
	jobs, err := h.apiUsecase.GetJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// Map the retrieved data to the output structure.
	var job []domain.Jobs
	for _, p := range jobs {
		job = append(job, domain.Jobs{
			ID:		p.ID,
			Name:	p.Name,
		})
	}
	c.JSON(http.StatusOK, job)
}

// GetPositions handles the HTTP request to retrieve position data.
func (h *ApiHandler) GetPositions(c *gin.Context) {
	// Call the use case to get position data.
	positions, err := h.apiUsecase.GetPositions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// Map the retrieved data to the output structure.
	var position []domain.Positions
	for _, p := range positions {
		position = append(position, domain.Positions{
			ID:		p.ID,
			Name:	p.Name,
		})
	}
	c.JSON(http.StatusOK, position)
}
