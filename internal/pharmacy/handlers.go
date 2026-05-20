package pharmacy

import (
	"net/http"
	"slices"
	"time"

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

// Orders handlers
func GetOrders(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		orders := store.Orders
		if orders == nil {
			orders = []Order{}
		}
		return nil, orders, http.StatusOK
	})
}

func GetOrder(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		orderId := c.Param("orderId")
		idx := slices.IndexFunc(store.Orders, func(o Order) bool { return o.Id == orderId })
		if idx < 0 {
			return nil, gin.H{"message": "Order not found"}, http.StatusNotFound
		}
		return nil, store.Orders[idx], http.StatusOK
	})
}

func CreateOrder(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		var order Order
		if err := c.ShouldBindJSON(&order); err != nil {
			return nil, gin.H{"message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		pharmacyId := c.Param("pharmacyId")
		if order.Id == "" || order.Id == "@new" {
			order.Id = uuid.NewString()
		}
		// ensure pharmacyId is set from path
		order.PharmacyId = pharmacyId
		if order.Status == "" {
			order.Status = OrderStatusCreated
		}
		now := time.Now().UTC()
		order.CreatedAt = now
		order.UpdatedAt = now

		conflict := slices.IndexFunc(store.Orders, func(o Order) bool { return o.Id == order.Id })
		if conflict >= 0 {
			return nil, gin.H{"message": "Order already exists"}, http.StatusConflict
		}
		store.Orders = append(store.Orders, order)
		return store, order, http.StatusOK
	})
}

func UpdateOrder(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		orderId := c.Param("orderId")
		var order Order
		if err := c.ShouldBindJSON(&order); err != nil {
			return nil, gin.H{"message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		idx := slices.IndexFunc(store.Orders, func(o Order) bool { return o.Id == orderId })
		if idx < 0 {
			return nil, gin.H{"message": "Order not found"}, http.StatusNotFound
		}
		// enforce id/path consistency
		order.Id = orderId
		order.PharmacyId = store.Id
		order.UpdatedAt = time.Now().UTC()
		// preserve CreatedAt if incoming doesn't set it
		if store.Orders[idx].CreatedAt.IsZero() == false && order.CreatedAt.IsZero() {
			order.CreatedAt = store.Orders[idx].CreatedAt
		}
		store.Orders[idx] = order
		return store, order, http.StatusOK
	})
}

func UpdateOrderStatus(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		orderId := c.Param("orderId")
		var payload struct {
			Status OrderStatus `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			return nil, gin.H{"message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		idx := slices.IndexFunc(store.Orders, func(o Order) bool { return o.Id == orderId })
		if idx < 0 {
			return nil, gin.H{"message": "Order not found"}, http.StatusNotFound
		}
		store.Orders[idx].Status = payload.Status
		store.Orders[idx].UpdatedAt = time.Now().UTC()
		return store, store.Orders[idx], http.StatusOK
	})
}

func DeleteOrder(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		orderId := c.Param("orderId")
		idx := slices.IndexFunc(store.Orders, func(o Order) bool { return o.Id == orderId })
		if idx < 0 {
			return nil, gin.H{"message": "Order not found"}, http.StatusNotFound
		}
		store.Orders = append(store.Orders[:idx], store.Orders[idx+1:]...)
		return store, nil, http.StatusNoContent
	})
}

// Dispensings handlers
func GetDispensings(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		dispensings := store.Dispensings
		if dispensings == nil {
			dispensings = []Dispensing{}
		}
		return nil, dispensings, http.StatusOK
	})
}

func CreateDispensing(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		var dispensing Dispensing
		if err := c.ShouldBindJSON(&dispensing); err != nil {
			return nil, gin.H{"message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		if dispensing.MedicineId == "" {
			return nil, gin.H{"message": "medicineId is required"}, http.StatusBadRequest
		}
		if dispensing.Quantity <= 0 {
			return nil, gin.H{"message": "quantity must be positive"}, http.StatusBadRequest
		}

		// find the medicine and check stock
		medIdx := slices.IndexFunc(store.Medicines, func(m Medicine) bool { return m.Id == dispensing.MedicineId })
		if medIdx < 0 {
			return nil, gin.H{"message": "Medicine not found"}, http.StatusNotFound
		}
		med := store.Medicines[medIdx]
		if med.CurrentStock < dispensing.Quantity {
			return nil, gin.H{"message": "Insufficient stock", "available": med.CurrentStock}, http.StatusUnprocessableEntity
		}

		// fill in derived fields
		dispensing.MedicineName = med.Name
		if dispensing.Id == "" || dispensing.Id == "@new" {
			dispensing.Id = uuid.NewString()
		}
		dispensing.DispensedAt = time.Now().UTC()

		// decrement stock
		store.Medicines[medIdx].CurrentStock -= dispensing.Quantity

		store.Dispensings = append(store.Dispensings, dispensing)
		return store, dispensing, http.StatusOK
	})
}

func GetOrderReceipt(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		orderId := c.Param("orderId")

		if store.StockReceipts == nil {
			return nil, gin.H{"message": "Receipt not found"}, http.StatusNotFound
		}

		receiptIdx := slices.IndexFunc(store.StockReceipts, func(r StockReceipt) bool { return r.OrderId == orderId })
		if receiptIdx < 0 {
			return nil, gin.H{"message": "Receipt not found"}, http.StatusNotFound
		}

		return nil, store.StockReceipts[receiptIdx], http.StatusOK
	})
}

func ReceiveOrder(ctx *gin.Context) {
	getPharmacyFunc(ctx, func(c *gin.Context, store *PharmacyStore) (*PharmacyStore, interface{}, int) {
		orderId := c.Param("orderId")

		orderIdx := slices.IndexFunc(store.Orders, func(o Order) bool { return o.Id == orderId })
		if orderIdx < 0 {
			return nil, gin.H{"message": "Order not found"}, http.StatusNotFound
		}

		order := store.Orders[orderIdx]
		if order.Status != OrderStatusDelivered {
			return nil, gin.H{"message": "Only delivered orders can be received into stock"}, http.StatusConflict
		}

		if len(order.Items) == 0 {
			return nil, gin.H{"message": "Order has no items"}, http.StatusBadRequest
		}

		if store.StockReceipts != nil {
			if receiptIdx := slices.IndexFunc(store.StockReceipts, func(r StockReceipt) bool { return r.OrderId == orderId }); receiptIdx >= 0 {
				return nil, nil, http.StatusNoContent
			}
		}

		receipt := StockReceipt{
			Id:         uuid.NewString(),
			PharmacyId: store.Id,
			OrderId:    orderId,
			ReceivedAt: time.Now().UTC(),
			Items:      make([]StockReceiptItem, 0, len(order.Items)),
		}

		medicineIndexes := make([]int, len(order.Items))
		for i, item := range order.Items {
			if item.Quantity <= 0 {
				return nil, gin.H{"message": "Order item quantity must be positive", "medicineId": item.MedicineId}, http.StatusBadRequest
			}

			medIdx := slices.IndexFunc(store.Medicines, func(m Medicine) bool { return m.Id == item.MedicineId })
			if medIdx < 0 {
				newMedId := item.MedicineId
				if newMedId == "" {
					newMedId = uuid.NewString()
				}
				newMed := Medicine{
					Id:           newMedId,
					Name:         item.MedicineName,
					CurrentStock: 0,
				}
				store.Medicines = append(store.Medicines, newMed)
				medIdx = len(store.Medicines) - 1
			}

			medicineIndexes[i] = medIdx
			receipt.Items = append(receipt.Items, StockReceiptItem{
				MedicineId:   store.Medicines[medIdx].Id,
				MedicineName: item.MedicineName,
				Quantity:     item.Quantity,
			})
		}

		for i, item := range order.Items {
			store.Medicines[medicineIndexes[i]].CurrentStock += item.Quantity
		}

		store.StockReceipts = append(store.StockReceipts, receipt)
		return store, nil, http.StatusNoContent
	})
}
