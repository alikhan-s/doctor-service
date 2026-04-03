package http

import (
	"net/http"

	"github.com/alikhan-s/doctor-s/internal/model"
	"github.com/alikhan-s/doctor-s/internal/usecase"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	usecase usecase.DoctorUseCase
}

func NewDoctorHandler(r *gin.Engine, u usecase.DoctorUseCase) {
	handler := &DoctorHandler{usecase: u}

	r.POST("/doctors", handler.Create)
	r.GET("/doctors/:id", handler.GetByID)
	r.GET("/doctors", handler.GetAll)
}

func (h *DoctorHandler) Create(c *gin.Context) {
	var doctor model.Doctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.usecase.CreateDoctor(c.Request.Context(), &doctor); err != nil {
		if err == model.ErrInvalidFullName || err == model.ErrInvalidEmail || err == model.ErrInvalidEmailFormat || err == model.ErrEmailExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create doctor in database"})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	doctor, err := h.usecase.GetDoctorByID(c.Request.Context(), id)
	if err != nil {
		if err == model.ErrDoctorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "there is no doctor like this"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) GetAll(c *gin.Context) {
	doctors, err := h.usecase.GetAllDoctors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}

	if doctors == nil {
		doctors = []*model.Doctor{}
	}

	c.JSON(http.StatusOK, doctors)
}
