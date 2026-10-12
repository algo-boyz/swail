package main

// Seller matches the frontend Seller interface
type Seller struct {
	ID       int    `json:"ID"`
	Username string `json:"Username"`
	IsPro    bool   `json:"IsPro"`
}

// Product matches the frontend Product interface
type Product struct {
	ItemID         string   `json:"ItemID"`
	Category       int      `json:"Category"`
	Title          string   `json:"Title"`
	Description    string   `json:"Description"`
	Weight         float64  `json:"Weight"` // grams
	LotURL         string   `json:"LotURL"`
	ImageURLs      []string `json:"ImageURLs"`
	MatchIDs       []string `json:"MatchIDs"`
	SellerID       int      `json:"SellerID"`
	AskingPrice    *int64   `json:"AskingPrice"`    // cents, nullable
	HighestBidEUR  *int64   `json:"HighestBidEUR"`  // cents, nullable
	DeliveryFeeEUR *int64   `json:"DeliveryFeeEUR"` // cents, nullable
	AiEstimate     *int64   `json:"AiEstimate"`     // cents, nullable
	MeltValue      int64    `json:"MeltValue"`      // cents
	FavoriteCount  int      `json:"FavoriteCount"`
	ObjectType     string   `json:"ObjectType"`
	ArbitrageRatio float64  `json:"ArbitrageRatio"`
	Confidence     float64  `json:"Confidence"`
}

// Listing is the combined view returned to the frontend
type Listing struct {
	Product Product `json:"Product"`
	Seller  Seller  `json:"Seller"`
}

// MpService exposes listing data to the frontend
type MpService struct{}

// GetListings returns the current set of Marktplaats deals.
// Replace the mock data with a real DB / scraper call later.
func (m *MpService) GetListings() []Listing {
	// Helper for optional int64 pointers
	ptr := func(v int64) *int64 { return &v }

	return []Listing{
		{
			Product: Product{
				ItemID:         "m123456789",
				Category:       12,
				Title:          "Gouden 10 Gulden Wilhelmina 1912",
				Description:    "Fraaie staat, originele randschrift.",
				Weight:         6.72,
				LotURL:         "https://www.marktplaats.nl/v/m123456789",
				ImageURLs:      []string{"https://images.marktplaats.com/api/images/m123456789/1"},
				MatchIDs:       []string{},
				SellerID:       8821,
				AskingPrice:    ptr(48500),
				HighestBidEUR:  nil,
				DeliveryFeeEUR: ptr(695),
				AiEstimate:     ptr(51200),
				MeltValue:      49800,
				FavoriteCount:  14,
				ObjectType:     "gold_coin",
				ArbitrageRatio: 0.027,
				Confidence:     0.91,
			},
			Seller: Seller{ID: 8821, Username: "goudhandelaar_nl", IsPro: true},
		},
		{
			Product: Product{
				ItemID:         "m987654321",
				Category:       12,
				Title:          "Zilveren Rijksdaalder 1786 Utrecht",
				Description:    "Mooie patina, geen restauratie.",
				Weight:         28.1,
				LotURL:         "https://www.marktplaats.nl/v/m987654321",
				ImageURLs:      []string{},
				MatchIDs:       []string{},
				SellerID:       4412,
				AskingPrice:    ptr(8900),
				HighestBidEUR:  nil,
				DeliveryFeeEUR: ptr(495),
				AiEstimate:     ptr(11200),
				MeltValue:      10450,
				FavoriteCount:  3,
				ObjectType:     "silver_coin",
				ArbitrageRatio: 0.174,
				Confidence:     0.78,
			},
			Seller: Seller{ID: 4412, Username: "oudemuntjes", IsPro: false},
		},
		{
			Product: Product{
				ItemID:         "m555666777",
				Category:       12,
				Title:          "Gouden 20 Frank Leopold II 1877",
				Description:    "Belgische gouden munt, goede staat.",
				Weight:         6.45,
				LotURL:         "https://www.marktplaats.nl/v/m555666777",
				ImageURLs:      []string{"https://images.marktplaats.com/api/images/m555666777/1"},
				MatchIDs:       []string{},
				SellerID:       1099,
				AskingPrice:    ptr(51200),
				HighestBidEUR:  nil,
				DeliveryFeeEUR: ptr(0),
				AiEstimate:     ptr(49800),
				MeltValue:      47900,
				FavoriteCount:  27,
				ObjectType:     "gold_coin",
				ArbitrageRatio: -0.064,
				Confidence:     0.85,
			},
			Seller: Seller{ID: 1099, Username: "antiek_goud", IsPro: true},
		},
	}
}
