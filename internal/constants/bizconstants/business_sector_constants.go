// internal/constants/bizconstants/business_sector_constants.go
package bizconstants

// ================================================
// BUSINESS SECTOR CONSTANTS (What the business DOES)
// ================================================

const (
	SectorRetail        = "retail"
	SectorWholesale     = "wholesale"
	SectorFashion       = "fashion"
	SectorBeauty        = "beauty"
	SectorFood          = "food"
	SectorHealth        = "health"
	SectorAgriculture   = "agriculture"
	SectorConstruction  = "construction"
	SectorRealEstate    = "real_estate"
	SectorTransport     = "transport"
	SectorLogistics     = "logistics"
	SectorHospitality   = "hospitality"
	SectorTourism       = "tourism"
	SectorEducation     = "education"
	SectorProfessional  = "professional"
	SectorFinancial     = "financial"
	SectorTechnology    = "technology"
	SectorTelecom       = "telecom"
	SectorEnergy        = "energy"
	SectorManufacturing = "manufacturing"
	SectorMining        = "mining"
	SectorAutomotive    = "automotive"
	SectorEntertainment = "entertainment"
	SectorMedia         = "media"
	SectorSports        = "sports"
	SectorCreative      = "creative"
	SectorCommunity     = "community"
	SectorEnvironment   = "environment"
	SectorSecurity      = "security"
	SectorCleaning      = "cleaning"
	SectorVeterinary    = "veterinary"
)

// BusinessSectorInfo holds all information about a business sector
type BusinessSectorInfo struct {
	Name        string
	DisplayName string
	Description string
	Icon        string
}

// BusinessSectors is the single source of truth for all business sectors
var BusinessSectors = map[string]BusinessSectorInfo{
	SectorRetail: {
		Name:        SectorRetail,
		DisplayName: "Retail & Consumer Goods",
		Description: "Selling physical goods directly to consumers (dukas, supermarkets, shops)",
		Icon:        "shopping_bag",
	},
	SectorWholesale: {
		Name:        SectorWholesale,
		DisplayName: "Wholesale & Distribution",
		Description: "Bulk selling and distribution to other businesses (Muthurwa, Wakulima, wholesalers)",
		Icon:        "warehouse",
	},
	SectorFashion: {
		Name:        SectorFashion,
		DisplayName: "Fashion & Apparel",
		Description: "Clothing, accessories, footwear, and fashion items (boutiques, mitumba, tailors)",
		Icon:        "style",
	},
	SectorBeauty: {
		Name:        SectorBeauty,
		DisplayName: "Beauty & Personal Care",
		Description: "Salons, barbers, spas, cosmetics, and grooming services",
		Icon:        "beauty",
	},
	SectorFood: {
		Name:        SectorFood,
		DisplayName: "Food & Beverage",
		Description: "Restaurants, cafes, bakeries, food production, and beverage services",
		Icon:        "restaurant",
	},
	SectorHealth: {
		Name:        SectorHealth,
		DisplayName: "Health & Wellness",
		Description: "Clinics, pharmacies, hospitals, fitness centers, and wellness services",
		Icon:        "health",
	},
	SectorAgriculture: {
		Name:        SectorAgriculture,
		DisplayName: "Agriculture & Farming",
		Description: "Farming, livestock, agri-inputs, and agricultural services",
		Icon:        "agriculture",
	},
	SectorConstruction: {
		Name:        SectorConstruction,
		DisplayName: "Construction & Building",
		Description: "Building, property development, construction services, and real estate",
		Icon:        "construction",
	},
	SectorRealEstate: {
		Name:        SectorRealEstate,
		DisplayName: "Real Estate & Property",
		Description: "Property sales, rentals, property management, and real estate services",
		Icon:        "real_estate",
	},
	SectorTransport: {
		Name:        SectorTransport,
		DisplayName: "Transportation",
		Description: "Transportation services, matatus, buses, and ride-hailing",
		Icon:        "transport",
	},
	SectorLogistics: {
		Name:        SectorLogistics,
		DisplayName: "Logistics & Supply Chain",
		Description: "Delivery, courier, logistics, and supply chain services",
		Icon:        "delivery",
	},
	SectorHospitality: {
		Name:        SectorHospitality,
		DisplayName: "Hospitality",
		Description: "Hotels, guest houses, accommodation, and hospitality services",
		Icon:        "hotel",
	},
	SectorTourism: {
		Name:        SectorTourism,
		DisplayName: "Tourism & Travel",
		Description: "Travel agencies, tour operators, safaris, and tourism services",
		Icon:        "tourism",
	},
	SectorEducation: {
		Name:        SectorEducation,
		DisplayName: "Education & Training",
		Description: "Schools, colleges, universities, tutoring, and professional training",
		Icon:        "school",
	},
	SectorProfessional: {
		Name:        SectorProfessional,
		DisplayName: "Professional Services",
		Description: "Legal, accounting, consulting, marketing, and business services",
		Icon:        "briefcase",
	},
	SectorFinancial: {
		Name:        SectorFinancial,
		DisplayName: "Financial Services",
		Description: "Banking, insurance, investment, SACCOs, and microfinance",
		Icon:        "finance",
	},
	SectorTechnology: {
		Name:        SectorTechnology,
		DisplayName: "Technology & Digital",
		Description: "Phones, computers, gadgets, repairs, and tech services",
		Icon:        "devices",
	},
	SectorTelecom: {
		Name:        SectorTelecom,
		DisplayName: "Telecommunications",
		Description: "Mobile network operators, internet service providers, and communication services",
		Icon:        "telecom",
	},
	SectorEnergy: {
		Name:        SectorEnergy,
		DisplayName: "Energy & Utilities",
		Description: "Energy generation, distribution, and utility services",
		Icon:        "energy",
	},
	SectorManufacturing: {
		Name:        SectorManufacturing,
		DisplayName: "Manufacturing & Production",
		Description: "Manufacturing, production, and assembly of goods",
		Icon:        "manufacturing",
	},
	SectorMining: {
		Name:        SectorMining,
		DisplayName: "Mining & Extraction",
		Description: "Mineral extraction, quarrying, and natural resource exploration",
		Icon:        "mining",
	},
	SectorAutomotive: {
		Name:        SectorAutomotive,
		DisplayName: "Automotive",
		Description: "Car sales, repairs, spare parts, and automotive services",
		Icon:        "car",
	},
	SectorEntertainment: {
		Name:        SectorEntertainment,
		DisplayName: "Entertainment",
		Description: "Events, music, media, and entertainment services",
		Icon:        "entertainment",
	},
	SectorMedia: {
		Name:        SectorMedia,
		DisplayName: "Media & Communications",
		Description: "TV, radio, print, and digital media services",
		Icon:        "media",
	},
	SectorSports: {
		Name:        SectorSports,
		DisplayName: "Sports & Fitness",
		Description: "Sports clubs, gyms, fitness centers, and sports training",
		Icon:        "sports",
	},
	SectorCreative: {
		Name:        SectorCreative,
		DisplayName: "Creative Arts & Design",
		Description: "Art, design, animation, fashion design, and creative services",
		Icon:        "creative",
	},
	SectorCommunity: {
		Name:        SectorCommunity,
		DisplayName: "Community & Non-Profit",
		Description: "NGOs, charities, community organizations, and social enterprises",
		Icon:        "community",
	},
	SectorEnvironment: {
		Name:        SectorEnvironment,
		DisplayName: "Environmental Services",
		Description: "Environmental consulting, conservation, and sustainability services",
		Icon:        "environment",
	},
	SectorSecurity: {
		Name:        SectorSecurity,
		DisplayName: "Security Services",
		Description: "Security guard services, alarm systems, and security solutions",
		Icon:        "security",
	},
	SectorCleaning: {
		Name:        SectorCleaning,
		DisplayName: "Cleaning & Sanitation",
		Description: "Cleaning services, waste management, and sanitation services",
		Icon:        "cleaning",
	},
	SectorVeterinary: {
		Name:        SectorVeterinary,
		DisplayName: "Veterinary Services",
		Description: "Animal health, veterinary clinics, and livestock services",
		Icon:        "veterinary",
	},
}

// AllBusinessSectors defines all valid business sector names
var AllBusinessSectors = []string{
	SectorRetail,
	SectorWholesale,
	SectorFashion,
	SectorBeauty,
	SectorFood,
	SectorHealth,
	SectorAgriculture,
	SectorConstruction,
	SectorRealEstate,
	SectorTransport,
	SectorLogistics,
	SectorHospitality,
	SectorTourism,
	SectorEducation,
	SectorProfessional,
	SectorFinancial,
	SectorTechnology,
	SectorTelecom,
	SectorEnergy,
	SectorManufacturing,
	SectorMining,
	SectorAutomotive,
	SectorEntertainment,
	SectorMedia,
	SectorSports,
	SectorCreative,
	SectorCommunity,
	SectorEnvironment,
	SectorSecurity,
	SectorCleaning,
	SectorVeterinary,
}

// BusinessSectorDisplayNames returns display names for business sectors
var BusinessSectorDisplayNames = map[string]string{
	SectorRetail:        "Retail & Consumer Goods",
	SectorWholesale:     "Wholesale & Distribution",
	SectorFashion:       "Fashion & Apparel",
	SectorBeauty:        "Beauty & Personal Care",
	SectorFood:          "Food & Beverage",
	SectorHealth:        "Health & Wellness",
	SectorAgriculture:   "Agriculture & Farming",
	SectorConstruction:  "Construction & Building",
	SectorRealEstate:    "Real Estate & Property",
	SectorTransport:     "Transportation",
	SectorLogistics:     "Logistics & Supply Chain",
	SectorHospitality:   "Hospitality",
	SectorTourism:       "Tourism & Travel",
	SectorEducation:     "Education & Training",
	SectorProfessional:  "Professional Services",
	SectorFinancial:     "Financial Services",
	SectorTechnology:    "Technology & Digital",
	SectorTelecom:       "Telecommunications",
	SectorEnergy:        "Energy & Utilities",
	SectorManufacturing: "Manufacturing & Production",
	SectorMining:        "Mining & Extraction",
	SectorAutomotive:    "Automotive",
	SectorEntertainment: "Entertainment",
	SectorMedia:         "Media & Communications",
	SectorSports:        "Sports & Fitness",
	SectorCreative:      "Creative Arts & Design",
	SectorCommunity:     "Community & Non-Profit",
	SectorEnvironment:   "Environmental Services",
	SectorSecurity:      "Security Services",
	SectorCleaning:      "Cleaning & Sanitation",
	SectorVeterinary:    "Veterinary Services",
}

// BusinessSectorDescriptions returns descriptions for business sectors
var BusinessSectorDescriptions = map[string]string{
	SectorRetail:        "Selling physical goods directly to consumers (dukas, supermarkets, shops)",
	SectorWholesale:     "Bulk selling and distribution to other businesses (Muthurwa, Wakulima, wholesalers)",
	SectorFashion:       "Clothing, accessories, footwear, and fashion items (boutiques, mitumba, tailors)",
	SectorBeauty:        "Salons, barbers, spas, cosmetics, and grooming services",
	SectorFood:          "Restaurants, cafes, bakeries, food production, and beverage services",
	SectorHealth:        "Clinics, pharmacies, hospitals, fitness centers, and wellness services",
	SectorAgriculture:   "Farming, livestock, agri-inputs, and agricultural services",
	SectorConstruction:  "Building, property development, construction services, and real estate",
	SectorRealEstate:    "Property sales, rentals, property management, and real estate services",
	SectorTransport:     "Transportation services, matatus, buses, and ride-hailing",
	SectorLogistics:     "Delivery, courier, logistics, and supply chain services",
	SectorHospitality:   "Hotels, guest houses, accommodation, and hospitality services",
	SectorTourism:       "Travel agencies, tour operators, safaris, and tourism services",
	SectorEducation:     "Schools, colleges, universities, tutoring, and professional training",
	SectorProfessional:  "Legal, accounting, consulting, marketing, and business services",
	SectorFinancial:     "Banking, insurance, investment, SACCOs, and microfinance",
	SectorTechnology:    "Phones, computers, gadgets, repairs, and tech services",
	SectorTelecom:       "Mobile network operators, internet service providers, and communication services",
	SectorEnergy:        "Energy generation, distribution, and utility services",
	SectorManufacturing: "Manufacturing, production, and assembly of goods",
	SectorMining:        "Mineral extraction, quarrying, and natural resource exploration",
	SectorAutomotive:    "Car sales, repairs, spare parts, and automotive services",
	SectorEntertainment: "Events, music, media, and entertainment services",
	SectorMedia:         "TV, radio, print, and digital media services",
	SectorSports:        "Sports clubs, gyms, fitness centers, and sports training",
	SectorCreative:      "Art, design, animation, fashion design, and creative services",
	SectorCommunity:     "NGOs, charities, community organizations, and social enterprises",
	SectorEnvironment:   "Environmental consulting, conservation, and sustainability services",
	SectorSecurity:      "Security guard services, alarm systems, and security solutions",
	SectorCleaning:      "Cleaning services, waste management, and sanitation services",
	SectorVeterinary:    "Animal health, veterinary clinics, and livestock services",
}

// BusinessSectorIcons returns icons for business sectors
var BusinessSectorIcons = map[string]string{
	SectorRetail:        "shopping_bag",
	SectorWholesale:     "warehouse",
	SectorFashion:       "style",
	SectorBeauty:        "beauty",
	SectorFood:          "restaurant",
	SectorHealth:        "health",
	SectorAgriculture:   "agriculture",
	SectorConstruction:  "construction",
	SectorRealEstate:    "real_estate",
	SectorTransport:     "transport",
	SectorLogistics:     "delivery",
	SectorHospitality:   "hotel",
	SectorTourism:       "tourism",
	SectorEducation:     "school",
	SectorProfessional:  "briefcase",
	SectorFinancial:     "finance",
	SectorTechnology:    "devices",
	SectorTelecom:       "telecom",
	SectorEnergy:        "energy",
	SectorManufacturing: "manufacturing",
	SectorMining:        "mining",
	SectorAutomotive:    "car",
	SectorEntertainment: "entertainment",
	SectorMedia:         "media",
	SectorSports:        "sports",
	SectorCreative:      "creative",
	SectorCommunity:     "community",
	SectorEnvironment:   "environment",
	SectorSecurity:      "security",
	SectorCleaning:      "cleaning",
	SectorVeterinary:    "veterinary",
}