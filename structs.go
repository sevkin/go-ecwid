package ecwid

type (
	OptionValue struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}

	ProductOptionChoice struct {
		Text              string `json:"text,omitempty"`
		PriceModifier     int    `json:"priceModifier,omitempty"`
		PriceModifierType string `json:"priceModifierType,omitempty"`
	}

	ProductOption struct {
		Type          string                `json:"type,omitempty"`
		Name          string                `json:"name,omitempty"`
		Required      bool                  `json:"required,omitempty"`
		Choices       []ProductOptionChoice `json:"choices,omitempty"`
		DefaultChoice uint                  `json:"defaultChoice,omitempty"`
	}

	ShippingSettings struct {
		Type            string   `json:"type,omitempty"`
		MethodMarkup    float32  `json:"methodMarkup,omitempty"`
		FlatRate        float32  `json:"flatRate,omitempty"`
		DisabledMethods []string `json:"disabledMethods,omitempty"`
		EnabledMethods  []string `json:"enabledMethods,omitempty"`
	}

	WholesalePrice struct {
		Quantity uint    `json:"quantity,omitempty"`
		Price    float32 `json:"price,omitempty"`
	}

	AttributeValue struct {
		ID    uint64 `json:"id"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
		Type  string `json:"type,omitempty"`
		Show  string `json:"show,omitempty"`
		// Alias string `json:"alias,omitempty"`
	}

	RelatedCategory struct {
		Enabled      bool   `json:"enabled,omitempty"`
		CategoryID   uint64 `json:"categoryId,omitempty"`
		ProductCount uint   `json:"productCount,omitempty"`
	}

	RelatedProducts struct {
		ProductIDs      []uint64        `json:"productIds,omitempty"`
		RelatedCategory RelatedCategory `json:"relatedCategory,omitempty"`
	}

	ProductDimensions struct {
		Length float32 `json:"length,omitempty"`
		Width  float32 `json:"width,omitempty"`
		Height float32 `json:"height,omitempty"`
	}

	TaxInfo struct {
		DefaultLocationIncludedTaxRate uint     `json:"defaultLocationIncludedTaxRate,omitempty"` // ???
		EnabledManualTaxes             []uint64 `json:"enabledManualTaxes,omitempty"`
	}

	ProductImage struct {
		ID               uint64 `json:"id"`
		OrderBy          uint   `json:"orderBy,omitempty"`
		IsMain           bool   `json:"isMain,omitempty"`
		Image160pxURL    string `json:"image160pxUrl,omitempty"`
		Image400pxURL    string `json:"image400pxUrl,omitempty"`
		Image800pxURL    string `json:"image800pxUrl,omitempty"`
		Image1500pxURL   string `json:"image1500pxUrl,omitempty"`
		ImageOriginalURL string `json:"imageOriginalUrl,omitempty"`
	}

	ProductMedia struct {
		Images []ProductImage `json:"images,omitempty"`
	}

	// TODO unfilled structs
	// GalleryImage struct {
	// }

	CategoriesInfo struct {
		ID      uint64 `json:"id"`
		Enabled bool   `json:"enabled,omitempty"`
	}

	// ProductFile struct {
	// }

	ProductVariation struct {
		ID                uint64  `json:"id"`
		CombinationNumber uint    `json:"combinationNumber,omitempty"`
		Sku               string  `json:"sku,omitempty"`
		ThumbnailURL      string  `json:"thumbnailUrl,omitempty"`
		ImageURL          string  `json:"imageUrl,omitempty"`
		SmallThumbnailURL string  `json:"smallThumbnailUrl,omitempty"`
		HdThumbnailURL    string  `json:"hdThumbnailUrl,omitempty"`
		OriginalImageURL  string  `json:"originalImageUrl,omitempty"`
		Quantity          uint    `json:"quantity,omitempty"`
		Unlimited         bool    `json:"unlimited,omitempty"`
		Price             float32 `json:"price,omitempty"`
		Weight            float32 `json:"weight,omitempty"`
		WarningLimit      uint    `json:"warningLimit,omitempty"`
		CompareToPrice    float32 `json:"compareToPrice,omitempty"`

		Options         []OptionValue    `json:"options,omitempty"`
		WholesalePrices []WholesalePrice `json:"wholesalePrices,omitempty"`
		Attributes      []AttributeValue `json:"attributes,omitempty"`
	}

	// FavoritesStats struct {
	// }
)
