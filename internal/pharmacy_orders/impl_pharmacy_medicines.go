package pharmacy_orders

import (
	"github.com/gin-gonic/gin"
	"github.com/xmikova/ambulance-webapi/internal/pharmacy"
)

type implPharmacyMedicinesAPI struct {
}

func NewPharmacyMedicinesApi() PharmacyMedicinesAPI {
	return &implPharmacyMedicinesAPI{}
}

func (m implPharmacyMedicinesAPI) GetMedicines(c *gin.Context) {
	pharmacy.GetMedicines(c)
}

func (m implPharmacyMedicinesAPI) CreateMedicine(c *gin.Context) {
	pharmacy.CreateMedicine(c)
}

func (m implPharmacyMedicinesAPI) GetMedicine(c *gin.Context) {
	pharmacy.GetMedicine(c)
}

func (m implPharmacyMedicinesAPI) UpdateMedicine(c *gin.Context) {
	pharmacy.UpdateMedicine(c)
}

func (m implPharmacyMedicinesAPI) DeleteMedicine(c *gin.Context) {
	pharmacy.DeleteMedicine(c)
}
