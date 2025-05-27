package ambulance_virtual_patient_list

import (
	"net/http"

	"github.com/Robert-Prikryl/ambulance-virtual-patient-api/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implVirtualPatientListAPI struct {
}

func NewPatientApi() VirtualPatientListAPI {
	return &implVirtualPatientListAPI{}
}

func (o implVirtualPatientListAPI) CreateVirtualPatient(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[VirtualPatient])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	virtualPatient := VirtualPatient{}
	err := c.BindJSON(&virtualPatient)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	if virtualPatient.Id == "" {
		virtualPatient.Id = uuid.New().String()
	}

	err = db.CreateDocument(c, virtualPatient.Id, &virtualPatient)

	switch err {
	case nil:
		c.JSON(
			http.StatusCreated,
			virtualPatient,
		)
	case db_service.ErrConflict:
		c.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Virtual patient already exists",
				"error":   err.Error(),
			},
		)
	default:
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create virtual patient in database",
				"error":   err.Error(),
			},
		)
	}
}

func (o implVirtualPatientListAPI) GetVirtualPatientList(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "db_service not found"},
		)
		return
	}

	db, ok := value.(db_service.DbService[VirtualPatient])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "db_service context is not of type db_service.DbService[VirtualPatient]"},
		)
		return
	}

	patients, err := db.ListDocuments(c.Request.Context())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, patients)
}

func (o implVirtualPatientListAPI) DeleteVirtualPatient(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[VirtualPatient])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	virtualPatientId := c.Param("virtualPatientId")
	err := db.DeleteDocument(c, virtualPatientId)

	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Virtual patient not found",
				"error":   err.Error(),
			},
		)
	default:
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete virtual patient from database",
				"error":   err.Error(),
			})
	}
}

func (o implVirtualPatientListAPI) UpdateVirtualPatient(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[VirtualPatient])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	virtualPatientId := c.Param("virtualPatientId")
	virtualPatient := VirtualPatient{}

	if err := c.BindJSON(&virtualPatient); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	// Ensure the ID in the path matches the ID in the body
	virtualPatient.Id = virtualPatientId

	err := db.UpdateDocument(c, virtualPatientId, &virtualPatient)
	switch err {
	case nil:
		c.JSON(http.StatusOK, virtualPatient)
	case db_service.ErrNotFound:
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Virtual patient not found",
				"error":   err.Error(),
			})
	default:
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update virtual patient in database",
				"error":   err.Error(),
			})
	}
}
