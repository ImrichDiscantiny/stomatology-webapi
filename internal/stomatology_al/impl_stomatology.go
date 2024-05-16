package stomatology_al

import (
	"net/http"

	"github.com/ImrichDiscantiny/stomatology-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Kópia zakomentovanej časti z api_stomatology.go
// CreateStomatology - Saves new stomatology definition
func (this *implStomatologyAPI) CreateStomatology(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
    if !exists {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "db not found",
                "error":   "db not found",
            })
        return
    }

    db, ok := value.(db_service.DbService[Stomatology])
    if !ok {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "db context is not of required type",
                "error":   "cannot cast db context to db_service.DbService",
            })
        return
    }

    stomatology := Stomatology{}
    err := ctx.BindJSON(&stomatology)
    if err != nil {
        ctx.JSON(
            http.StatusBadRequest,
            gin.H{
                "status":  "Bad Request",
                "message": "Invalid request body",
                "error":   err.Error(),
            })
        return
    }

    if stomatology.Id == "" {
        stomatology.Id = uuid.New().String()
    }
    // fmt.Println("Deleting stomatology with ID:", stomatology.DescriptionAppointment)
    err = db.CreateDocument(ctx, stomatology.Id, &stomatology)

    switch err {
    case nil:
        ctx.JSON(
            http.StatusCreated,
            stomatology,
        )
    case db_service.ErrConflict:
        ctx.JSON(
            http.StatusConflict,
            gin.H{
                "status":  "Conflict",
                "message": "Stomatology already exists",
                "error":   err.Error(),
            },
        )
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to create stomatology in database",
                "error":   err.Error(),
            },
        )
    }
}

// DeleteStomatology - Deletes specific stomatology
func (this *implStomatologyAPI) DeleteStomatology(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
    if !exists {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "db_service not found",
                "error":   "db_service not found",
            })
        return
    }

    db, ok := value.(db_service.DbService[Stomatology])
    if !ok {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "db_service context is not of type db_service.DbService",
                "error":   "cannot cast db_service context to db_service.DbService",
            })
        return
    }

    stomatologyId := ctx.Param("stomatologyId")

    err := db.DeleteDocument(ctx, stomatologyId)

    switch err {
    case nil:
        ctx.AbortWithStatus(http.StatusNoContent)
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Stomatology not found",
                "error":   err.Error(),
            },
        )
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to delete stomatology from database",
                "error":   err.Error(),
            })
    }
}