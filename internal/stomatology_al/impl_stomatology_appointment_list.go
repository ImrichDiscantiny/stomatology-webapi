package stomatology_al

import (
	"fmt"
	"net/http"

	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FormatDateString(y int, m time.Month, d int) string {
	return fmt.Sprintf("%04d-%02d-%02d", y, int(m), d)
}

func isNotOldDate(appDate string) bool {
	now := time.Now()

	y, m, d := now.Date()

	dateStr := FormatDateString(y, m, d)

	CurrDate, err1 := time.Parse("2006-01-02", dateStr)
	nextDate, err2 := time.Parse("2006-01-02", appDate)

	if err1 != nil {
		// fmt.Println(err1)
		return true
	}

	if err2 != nil {
		// fmt.Println("error2")
		return true
	}

	if nextDate.Before(CurrDate) {
		// fmt.Println(appDate, "is before", dateStr)
		return true
	} else {
		// fmt.Println(dateStr, "is before", appDate)
		return false
	}

}

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_ambulance_Appointment_list.go

// CreateAppointmentListEntry - Saves new entry into Appointment list
func (this *implStomatologyAppointmentListAPI) CreateAppointmentListEntry(ctx *gin.Context) {
	updateStomatologyFunc(ctx, func(c *gin.Context, stomatology *Stomatology) (*Stomatology, interface{}, int) {
		var entry AppointmentListEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.Id == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Patient ID is required",
			}, http.StatusBadRequest
		}

		isNotOld := isNotOldDate(entry.Date)

		if isNotOld {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Cannot add old appointments",
			}, http.StatusBadRequest
		}

		if entry.Id == "" || entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(stomatology.AppointmentList, func(appointment AppointmentListEntry) bool {
			return entry.Id == appointment.Id || (entry.Date == appointment.Date && entry.Duration == appointment.Duration)
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		stomatology.AppointmentList = append(stomatology.AppointmentList, entry)
		// stomatology.reconcileAppointmentList()

		// entry was copied by value return reconciled value from the list
		entryIndx := slices.IndexFunc(stomatology.AppointmentList, func(appointment AppointmentListEntry) bool {
			return entry.Id == appointment.Id
		})
		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}
		return stomatology, stomatology.AppointmentList[entryIndx], http.StatusOK
	})
}

// DeleteAppointmentListEntry - Deletes specific entry
func (this *implStomatologyAppointmentListAPI) DeleteAppointmentListEntry(ctx *gin.Context) {
	updateStomatologyFunc(ctx, func(c *gin.Context, stomatology *Stomatology) (*Stomatology, interface{}, int) {
		entryId := ctx.Param("id")

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

		isNotOld := isNotOldDate(stomatology.AppointmentList[entryIndx].Date)

		if isNotOld {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Cannot delete old appointments",
			}, http.StatusBadRequest
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

		pivotDate := ctx.Param("appointmentsDate")

		date, err := time.Parse("2006-01-02", pivotDate)

		if err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Bad date",
			}, http.StatusBadRequest
		}

		var offsetStart int
		var offsetEnd int

		// location, err := time.LoadLocation("Europe/Bratislava")
		// if err != nil {
		//     return nil, gin.H{
		//         "status":  http.StatusBadRequest,
		//         "message": "Bad Slovak date",
		//     }, http.StatusBadRequest
		// }

		weekday := date.Weekday()

		switch weekday {
		case time.Monday:
			offsetStart = 0
			offsetEnd = 4
		case time.Tuesday:
			offsetStart = 1
			offsetEnd = 3
		case time.Wednesday:
			offsetStart = 2
			offsetEnd = 2
		case time.Thursday:
			offsetStart = 3
			offsetEnd = 1
		case time.Friday:
			offsetStart = 4
			offsetEnd = 0
		case time.Saturday:
			offsetStart = 5
			offsetEnd = 0
		case time.Sunday:
			offsetStart = 6
			offsetEnd = 0
		}

		startDate := date.AddDate(0, 0, -offsetStart)
		endDate := date.AddDate(0, 0, offsetEnd)

		fmt.Println("start: ", startDate, "endDate: ", endDate)
		var filtered []AppointmentListEntry

		for _, appointment := range result {
			currDate, err := time.Parse("2006-01-02", appointment.Date)

			if err == nil && (currDate.After(startDate) || currDate.Equal(startDate)) &&
				(currDate.Before(endDate) || currDate.Equal(endDate)) {
				filtered = append(filtered, appointment)
			}
		}

		return nil, filtered, http.StatusOK
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

		entryId := ctx.Param("id")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Entry ID is required",
			}, http.StatusBadRequest
		}

		isNotOld := isNotOldDate(entry.Date)

		if isNotOld {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Cannot add old appointments",
			}, http.StatusBadRequest
		}

		entryIndxIn := slices.IndexFunc(stomatology.AppointmentList, func(appointment AppointmentListEntry) bool {
			return entryId == appointment.Id
		})

		if entryIndxIn < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		isNotOld = isNotOldDate(stomatology.AppointmentList[entryIndxIn].Date)

		if isNotOld {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Cannot update old appointments",
			}, http.StatusBadRequest
		}

		DateMatchIndx := slices.IndexFunc(stomatology.AppointmentList, func(appointment AppointmentListEntry) bool {
			return entryId != appointment.Id && entry.Date == appointment.Date && entry.Duration == appointment.Duration
		})

		if DateMatchIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Time of Entry already exists, pick another one",
			}, http.StatusNotFound
		}

		stomatology.AppointmentList[entryIndxIn] = entry
		return stomatology, stomatology.AppointmentList[entryIndxIn], http.StatusOK
	})
}
