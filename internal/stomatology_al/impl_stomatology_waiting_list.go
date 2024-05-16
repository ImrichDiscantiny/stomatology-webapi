package stomatology_al

import (
	"net/http"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_ambulance_Appointment_list.go

// CreateAppointmentListEntry - Saves new entry into Appointment list
func (this *implStomatologyAppointmentListAPI) CreateAppointmentListEntry(ctx *gin.Context) {
  updateStomatologyFunc(ctx, func(c *gin.Context, stomatology *Stomatology) (*Stomatology,  interface{},  int){
    var entry AppointmentListEntry

    if err := c.ShouldBindJSON(&entry); err != nil {
        return nil, gin.H{
            "status": http.StatusBadRequest,
            "message": "Invalid request body",
            "error": err.Error(),
        }, http.StatusBadRequest
    }

    if entry.Id == "" {
        return nil, gin.H{
            "status": http.StatusBadRequest,
            "message": "Patient ID is required",
        }, http.StatusBadRequest
    }

    if entry.Id == "" || entry.Id == "@new" {
        entry.Id = uuid.NewString()
    }

    conflictIndx := slices.IndexFunc( stomatology.AppointmentList, func(appointment AppointmentListEntry) bool {
        return entry.Id == appointment.Id || entry.Id == appointment.Id
    })

    if conflictIndx >= 0 {
        return nil, gin.H{
            "status": http.StatusConflict,
            "message": "Entry already exists",
        }, http.StatusConflict
    }

    stomatology.AppointmentList = append(stomatology.AppointmentList, entry)
    // stomatology.reconcileAppointmentList()

    // entry was copied by value return reconciled value from the list
    entryIndx := slices.IndexFunc( stomatology.AppointmentList, func(appointment AppointmentListEntry) bool {
        return entry.Id == appointment.Id
    })
    if entryIndx < 0 {
        return nil, gin.H{
            "status": http.StatusInternalServerError,
            "message": "Failed to save entry",
        }, http.StatusInternalServerError
    }
    return stomatology, stomatology.AppointmentList[entryIndx], http.StatusOK
})
}

// DeleteAppointmentListEntry - Deletes specific entry
func (this *implStomatologyAppointmentListAPI) DeleteAppointmentListEntry(ctx *gin.Context) {
      updateStomatologyFunc(ctx, func(c *gin.Context, stomatology *Stomatology) (*Stomatology, interface{}, int) {
        entryId := ctx.Param("entryId")

        if entryId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Entry ID is required",
            }, http.StatusBadRequest
        }

        entryIndx := slices.IndexFunc(stomatology.AppointmentList, func(appointment AppointmentListEntry) bool {
            return entryId == appointment.Id
        })

        if entryIndx < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Entry not found",
            }, http.StatusNotFound
        }

        stomatology.AppointmentList = append(stomatology.AppointmentList[:entryIndx], stomatology.AppointmentList[entryIndx+1:]...)
        // stomatology.reconcileWaitingList()
        return stomatology, nil, http.StatusNoContent
    })
}

// GetAppointmentListEntries - Provides the ambulance Appointment list
func (this *implStomatologyAppointmentListAPI) GetAppointmentListEntries(ctx *gin.Context) {
  updateStomatologyFunc(ctx, func(c *gin.Context, stomatology *Stomatology) (*Stomatology, interface{}, int) {
    result := stomatology.AppointmentList
    if result == nil {
        result = []AppointmentListEntry{}
    }
    // return nil ambulance - no need to update it in db
    return nil, result, http.StatusOK
})
}

// UpdateAppointmentListEntry - Updates specific entry
func (this *implStomatologyAppointmentListAPI) UpdateAppointmentListEntry(ctx *gin.Context) {
  updateStomatologyFunc(ctx, func(c *gin.Context, stomatology *Stomatology) (*Stomatology, interface{}, int) {
    var entry AppointmentListEntry

    if err := c.ShouldBindJSON(&entry); err != nil {
        return nil, gin.H{
            "status":  http.StatusBadRequest,
            "message": "Invalid request body",
            "error":   err.Error(),
        }, http.StatusBadRequest
    }

    entryId := ctx.Param("entryId")

    if entryId == "" {
        return nil, gin.H{
            "status":  http.StatusBadRequest,
            "message": "Entry ID is required",
        }, http.StatusBadRequest
    }

    entryIndx := slices.IndexFunc(stomatology.AppointmentList, func(waiting AppointmentListEntry) bool {
        return entryId == waiting.Id
    })

    if entryIndx < 0 {
        return nil, gin.H{
            "status":  http.StatusNotFound,
            "message": "Entry not found",
        }, http.StatusNotFound
    }

    if entry.Id != "" {
        stomatology.AppointmentList[entryIndx].Id = entry.Id
    }

    if entry.Id != "" {
        stomatology.AppointmentList[entryIndx].Id = entry.Id
    }

    // if entry.WaitingSince.After(time.Time{}) {
    //     ambulance.WaitingList[entryIndx].WaitingSince = entry.WaitingSince
    // }

    // if entry.EstimatedDurationMinutes > 0 {
    //     ambulance.WaitingList[entryIndx].EstimatedDurationMinutes = entry.EstimatedDurationMinutes
    // }

    // ambulance.reconcileWaitingList()
    return stomatology, stomatology.AppointmentList[entryIndx], http.StatusOK
})
}