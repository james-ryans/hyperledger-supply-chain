package model

type Seed struct {
	ID          string  `json:"id"`
	SupplierID  string  `json:"supplier_id"`
	VarietyName string  `json:"variety_name"`
	PlantAge    float32 `json:"plant_age"`
	PlantShape  string  `json:"plant_shape"`
	PlantHeight float32 `json:"plant_height"`
	LeafShape   string  `json:"leaf_shape"`
}

type SeedResponse struct {
	ID          string  `json:"id"`
	VarietyName string  `json:"variety_name"`
	PlantAge    float32 `json:"plant_age"`
	PlantShape  string  `json:"plant_shape"`
	PlantHeight float32 `json:"plant_height"`
	LeafShape   string  `json:"leaf_shape"`
}

func ToSeedResponse(seed *Seed) SeedResponse {
	return SeedResponse{
		ID:          seed.ID,
		VarietyName: seed.VarietyName,
		PlantAge:    seed.PlantAge,
		PlantShape:  seed.PlantShape,
		PlantHeight: seed.PlantHeight,
		LeafShape:   seed.LeafShape,
	}
}

func ToSeedsResponse(seeds *[]Seed) []SeedResponse {
	res := make([]SeedResponse, 0)
	for _, seed := range *seeds {
		res = append(res, ToSeedResponse(&seed))
	}

	return res
}

type SeedService interface {
	GetAllSeeds(channelID, ID string) (*[]Seed, error)
	GetSeedByID(channelID, ID string) (*Seed, error)
	CreateSeed(channelID string, seed *Seed) (*Seed, error)
	UpdateSeed(channelID string, seed *Seed) error
	DeleteSeed(channelID, ID string) error
}

type SeedRepository interface {
	FindAll(channelID, ID string) (*[]Seed, error)
	FindByID(channelID, ID string) (*Seed, error)
	Create(channelID string, seed *Seed) error
	Update(channelID string, seed *Seed) error
	Delete(channelID, ID string) error
}
