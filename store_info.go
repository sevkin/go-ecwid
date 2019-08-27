package ecwid

type (
	// StoreProfile contains a basic information about an Ecwid store:
	// settings, store location, email, etc.
	StoreProfile struct {
		GeneralInfo GeneralInfo `json:"generalInfo"`
		Account     Account     `json:"account"`

		// TODO add more fields of Profile
		// settings	<Settings>	Store general settings
		// mailNotifications	<MailNotifications>	Mail notifications settings
		// company	<Company>	Company info
		// formatsAndUnits	<FormatsAndUnits>	Store formats/untis settings
		// languages	<Languages>	Store language settings
		// shipping	<Shipping>	Store shipping settings (only handling fee is included at the moment)
		// taxSettings	<TaxSettings>	Store taxes settings
		// zones	Array<Zone>	Store destination zones
		// businessRegistrationID	<BusinessRegistrationID>	Company registration ID, e.g. VAT reg number or company ID, which is set under Settings / Invoice in Control panel
		// legalPagesSettings	<LegalPagesSettingsDetails>	Legal pages settings for a store (System Settings → General → Legal Pages)
		// payment	<PaymentInfo>	Store payment settings information
		// featureToggles	<FeatureTogglesInfo>	Information about enabled/disabled new store features and their visibility in Ecwid Control Panel. Not provided via public token. Some of them are available in Ecwid JS API
		// designSettings	<DesignSettingsInfo>	Design settings of an Ecwid store. Can be overriden by updating store profile or by customizing design via JS config in storefront.
		// productFiltersSettings	<ProductFiltersSettings>	Settings for product filters in a store
	}

	// GeneralInfo - store basic data
	GeneralInfo struct {
		StoreID     uint64          `json:"storeId"`
		StoreURL    string          `json:"storeUrl"`
		StarterSite InstantSiteInfo `json:"starterSite"`
	}

	// InstantSiteInfo - details of Ecwid Instant site for account
	InstantSiteInfo struct {
		EcwidSubdomain string `json:"ecwidSubdomain"`
		CustomDomain   string `json:"customDomain"`
		GeneratedURL   string `json:"generatedUrl"`
		StoreLogoURL   string `json:"storeLogoUrl"`
	}

	// Account - store owner’s account data
	Account struct {
		AccountName       string   `json:"accountName"`
		AccountNickName   string   `json:"accountNickName"`
		AccountEmail      string   `json:"accountEmail"`
		AvailableFeatures []string `json:"availableFeatures"`
		WhiteLabel        bool     `json:"whiteLabel"`
	}
)

// StoreProfileGet returns basic information about an Ecwid store:
// settings, store location, email, etc.
func (c *Client) StoreProfileGet() (*StoreProfile, error) {
	response, err := c.R().Get("/profile")

	var result StoreProfile
	return &result, responseUnmarshal(response, err, &result)
}

// TODO add more store related methods
// func (c *Client) StoreProfileUpdate

// func (c *Client) ShippingOptionsGet
// func (c *Client) ShippingOptionAdd
// func (c *Client) ShippingOptionUpdate

// func (c *Client) StoreLogoUpload
// func (c *Client) StoreLogoRemove

// func (c *Client) InvoiceLogoUpload
// func (c *Client) InvoiceLogoRemove

// func (c *Client) EmailLogoUpload
// func (c *Client) EmailLogoRemove

// func (c *Client) StoreUpdateStatisticsGet

// func (c *Client) DeletedItemsStatisticsGet
