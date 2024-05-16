package stomatology_al

import (
	"net/http"

	"github.com/ImrichDiscantiny/stomatology-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
)

type stomatologyUpdater = func(
    ctx *gin.Context,
    stomatology *Stomatology,
) (updateStomatology *Stomatology, responseContent interface{}, status int)

func updateStomatologyFunc(ctx *gin.Context, updater stomatologyUpdater) {
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

    stomatology, err := db.FindDocument(ctx, stomatologyId)

    switch err {
    case nil:
        // continue
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Stoamtology not found",
                "error":   err.Error(),
            },
        )
        return
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to load stoamtology from database",
                "error":   err.Error(),
            })
        return
    }

    if !ok {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "Failed to cast stoamtology from database",
                "error":   "Failed to cast stoamtology from database",
            })
        return
    }

    updatedStomatology, responseObject, status := updater(ctx, stomatology)

    if updatedStomatology != nil {
        err = db.UpdateDocument(ctx, stomatologyId, updatedStomatology)
    } else {
        err = nil // redundant but for clarity
    }

    switch err {
    case nil:
        if responseObject != nil {
            ctx.JSON(status, responseObject)
        } else {
            ctx.AbortWithStatus(status)
        }
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Ambulance was deleted while processing the request",
                "error":   err.Error(),
            },
        )
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to update ambulance in database",
                "error":   err.Error(),
            })
    }

}