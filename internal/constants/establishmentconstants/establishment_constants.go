// internal/constants/establishmentconstants/establishment_constants.go
package establishmentconstants

// ================================================
// ESTABLISHMENT TYPE CONSTANTS
// ================================================

const (
	EstablishmentShop           = "shop"
	EstablishmentStall          = "stall"
	EstablishmentKiosk          = "kiosk"
	EstablishmentMarket         = "market"
	EstablishmentBooth          = "booth"
	EstablishmentOutlet         = "outlet"
	EstablishmentSupermarket    = "supermarket"
	EstablishmentOnlineStore    = "online_store"
	EstablishmentDigitalService = "digital_service"
	EstablishmentDeliveryOnly   = "delivery_only"
	EstablishmentVirtualOffice  = "virtual_office"
	EstablishmentMobileVendor   = "mobile_vendor"
	EstablishmentHomeBased      = "home_based"
	EstablishmentAgent          = "agent"
	EstablishmentHybrid         = "hybrid"
)

// EstablishmentTypeInfo holds all information about an establishment type
type EstablishmentTypeInfo struct {
	Name        string
	DisplayName string
	Description string
	Category    string // physical, digital, mobile, home, hybrid
}

// EstablishmentTypes is the single source of truth for all establishment types
var EstablishmentTypes = map[string]EstablishmentTypeInfo{
	EstablishmentShop: {
		Name:        EstablishmentShop,
		DisplayName: "Shop",
		Description: "Permanent storefront or retail space",
		Category:    "physical",
	},
	EstablishmentStall: {
		Name:        EstablishmentStall,
		DisplayName: "Stall",
		Description: "Small unit in a market or open-air space",
		Category:    "physical",
	},
	EstablishmentKiosk: {
		Name:        EstablishmentKiosk,
		DisplayName: "Kiosk",
		Description: "Small standalone structure, often street-side",
		Category:    "physical",
	},
	EstablishmentMarket: {
		Name:        EstablishmentMarket,
		DisplayName: "Market",
		Description: "Collection of vendors selling various items",
		Category:    "physical",
	},
	EstablishmentBooth: {
		Name:        EstablishmentBooth,
		DisplayName: "Booth",
		Description: "Temporary stall for events, exhibitions, or trade fairs",
		Category:    "physical",
	},
	EstablishmentOutlet: {
		Name:        EstablishmentOutlet,
		DisplayName: "Outlet",
		Description: "Brand store or official retail outlet",
		Category:    "physical",
	},
	EstablishmentSupermarket: {
		Name:        EstablishmentSupermarket,
		DisplayName: "Supermarket",
		Description: "Large grocery store with wide selection of products",
		Category:    "physical",
	},
	EstablishmentOnlineStore: {
		Name:        EstablishmentOnlineStore,
		DisplayName: "Online Store",
		Description: "Business operating entirely online (social media sellers, e-commerce)",
		Category:    "digital",
	},
	EstablishmentDigitalService: {
		Name:        EstablishmentDigitalService,
		DisplayName: "Digital Service",
		Description: "Freelancers, consultants, content creators, remote services",
		Category:    "digital",
	},
	EstablishmentDeliveryOnly: {
		Name:        EstablishmentDeliveryOnly,
		DisplayName: "Delivery Only",
		Description: "Business that delivers directly to customers (food, groceries, courier)",
		Category:    "digital",
	},
	EstablishmentVirtualOffice: {
		Name:        EstablishmentVirtualOffice,
		DisplayName: "Virtual Office",
		Description: "Business with virtual office address (no physical storefront)",
		Category:    "digital",
	},
	EstablishmentMobileVendor: {
		Name:        EstablishmentMobileVendor,
		DisplayName: "Mobile Vendor",
		Description: "Street vendors, hawkers, mobile food trucks",
		Category:    "mobile",
	},
	EstablishmentHomeBased: {
		Name:        EstablishmentHomeBased,
		DisplayName: "Home-Based",
		Description: "Business operating from home (baking, tailoring, catering, daycare)",
		Category:    "home",
	},
	EstablishmentAgent: {
		Name:        EstablishmentAgent,
		DisplayName: "Agent",
		Description: "M-Pesa agents, insurance agents, booking agents",
		Category:    "physical",
	},
	EstablishmentHybrid: {
		Name:        EstablishmentHybrid,
		DisplayName: "Hybrid",
		Description: "Both physical store and online presence",
		Category:    "hybrid",
	},
}

// ================================================
// HELPER SLICES AND MAPS
// ================================================

// AllEstablishmentTypes defines all valid establishment type names
var AllEstablishmentTypes = []string{
	EstablishmentShop,
	EstablishmentStall,
	EstablishmentKiosk,
	EstablishmentMarket,
	EstablishmentBooth,
	EstablishmentOutlet,
	EstablishmentSupermarket,
	EstablishmentOnlineStore,
	EstablishmentDigitalService,
	EstablishmentDeliveryOnly,
	EstablishmentVirtualOffice,
	EstablishmentMobileVendor,
	EstablishmentHomeBased,
	EstablishmentAgent,
	EstablishmentHybrid,
}

// EstablishmentTypeDisplayNames returns display names for establishment types
var EstablishmentTypeDisplayNames = map[string]string{
	EstablishmentShop:           "Shop",
	EstablishmentStall:          "Stall",
	EstablishmentKiosk:          "Kiosk",
	EstablishmentMarket:         "Market",
	EstablishmentBooth:          "Booth",
	EstablishmentOutlet:         "Outlet",
	EstablishmentSupermarket:    "Supermarket",
	EstablishmentOnlineStore:    "Online Store",
	EstablishmentDigitalService: "Digital Service",
	EstablishmentDeliveryOnly:   "Delivery Only",
	EstablishmentVirtualOffice:  "Virtual Office",
	EstablishmentMobileVendor:   "Mobile Vendor",
	EstablishmentHomeBased:      "Home-Based",
	EstablishmentAgent:          "Agent",
	EstablishmentHybrid:         "Hybrid",
}

// EstablishmentTypeDescriptions returns descriptions for establishment types
var EstablishmentTypeDescriptions = map[string]string{
	EstablishmentShop:           "Permanent storefront or retail space",
	EstablishmentStall:          "Small unit in a market or open-air space",
	EstablishmentKiosk:          "Small standalone structure, often street-side",
	EstablishmentMarket:         "Collection of vendors selling various items",
	EstablishmentBooth:          "Temporary stall for events, exhibitions, or trade fairs",
	EstablishmentOutlet:         "Brand store or official retail outlet",
	EstablishmentSupermarket:    "Large grocery store with wide selection of products",
	EstablishmentOnlineStore:    "Business operating entirely online (social media sellers, e-commerce)",
	EstablishmentDigitalService: "Freelancers, consultants, content creators, remote services",
	EstablishmentDeliveryOnly:   "Business that delivers directly to customers (food, groceries, courier)",
	EstablishmentVirtualOffice:  "Business with virtual office address (no physical storefront)",
	EstablishmentMobileVendor:   "Street vendors, hawkers, mobile food trucks",
	EstablishmentHomeBased:      "Business operating from home (baking, tailoring, catering, daycare)",
	EstablishmentAgent:          "M-Pesa agents, insurance agents, booking agents",
	EstablishmentHybrid:         "Both physical store and online presence",
}

// EstablishmentTypeIcons returns icons for establishment types
var EstablishmentTypeIcons = map[string]string{
	EstablishmentShop:           "store",
	EstablishmentStall:          "stall",
	EstablishmentKiosk:          "kiosk",
	EstablishmentMarket:         "market",
	EstablishmentBooth:          "booth",
	EstablishmentOutlet:         "outlet",
	EstablishmentSupermarket:    "supermarket",
	EstablishmentOnlineStore:    "online_store",
	EstablishmentDigitalService: "digital_service",
	EstablishmentDeliveryOnly:   "delivery",
	EstablishmentVirtualOffice:  "virtual_office",
	EstablishmentMobileVendor:   "mobile_vendor",
	EstablishmentHomeBased:      "home",
	EstablishmentAgent:          "agent",
	EstablishmentHybrid:         "hybrid",
}

// EstablishmentTypeCategories returns category groups for establishment types
var EstablishmentTypeCategories = map[string][]string{
	"physical": {
		EstablishmentShop,
		EstablishmentStall,
		EstablishmentKiosk,
		EstablishmentMarket,
		EstablishmentBooth,
		EstablishmentOutlet,
		EstablishmentSupermarket,
		EstablishmentAgent,
	},
	"digital": {
		EstablishmentOnlineStore,
		EstablishmentDigitalService,
		EstablishmentDeliveryOnly,
		EstablishmentVirtualOffice,
	},
	"mobile": {
		EstablishmentMobileVendor,
	},
	"home": {
		EstablishmentHomeBased,
	},
	"hybrid": {
		EstablishmentHybrid,
	},
}