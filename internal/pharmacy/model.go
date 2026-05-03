package pharmacy

import "time"

type Medicine struct {
	Id              string `json:"id" bson:"id"`
	Name            string `json:"name" bson:"name"`
	ActiveSubstance string `json:"activeSubstance" bson:"activeSubstance"`
	Dosage          string `json:"dosage" bson:"dosage"`
	BatchNumber     string `json:"batchNumber" bson:"batchNumber"`
	ExpiryDate      string `json:"expiryDate" bson:"expiryDate"`
	MinStock        int    `json:"minStock" bson:"minStock"`
	CurrentStock    int    `json:"currentStock" bson:"currentStock"`
}

type PharmacyStore struct {
	Id        string     `json:"id" bson:"id"`
	Medicines []Medicine `json:"medicines" bson:"medicines"`
	Orders    []Order    `json:"orders" bson:"orders"`
}

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "created"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	MedicineId   string  `json:"medicineId" bson:"medicineId"`
	MedicineName string  `json:"medicineName" bson:"medicineName"`
	Quantity     int     `json:"quantity" bson:"quantity"`
	Unit         string  `json:"unit" bson:"unit"`
	UnitPrice    float64 `json:"unitPrice" bson:"unitPrice"`
	TotalPrice   float64 `json:"totalPrice" bson:"totalPrice"`
}

type Order struct {
	Id         string      `json:"id" bson:"id"`
	PharmacyId string      `json:"pharmacyId" bson:"pharmacyId"`
	SupplierId string      `json:"supplierId,omitempty" bson:"supplierId,omitempty"`
	Items      []OrderItem `json:"items" bson:"items"`
	Status     OrderStatus `json:"status" bson:"status"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt" bson:"updatedAt"`
	CreatedBy  string      `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	Notes      string      `json:"notes,omitempty" bson:"notes,omitempty"`
}
