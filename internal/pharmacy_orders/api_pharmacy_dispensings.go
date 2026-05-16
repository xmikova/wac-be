package pharmacy_orders

import (
	"github.com/gin-gonic/gin"
)

type PharmacyDispensingsAPI interface {

	// GetDispensings Get /api/pharmacy/:pharmacyId/dispensings
	// Provides the list of medicine dispensings
	GetDispensings(c *gin.Context)

	// CreateDispensing Post /api/pharmacy/:pharmacyId/dispensings
	// Records a new medicine dispensing and decrements stock
	CreateDispensing(c *gin.Context)
}
