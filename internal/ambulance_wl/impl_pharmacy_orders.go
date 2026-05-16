package ambulance_wl

import (
	"github.com/gin-gonic/gin"
	"github.com/xmikova/ambulance-webapi/internal/pharmacy"
)

type implPharmacyOrdersAPI struct {
}

func NewPharmacyOrdersApi() PharmacyOrdersAPI {
	return &implPharmacyOrdersAPI{}
}

func (o implPharmacyOrdersAPI) CreateOrder(c *gin.Context) {
	pharmacy.CreateOrder(c)
}

func (o implPharmacyOrdersAPI) DeleteOrder(c *gin.Context) {
	pharmacy.DeleteOrder(c)
}

func (o implPharmacyOrdersAPI) GetOrder(c *gin.Context) {
	pharmacy.GetOrder(c)
}

func (o implPharmacyOrdersAPI) GetOrders(c *gin.Context) {
	pharmacy.GetOrders(c)
}

func (o implPharmacyOrdersAPI) UpdateOrder(c *gin.Context) {
	pharmacy.UpdateOrder(c)
}

func (o implPharmacyOrdersAPI) UpdateOrderStatus(c *gin.Context) {
	pharmacy.UpdateOrderStatus(c)
}
