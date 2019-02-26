package ecwid

type (
	OptionValue struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}

	ProductOptionChoice struct {
		Text              string `json:"text,omitempty"`
		PriceModifier     int    `json:"priceModifier"`
		PriceModifierType string `json:"priceModifierType,omitempty"`
	}

	ProductOption struct {
		Type          string                `json:"type,omitempty"`
		Name          string                `json:"name,omitempty"`
		Required      bool                  `json:"required"`
		Choices       []ProductOptionChoice `json:"choices,omitempty"`
		DefaultChoice uint                  `json:"defaultChoice"`
	}

	ShippingSettings struct {
		Type            string   `json:"type,omitempty"`
		MethodMarkup    float32  `json:"methodMarkup,omitempty"`
		FlatRate        float32  `json:"flatRate,omitempty"`
		DisabledMethods []string `json:"disabledMethods,omitempty"`
		EnabledMethods  []string `json:"enabledMethods,omitempty"`
	}

	WholesalePrice struct {
		Quantity uint    `json:"quantity"`
		Price    float32 `json:"price,omitempty"`
	}

	// //

	Attribute struct {
		ID    uint64 `json:"id,omitempty"` // mandatory for update exist attributes can`t be 0
		Name  string `json:"name,omitempty"`
		Alias string `json:"alias,omitempty"`
		Type  string `json:"type,omitempty"` // CUSTOM, UPC, BRAND, GENDER, AGE_GROUP, COLOR, SIZE, PRICE_PER_UNIT, UNITS_IN_PRODUCT
		Show  string `json:"show,omitempty"` // NOTSHOW, DESCR, PRICE
	}

	AttributeValue struct {
		Attribute
		Value string `json:"value,omitempty"`
	}

	// //

	RelatedCategory struct {
		Enabled      bool   `json:"enabled"`
		CategoryID   uint64 `json:"categoryId"`
		ProductCount uint   `json:"productCount"`
	}

	RelatedProducts struct {
		ProductIDs      []uint64        `json:"productIds,omitempty"`
		RelatedCategory RelatedCategory `json:"relatedCategory"`
	}

	ProductDimensions struct {
		Length float32 `json:"length"`
		Width  float32 `json:"width"`
		Height float32 `json:"height"`
	}

	TaxInfo struct {
		DefaultLocationIncludedTaxRate uint     `json:"defaultLocationIncludedTaxRate"` // ???
		EnabledManualTaxes             []uint64 `json:"enabledManualTaxes,omitempty"`
	}

	ProductImage struct {
		ID               uint64 `json:"id"`
		OrderBy          uint   `json:"orderBy"`
		IsMain           bool   `json:"isMain"`
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
		Enabled bool   `json:"enabled"`
	}

	// ProductFile struct {
	// }

	ProductVariation struct {
		ID                uint64  `json:"id"`
		CombinationNumber uint    `json:"combinationNumber"`
		Sku               string  `json:"sku,omitempty"`
		ThumbnailURL      string  `json:"thumbnailUrl,omitempty"`
		ImageURL          string  `json:"imageUrl,omitempty"`
		SmallThumbnailURL string  `json:"smallThumbnailUrl,omitempty"`
		HdThumbnailURL    string  `json:"hdThumbnailUrl,omitempty"`
		OriginalImageURL  string  `json:"originalImageUrl,omitempty"`
		Quantity          uint    `json:"quantity"`
		Unlimited         bool    `json:"unlimited"`
		Price             float32 `json:"price,omitempty"`
		Weight            float32 `json:"weight,omitempty"`
		WarningLimit      uint    `json:"warningLimit"`
		CompareToPrice    float32 `json:"compareToPrice,omitempty"`

		Options         []OptionValue    `json:"options,omitempty"`
		WholesalePrices []WholesalePrice `json:"wholesalePrices,omitempty"`
		Attributes      []AttributeValue `json:"attributes,omitempty"`
	}

	// FavoritesStats struct {
	// }

	ImageDetails struct {
		URL    string `json:"url"`
		Width  uint   `json:"width"`
		Height uint   `json:"height"`
	}
)
