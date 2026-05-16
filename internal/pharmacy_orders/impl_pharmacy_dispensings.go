package pharmacy_orders

import (
	"github.com/gin-gonic/gin"
	"github.com/xmikova/ambulance-webapi/internal/pharmacy"
)

type implPharmacyDispensingsAPI struct {
}

func NewPharmacyDispensingsApi() PharmacyDispensingsAPI {
	return &implPharmacyDispensingsAPI{}
}

func (d implPharmacyDispensingsAPI) GetDispensings(c *gin.Context) {
	pharmacy.GetDispensings(c)
}

func (d implPharmacyDispensingsAPI) CreateDispensing(c *gin.Context) {
	pharmacy.CreateDispensing(c)
}
