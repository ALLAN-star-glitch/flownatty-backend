// internal/constants/productconstants/product_subcategory_constants.go
package productconstants

// ================================================
// PRODUCT/SERVICE SUBCATEGORY CONSTANTS
// ================================================

const (
	// Women's Clothing
	SubcategoryDresses         = "dresses"
	SubcategoryTopsBlouses     = "tops_blouses"
	SubcategorySkirts          = "skirts"
	SubcategoryPantsTrousers   = "pants_trousers"
	SubcategoryJacketsCoats    = "jackets_coats"
	SubcategoryTraditionalWear = "traditional_wear"
	SubcategoryActivewear      = "activewear"

	// Men's Clothing
	SubcategoryShirts          = "shirts"
	SubcategoryTShirts         = "t_shirts"
	SubcategorySuitsBlazers    = "suits_blazers"
	SubcategoryMensPants       = "mens_pants"
	SubcategoryTraditionalMen  = "traditional_men"

	// Shoes & Footwear
	SubcategoryWomensShoes     = "womens_shoes"
	SubcategoryMensShoes       = "mens_shoes"
	SubcategoryKidsShoes       = "kids_shoes"
	SubcategoryAthleticShoes   = "athletic_shoes"
	SubcategorySandals         = "sandals"

	// Accessories
	SubcategoryBags            = "bags"
	SubcategoryBelts           = "belts"
	SubcategoryHatsCaps        = "hats_caps"
	SubcategoryScarvesWraps    = "scarves_wraps"
	SubcategoryWatches         = "watches"

	// Makeup
	SubcategoryFoundation      = "foundation"
	SubcategoryLipstick        = "lipstick"
	SubcategoryMascara         = "mascara"
	SubcategoryEyeshadow       = "eyeshadow"
	SubcategoryConcealer       = "concealer"

	// Skincare
	SubcategoryMoisturizers    = "moisturizers"
	SubcategoryCleansers       = "cleansers"
	SubcategorySunscreen       = "sunscreen"
	SubcategoryAntiAging       = "anti_aging"

	// Electronics
	SubcategorySmartphones     = "smartphones"
	SubcategoryPhoneCases      = "phone_cases"
	SubcategoryChargers        = "chargers"
	SubcategoryPowerBanks      = "power_banks"
	SubcategoryScreenProtectors = "screen_protectors"

	// Health Services
	SubcategoryGeneralPractice = "general_practice"
	SubcategoryPediatrics      = "pediatrics"
	SubcategoryCardiology      = "cardiology"
	SubcategoryOrthopedics     = "orthopedics"
	SubcategoryGynecology      = "gynecology"

	// Beauty Services
	SubcategoryHaircuts        = "haircuts"
	SubcategoryBraidsWeaves    = "braids_weaves"
	SubcategoryLocsDreads      = "locs_dreads"
	SubcategoryHairColor       = "hair_color"
	SubcategoryHairTreatments  = "hair_treatments"

	// Event Services
	SubcategoryWeddingPlanning = "wedding_planning"
	SubcategoryCorporateEvents = "corporate_events"
	SubcategoryBirthdayParties = "birthday_parties"
	SubcategoryAnniversaries   = "anniversaries"
	SubcategoryGraduationParties = "graduation_parties"

	// Plumbing Services
	SubcategoryPlumbingRepairs = "plumbing_repairs"
	SubcategoryPlumbingInstall = "plumbing_install"
	SubcategoryDrainage        = "drainage"

	// Electrical Services
	SubcategoryElectricalWiring = "electrical_wiring"
	SubcategoryElectricalInstall = "electrical_install"
	SubcategoryElectricalRepairs = "electrical_repairs"

	// Fitness Services
	SubcategoryOneOnOneTraining = "one_on_one_training"
	SubcategoryGroupTraining    = "group_training"
	SubcategoryOnlineTraining   = "online_training"

	// Yoga & Pilates
	SubcategoryHathaYoga       = "hatha_yoga"
	SubcategoryVinyasaYoga     = "vinyasa_yoga"
	SubcategoryPilates         = "pilates"

	// Professional Services
	SubcategoryCorporateLaw    = "corporate_law"
	SubcategoryFamilyLaw       = "family_law"
	SubcategoryCriminalLaw     = "criminal_law"
	SubcategoryPropertyLaw     = "property_law"

	// IT Services
	SubcategoryWebDevelopment  = "web_development"
	SubcategoryAppDevelopment  = "app_development"
	SubcategorySoftwareDevelopment = "software_development"
	SubcategoryITSupport       = "it_support"
)

// ProductSubcategoryInfo holds all information about a product/service subcategory
type ProductSubcategoryInfo struct {
	Name        string
	DisplayName string
	Description string
	Category    string // The category this subcategory belongs to
	Icon        string
}

// ProductSubcategories is the single source of truth for all product/service subcategories
var ProductSubcategories = map[string]ProductSubcategoryInfo{
	// ========================================
	// WOMEN'S CLOTHING SUBCATEGORIES
	// ========================================
	SubcategoryDresses: {
		Name:        SubcategoryDresses,
		DisplayName: "Dresses",
		Description: "Maxi, midi, mini, cocktail, evening dresses",
		Category:    CategoryWomenClothing,
		Icon:        "dresses",
	},
	SubcategoryTopsBlouses: {
		Name:        SubcategoryTopsBlouses,
		DisplayName: "Tops & Blouses",
		Description: "Shirts, blouses, crop tops, tank tops",
		Category:    CategoryWomenClothing,
		Icon:        "tops",
	},
	SubcategorySkirts: {
		Name:        SubcategorySkirts,
		DisplayName: "Skirts",
		Description: "Mini, midi, maxi, pencil, A-line skirts",
		Category:    CategoryWomenClothing,
		Icon:        "skirts",
	},
	SubcategoryPantsTrousers: {
		Name:        SubcategoryPantsTrousers,
		DisplayName: "Pants & Trousers",
		Description: "Jeans, trousers, leggings, palazzos",
		Category:    CategoryWomenClothing,
		Icon:        "pants",
	},
	SubcategoryJacketsCoats: {
		Name:        SubcategoryJacketsCoats,
		DisplayName: "Jackets & Coats",
		Description: "Blazers, jackets, coats, cardigans",
		Category:    CategoryWomenClothing,
		Icon:        "jackets",
	},
	SubcategoryTraditionalWear: {
		Name:        SubcategoryTraditionalWear,
		DisplayName: "Traditional Wear",
		Description: "Kitenge, kanga, khanga, cultural attire",
		Category:    CategoryWomenClothing,
		Icon:        "traditional_women",
	},
	SubcategoryActivewear: {
		Name:        SubcategoryActivewear,
		DisplayName: "Activewear",
		Description: "Sports bras, leggings, gym wear",
		Category:    CategoryWomenClothing,
		Icon:        "activewear",
	},

	// ========================================
	// MEN'S CLOTHING SUBCATEGORIES
	// ========================================
	SubcategoryShirts: {
		Name:        SubcategoryShirts,
		DisplayName: "Shirts",
		Description: "Casual, formal, polo, button-down shirts",
		Category:    CategoryMenClothing,
		Icon:        "shirts",
	},
	SubcategoryTShirts: {
		Name:        SubcategoryTShirts,
		DisplayName: "T-shirts",
		Description: "Plain, graphic, polo, henley t-shirts",
		Category:    CategoryMenClothing,
		Icon:        "tshirts",
	},
	SubcategorySuitsBlazers: {
		Name:        SubcategorySuitsBlazers,
		DisplayName: "Suits & Blazers",
		Description: "Business suits, blazers, sports jackets",
		Category:    CategoryMenClothing,
		Icon:        "suits",
	},
	SubcategoryMensPants: {
		Name:        SubcategoryMensPants,
		DisplayName: "Pants",
		Description: "Jeans, chinos, dress pants, shorts",
		Category:    CategoryMenClothing,
		Icon:        "mens_pants",
	},
	SubcategoryTraditionalMen: {
		Name:        SubcategoryTraditionalMen,
		DisplayName: "Traditional Wear",
		Description: "Cultural attire, kanzus, agbadas",
		Category:    CategoryMenClothing,
		Icon:        "traditional_men",
	},

	// ========================================
	// SHOES & FOOTWEAR SUBCATEGORIES
	// ========================================
	SubcategoryWomensShoes: {
		Name:        SubcategoryWomensShoes,
		DisplayName: "Women's Shoes",
		Description: "Heels, flats, sandals, boots, sneakers",
		Category:    CategoryShoesFootwear,
		Icon:        "womens_shoes",
	},
	SubcategoryMensShoes: {
		Name:        SubcategoryMensShoes,
		DisplayName: "Men's Shoes",
		Description: "Loafers, oxfords, sneakers, boots, sandals",
		Category:    CategoryShoesFootwear,
		Icon:        "mens_shoes",
	},
	SubcategoryKidsShoes: {
		Name:        SubcategoryKidsShoes,
		DisplayName: "Children's Shoes",
		Description: "Kids' shoes, school shoes, sandals, boots",
		Category:    CategoryShoesFootwear,
		Icon:        "kids_shoes",
	},
	SubcategoryAthleticShoes: {
		Name:        SubcategoryAthleticShoes,
		DisplayName: "Athletic Shoes",
		Description: "Running shoes, training shoes, sports shoes",
		Category:    CategoryShoesFootwear,
		Icon:        "athletic_shoes",
	},
	SubcategorySandals: {
		Name:        SubcategorySandals,
		DisplayName: "Traditional Sandals",
		Description: "Maasai sandals, leather sandals",
		Category:    CategoryShoesFootwear,
		Icon:        "sandals",
	},

	// ========================================
	// ACCESSORIES SUBCATEGORIES
	// ========================================
	SubcategoryBags: {
		Name:        SubcategoryBags,
		DisplayName: "Bags",
		Description: "Handbags, backpacks, totes, clutches, travel bags",
		Category:    CategoryAccessories,
		Icon:        "bags",
	},
	SubcategoryBelts: {
		Name:        SubcategoryBelts,
		DisplayName: "Belts",
		Description: "Leather, fabric, formal, casual belts",
		Category:    CategoryAccessories,
		Icon:        "belts",
	},
	SubcategoryHatsCaps: {
		Name:        SubcategoryHatsCaps,
		DisplayName: "Hats & Caps",
		Description: "Sun hats, fedoras, caps, beanies",
		Category:    CategoryAccessories,
		Icon:        "hats",
	},
	SubcategoryScarvesWraps: {
		Name:        SubcategoryScarvesWraps,
		DisplayName: "Scarves & Wraps",
		Description: "Winter scarves, summer wraps, pashminas",
		Category:    CategoryAccessories,
		Icon:        "scarves",
	},
	SubcategoryWatches: {
		Name:        SubcategoryWatches,
		DisplayName: "Watches",
		Description: "Smart watches, luxury, sport, casual",
		Category:    CategoryAccessories,
		Icon:        "watches",
	},
}

// AllProductSubcategories returns all valid product/service subcategory names
func AllProductSubcategories() []string {
	subcategories := make([]string, 0, len(ProductSubcategories))
	for key := range ProductSubcategories {
		subcategories = append(subcategories, key)
	}
	return subcategories
}

// IsValidProductSubcategory checks if a product/service subcategory is valid
func IsValidProductSubcategory(subcategory string) bool {
	_, exists := ProductSubcategories[subcategory]
	return exists
}

// GetProductSubcategoryInfo returns the info for a product/service subcategory
func GetProductSubcategoryInfo(subcategory string) (ProductSubcategoryInfo, bool) {
	info, exists := ProductSubcategories[subcategory]
	return info, exists
}

// GetProductSubcategoriesByCategory returns subcategories for a given category
func GetProductSubcategoriesByCategory(category string) []string {
	var result []string
	for key, info := range ProductSubcategories {
		if info.Category == category {
			result = append(result, key)
		}
	}
	return result
}