package pharmacy

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xmikova/ambulance-webapi/internal/db_service"
)

func getPharmacyFunc(ctx *gin.Context, fn func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int)) {
	value, exists := ctx.Get("pharmacy_db_service")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "pharmacy db service not found"})
		return
	}
	db, ok := value.(db_service.DbService[PharmacyStore])
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "pharmacy db service type error"})
		return
	}

	pharmacyId := ctx.Param("pharmacyId")
	store, err := db.FindDocument(ctx, pharmacyId)
	switch err {
	case nil:
	case db_service.ErrNotFound:
		// auto-create on first access (useful for local dev without init container)
		store = &PharmacyStore{Id: pharmacyId, Medicines: []Medicine{}}
		if createErr := db.CreateDocument(ctx, pharmacyId, store); createErr != nil && createErr != db_service.ErrConflict {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": createErr.Error()})
			return
		}
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	updatedStore, responseObject, statusCode := fn(ctx, store)
	if updatedStore != nil {
		if err := db.UpdateDocument(ctx, pharmacyId, updatedStore); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	if responseObject != nil {
		ctx.JSON(statusCode, responseObject)
	} else {
		ctx.AbortWithStatus(statusCode)
	}
}

func GetMedicines(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		medicines := store.Medicines
		if medicines == nil {
			medicines = []Medicine{}
		}
		return nil, medicines, http.StatusOK
	})
}

func GetMedicine(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		medicineId := c.Param("medicineId")
		idx := slices.IndexFunc(store.Medicines, func(m Medicine) bool { return m.Id == medicineId })
		if idx < 0 {
			return nil, gin.H{"message": "Medicine not found"}, http.StatusNotFound
		}
		return nil, store.Medicines[idx], http.StatusOK
	})
}

func CreateMedicine(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		var medicine Medicine
		if err := c.ShouldBindJSON(&medicine); err != nil {
			return nil, gin.H{"message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		if medicine.Name == "" {
			return nil, gin.H{"message": "Name is required"}, http.StatusBadRequest
		}
		if medicine.Id == "" || medicine.Id == "@new" {
			medicine.Id = uuid.NewString()
		}
		conflict := slices.IndexFunc(store.Medicines, func(m Medicine) bool { return m.Id == medicine.Id })
		if conflict >= 0 {
			return nil, gin.H{"message": "Medicine already exists"}, http.StatusConflict
		}
		store.Medicines = append(store.Medicines, medicine)
		return store, medicine, http.StatusOK
	})
}

func UpdateMedicine(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		medicineId := c.Param("medicineId")
		var medicine Medicine
		if err := c.ShouldBindJSON(&medicine); err != nil {
			return nil, gin.H{"message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		idx := slices.IndexFunc(store.Medicines, func(m Medicine) bool { return m.Id == medicineId })
		if idx < 0 {
			return nil, gin.H{"message": "Medicine not found"}, http.StatusNotFound
		}
		medicine.Id = medicineId
		store.Medicines[idx] = medicine
		return store, medicine, http.StatusOK
	})
}

func DeleteMedicine(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		medicineId := c.Param("medicineId")
		idx := slices.IndexFunc(store.Medicines, func(m Medicine) bool { return m.Id == medicineId })
		if idx < 0 {
			return nil, gin.H{"message": "Medicine not found"}, http.StatusNotFound
		}
		store.Medicines = append(store.Medicines[:idx], store.Medicines[idx+1:]...)
		return store, nil, http.StatusNoContent
	})
}
