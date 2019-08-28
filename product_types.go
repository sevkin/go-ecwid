package ecwid

type (
	// OptionValue is set of options that identifies this variation.
	OptionValue struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}

	// ProductOptionChoice - possible option selections for the types SELECT, CHECKBOX or RADIO.
	// This field is omitted for the product option with no selection
	// (e.g. text, datepicker or upload file options)
	ProductOptionChoice struct {
		Text              string       `json:"text,omitempty"`
		PriceModifier     int          `json:"priceModifier"`
		PriceModifierType ModifierType `json:"priceModifierType,omitempty"`
	}

	// ProductOption ...
	ProductOption struct {
		Type          ProductOptionType     `json:"type,omitempty"` // TODO One of SELECT, RADIO, CHECKBOX, TEXTFIELD, TEXTAREA, DATE, FILES
		Name          string                `json:"name,omitempty"`
		Required      bool                  `json:"required"`
		Choices       []ProductOptionChoice `json:"choices,omitempty"`
		DefaultChoice uint                  `json:"defaultChoice"`
	}

	// ShippingSettings of product
	ShippingSettings struct {
		Type            ShippingSettingsType `json:"type,omitempty"` // TODO One of: "GLOBAL_METHODS", "SELECTED_METHODS", "FLAT_RATE", "FREE_SHIPPING". "GLOBAL_METHODS"
		MethodMarkup    float32              `json:"methodMarkup,omitempty"`
		FlatRate        float32              `json:"flatRate,omitempty"`
		DisabledMethods []string             `json:"disabledMethods,omitempty"`
		EnabledMethods  []string             `json:"enabledMethods,omitempty"`
	}

	// WholesalePrice is element of array of variation’s wholesale price tiers
	// (quantity limit and price).
	WholesalePrice struct {
		Quantity uint    `json:"quantity"`
		Price    float32 `json:"price,omitempty"`
	}

	// RelatedCategory describes the “N random related products from a category” option
	RelatedCategory struct {
		Enabled      bool   `json:"enabled"`
		CategoryID   uint64 `json:"categoryId"`
		ProductCount uint   `json:"productCount"`
	}

	// RelatedProducts related or “You may also like” products of the product
	RelatedProducts struct {
		ProductIDs      []uint64        `json:"productIds,omitempty"`
		RelatedCategory RelatedCategory `json:"relatedCategory"`
	}

	// ProductDimensions is product dimensions info
	ProductDimensions struct {
		Length float32 `json:"length"`
		Width  float32 `json:"width"`
		Height float32 `json:"height"`
	}

	// TaxInfo contains detailed information about product’s taxes
	TaxInfo struct {
		DefaultLocationIncludedTaxRate uint     `json:"defaultLocationIncludedTaxRate"` // ???
		EnabledManualTaxes             []uint64 `json:"enabledManualTaxes,omitempty"`
	}

	// ProductImage contains images of product and their details
	ProductImage struct {
		ID               *string `json:"id"`
		OrderBy          *uint   `json:"orderBy"`
		IsMain           *bool   `json:"isMain"`
		Image160pxURL    string  `json:"image160pxUrl,omitempty"`
		Image400pxURL    string  `json:"image400pxUrl,omitempty"`
		Image800pxURL    string  `json:"image800pxUrl,omitempty"`
		Image1500pxURL   string  `json:"image1500pxUrl,omitempty"`
		ImageOriginalURL string  `json:"imageOriginalUrl,omitempty"`
	}

	// ProductMedia contains media files for a product (images)
	ProductMedia struct {
		Images []ProductImage `json:"images,omitempty"`
	}

	// TODO unfilled structs
	// GalleryImage struct {
	// }

	// CategoriesInfo ...
	CategoriesInfo struct {
		ID      uint64 `json:"id"`
		Enabled bool   `json:"enabled"`
	}

	// ProductFile struct {
	// }

	// ProductVariation ...
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
		Attributes      []Attribute      `json:"attributes,omitempty"`
	}

	// FavoritesStats struct {
	// }

	// ImageDetails is thumbnail image data.
	// The thumbnail size is specified in the store settings.
	ImageDetails struct {
		URL    string `json:"url"`
		Width  uint   `json:"width"`
		Height uint   `json:"height"`
	}
)

// Custom fields

type (
	// ProductOptionType ...
	ProductOptionType string

	// ShippingSettingsType ...
	ShippingSettingsType string
)

// ProductOptionType types
const (
	ProductOptionSelect    ProductOptionType = "SELECT"
	ProductOptionRadio     ProductOptionType = "RADIO"
	ProductOptionCheckbox  ProductOptionType = "CHECKBOX"
	ProductOptionTextfield ProductOptionType = "TEXTFIELD"
	ProductOptionTextarea  ProductOptionType = "TEXTAREA"
	ProductOptionDate      ProductOptionType = "DATE"
	ProductOptionFiles     ProductOptionType = "FILES"
)

// ShippingSettingsType types
const (
	ShippingSettingGlobalMethods   ShippingSettingsType = "GLOBAL_METHODS"
	ShippingSettingSelectedMethods ShippingSettingsType = "SELECTED_METHODS"
	ShippingSettingFlatFate        ShippingSettingsType = "FLAT_RATE"
	ShippingSettingFreeShipping    ShippingSettingsType = "FREE_SHIPPING"
)
