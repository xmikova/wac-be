package pharmacy

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
}
