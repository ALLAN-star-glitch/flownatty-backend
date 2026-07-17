// internal/constants/productconstants/product_category_constants.go
package productconstants

// ================================================
// PRODUCT/SERVICE CATEGORY CONSTANTS
// ================================================

const (
	// Fashion & Apparel
	CategoryWomenClothing    = "women_clothing"
	CategoryMenClothing      = "men_clothing"
	CategoryChildrenClothing = "children_clothing"
	CategoryShoesFootwear    = "shoes_footwear"
	CategoryAccessories      = "accessories"
	CategoryTraditionalWear  = "traditional_wear"

	// Beauty & Personal Care
	CategoryMakeupCosmetics  = "makeup_cosmetics"
	CategorySkincare         = "skincare"
	CategoryHaircare         = "haircare"
	CategoryPerfumes         = "perfumes"
	CategoryNailCare         = "nail_care"
	CategoryMenGrooming      = "men_grooming"

	// Electronics & Technology
	CategorySmartphones      = "smartphones"
	CategoryComputers        = "computers"
	CategoryAudio            = "audio"
	CategoryTVEntertainment  = "tv_entertainment"
	CategoryGaming           = "gaming"
	CategoryCameraPhoto      = "camera_photo"

	// Home & Living
	CategoryFurniture        = "furniture"
	CategoryHomeDecor        = "home_decor"
	CategoryKitchenware      = "kitchenware"
	CategoryBeddingLinens    = "bedding_linens"
	CategoryLighting         = "lighting"

	// Food & Groceries
	CategoryFreshProduce     = "fresh_produce"
	CategoryMeatPoultry      = "meat_poultry"
	CategoryDairyBakery      = "dairy_bakery"
	CategoryPantryStaples    = "pantry_staples"
	CategoryBeverages        = "beverages"
	CategorySnacksSweets     = "snacks_sweets"

	// Baby & Kids
	CategoryBabyProducts     = "baby_products"
	CategoryToysGames        = "toys_games"
	CategorySchoolSupplies   = "school_supplies"

	// Health & Wellness
	CategoryPharmaceuticals  = "pharmaceuticals"
	CategoryVitaminsSupps    = "vitamins_supplements"
	CategoryFitnessEquip     = "fitness_equipment"

	// Automotive
	CategoryCarParts         = "car_parts"
	CategoryMotorcycleParts  = "motorcycle_parts"
	CategoryCarAccessories   = "car_accessories"

	// Services - Beauty
	CategoryHairServices     = "hair_services"
	CategoryBarberServices   = "barber_services"
	CategorySpaMassage       = "spa_massage"
	CategoryNailServices     = "nail_services"
	CategoryMakeupServices   = "makeup_services"

	// Services - Health
	CategoryMedicalConsults  = "medical_consultations"
	CategoryDentalServices   = "dental_services"
	CategoryPhysiotherapy    = "physiotherapy"
	CategoryMentalHealth     = "mental_health"
	CategoryNutrition        = "nutrition"

	// Services - Fitness
	CategoryPersonalTraining = "personal_training"
	CategoryYogaPilates      = "yoga_pilates"
	CategorySwimmingLessons  = "swimming_lessons"
	CategoryDanceClasses     = "dance_classes"

	// Services - Home
	CategoryPlumbing         = "plumbing"
	CategoryElectrical       = "electrical"
	CategoryCleaning         = "cleaning"
	CategoryPaintingDecor    = "painting_decor"
	CategoryCarpentry        = "carpentry"

	// Services - Professional
	CategoryLegal            = "legal"
	CategoryAccounting       = "accounting"
	CategoryITServices       = "it_services"
	CategoryMarketing        = "marketing"

	// Services - Events
	CategoryEventPlanning    = "event_planning"
	CategoryCatering         = "catering"
	CategoryPhotography      = "photography"
	CategoryVideography      = "videography"

	// Services - Education
	CategoryTutoring         = "tutoring"
	CategoryMusicLessons     = "music_lessons"
	CategoryDrivingLessons   = "driving_lessons"
	CategoryCodingClasses    = "coding_classes"
	CategoryLanguageClasses  = "language_classes"

	// Other Services
	CategoryCourierDelivery  = "courier_delivery"
	CategoryTranslation      = "translation"
)

// ProductCategoryInfo holds all information about a product/service category
type ProductCategoryInfo struct {
	Name        string
	DisplayName string
	Description string
	Type        string // "product" or "service"
	Icon        string
}

// ProductCategories is the single source of truth for all product/service categories
var ProductCategories = map[string]ProductCategoryInfo{
	// ========================================
	// FASHION & APPAREL (Products)
	// ========================================
	CategoryWomenClothing: {
		Name:        CategoryWomenClothing,
		DisplayName: "Women's Clothing",
		Description: "Dresses, tops, skirts, traditional wear (kangas, kitenge)",
		Type:        "product",
		Icon:        "women",
	},
	CategoryMenClothing: {
		Name:        CategoryMenClothing,
		DisplayName: "Men's Clothing",
		Description: "Shirts, suits, casual wear, traditional attire",
		Type:        "product",
		Icon:        "men",
	},
	CategoryChildrenClothing: {
		Name:        CategoryChildrenClothing,
		DisplayName: "Children's Clothing",
		Description: "Kids' clothes, school uniforms, baby wear",
		Type:        "product",
		Icon:        "children",
	},
	CategoryShoesFootwear: {
		Name:        CategoryShoesFootwear,
		DisplayName: "Shoes & Footwear",
		Description: "All types of footwear including traditional sandals",
		Type:        "product",
		Icon:        "shoes",
	},
	CategoryAccessories: {
		Name:        CategoryAccessories,
		DisplayName: "Accessories",
		Description: "Bags, belts, hats, scarves, jewelry, watches",
		Type:        "product",
		Icon:        "accessories",
	},
	CategoryTraditionalWear: {
		Name:        CategoryTraditionalWear,
		DisplayName: "Traditional Wear",
		Description: "Kitenge, kanga, khanga, and cultural attire",
		Type:        "product",
		Icon:        "traditional",
	},

	// ========================================
	// BEAUTY & PERSONAL CARE (Products)
	// ========================================
	CategoryMakeupCosmetics: {
		Name:        CategoryMakeupCosmetics,
		DisplayName: "Makeup & Cosmetics",
		Description: "Foundation, lipstick, mascara, eyeliner, beauty products",
		Type:        "product",
		Icon:        "makeup",
	},
	CategorySkincare: {
		Name:        CategorySkincare,
		DisplayName: "Skincare",
		Description: "Moisturizers, serums, sunscreens, cleansers, toners",
		Type:        "product",
		Icon:        "skincare",
	},
	CategoryHaircare: {
		Name:        CategoryHaircare,
		DisplayName: "Haircare",
		Description: "Shampoo, conditioners, hair oils, styling products",
		Type:        "product",
		Icon:        "haircare",
	},
	CategoryPerfumes: {
		Name:        CategoryPerfumes,
		DisplayName: "Perfumes & Fragrances",
		Description: "Perfumes, body sprays, essential oils",
		Type:        "product",
		Icon:        "perfume",
	},
	CategoryNailCare: {
		Name:        CategoryNailCare,
		DisplayName: "Nail Care",
		Description: "Nail polish, nail art, treatments",
		Type:        "product",
		Icon:        "nails",
	},
	CategoryMenGrooming: {
		Name:        CategoryMenGrooming,
		DisplayName: "Men's Grooming",
		Description: "Shaving products, beard care, cologne",
		Type:        "product",
		Icon:        "grooming",
	},

	// ========================================
	// ELECTRONICS & TECHNOLOGY (Products)
	// ========================================
	CategorySmartphones: {
		Name:        CategorySmartphones,
		DisplayName: "Smartphones & Accessories",
		Description: "Phones, cases, screen protectors, chargers, power banks",
		Type:        "product",
		Icon:        "phone",
	},
	CategoryComputers: {
		Name:        CategoryComputers,
		DisplayName: "Laptops & Computers",
		Description: "Laptops, desktops, tablets, computer accessories",
		Type:        "product",
		Icon:        "computer",
	},
	CategoryAudio: {
		Name:        CategoryAudio,
		DisplayName: "Audio & Headphones",
		Description: "Earphones, headphones, speakers, sound systems",
		Type:        "product",
		Icon:        "audio",
	},
	CategoryTVEntertainment: {
		Name:        CategoryTVEntertainment,
		DisplayName: "TV & Home Entertainment",
		Description: "Televisions, projectors, streaming devices, sound bars",
		Type:        "product",
		Icon:        "tv",
	},
	CategoryGaming: {
		Name:        CategoryGaming,
		DisplayName: "Gaming",
		Description: "Consoles, games, controllers, gaming accessories",
		Type:        "product",
		Icon:        "gaming",
	},
	CategoryCameraPhoto: {
		Name:        CategoryCameraPhoto,
		DisplayName: "Camera & Photography",
		Description: "Cameras, lenses, photography accessories",
		Type:        "product",
		Icon:        "camera",
	},

	// ========================================
	// HOME & LIVING (Products)
	// ========================================
	CategoryFurniture: {
		Name:        CategoryFurniture,
		DisplayName: "Furniture",
		Description: "Beds, sofas, tables, chairs, wardrobes, storage",
		Type:        "product",
		Icon:        "furniture",
	},
	CategoryHomeDecor: {
		Name:        CategoryHomeDecor,
		DisplayName: "Home Decor",
		Description: "Curtains, rugs, wall art, vases, mirrors, ornaments",
		Type:        "product",
		Icon:        "decor",
	},
	CategoryKitchenware: {
		Name:        CategoryKitchenware,
		DisplayName: "Kitchenware",
		Description: "Pots, pans, utensils, cutlery, cookware, bakeware",
		Type:        "product",
		Icon:        "kitchen",
	},
	CategoryBeddingLinens: {
		Name:        CategoryBeddingLinens,
		DisplayName: "Bedding & Linens",
		Description: "Sheets, pillows, blankets, towels, duvets",
		Type:        "product",
		Icon:        "bedding",
	},
	CategoryLighting: {
		Name:        CategoryLighting,
		DisplayName: "Lighting",
		Description: "Lamps, ceiling lights, bulbs, decorative lighting",
		Type:        "product",
		Icon:        "lighting",
	},

	// ========================================
	// FOOD & GROCERIES (Products)
	// ========================================
	CategoryFreshProduce: {
		Name:        CategoryFreshProduce,
		DisplayName: "Fresh Produce",
		Description: "Vegetables, fruits, herbs, fresh farm produce",
		Type:        "product",
		Icon:        "produce",
	},
	CategoryMeatPoultry: {
		Name:        CategoryMeatPoultry,
		DisplayName: "Meat & Poultry",
		Description: "Beef, chicken, goat, fish, sausages, nyama",
		Type:        "product",
		Icon:        "meat",
	},
	CategoryDairyBakery: {
		Name:        CategoryDairyBakery,
		DisplayName: "Dairy & Bakery",
		Description: "Milk, yogurt, cheese, bread, cakes, pastries",
		Type:        "product",
		Icon:        "dairy",
	},
	CategoryPantryStaples: {
		Name:        CategoryPantryStaples,
		DisplayName: "Pantry Staples",
		Description: "Rice, flour, sugar, cooking oil, salt, spices",
		Type:        "product",
		Icon:        "pantry",
	},
	CategoryBeverages: {
		Name:        CategoryBeverages,
		DisplayName: "Beverages",
		Description: "Sodas, juices, water, energy drinks, tea, coffee",
		Type:        "product",
		Icon:        "beverage",
	},
	CategorySnacksSweets: {
		Name:        CategorySnacksSweets,
		DisplayName: "Snacks & Sweets",
		Description: "Chips, biscuits, sweets, nuts, popcorn",
		Type:        "product",
		Icon:        "snacks",
	},

	// ========================================
	// BABY & KIDS (Products)
	// ========================================
	CategoryBabyProducts: {
		Name:        CategoryBabyProducts,
		DisplayName: "Baby Products",
		Description: "Diapers, baby food, baby gear, strollers, car seats",
		Type:        "product",
		Icon:        "baby",
	},
	CategoryToysGames: {
		Name:        CategoryToysGames,
		DisplayName: "Toys & Games",
		Description: "Educational toys, dolls, cars, puzzles, games",
		Type:        "product",
		Icon:        "toys",
	},
	CategorySchoolSupplies: {
		Name:        CategorySchoolSupplies,
		DisplayName: "School Supplies",
		Description: "Backpacks, notebooks, stationery, art supplies",
		Type:        "product",
		Icon:        "school_supplies",
	},

	// ========================================
	// HEALTH & WELLNESS (Products)
	// ========================================
	CategoryPharmaceuticals: {
		Name:        CategoryPharmaceuticals,
		DisplayName: "Pharmaceuticals",
		Description: "Prescription and OTC medicine, health products",
		Type:        "product",
		Icon:        "medicine",
	},
	CategoryVitaminsSupps: {
		Name:        CategoryVitaminsSupps,
		DisplayName: "Vitamins & Supplements",
		Description: "Vitamins, minerals, protein, herbal supplements",
		Type:        "product",
		Icon:        "supplements",
	},
	CategoryFitnessEquip: {
		Name:        CategoryFitnessEquip,
		DisplayName: "Fitness Equipment",
		Description: "Weights, bands, yoga mats, gym equipment",
		Type:        "product",
		Icon:        "fitness",
	},

	// ========================================
	// AUTOMOTIVE (Products)
	// ========================================
	CategoryCarParts: {
		Name:        CategoryCarParts,
		DisplayName: "Car Parts & Accessories",
		Description: "Spare parts, seat covers, mats, car electronics",
		Type:        "product",
		Icon:        "car_parts",
	},
	CategoryMotorcycleParts: {
		Name:        CategoryMotorcycleParts,
		DisplayName: "Motorcycle Parts",
		Description: "Motorcycle spare parts, helmets, accessories",
		Type:        "product",
		Icon:        "motorcycle",
	},
	CategoryCarAccessories: {
		Name:        CategoryCarAccessories,
		DisplayName: "Car Accessories",
		Description: "Seat covers, steering covers, car electronics",
		Type:        "product",
		Icon:        "car_accessories",
	},

	// ========================================
	// SERVICES - BEAUTY
	// ========================================
	CategoryHairServices: {
		Name:        CategoryHairServices,
		DisplayName: "Hair Services",
		Description: "Haircuts, styling, weaves, braids, locs, extensions",
		Type:        "service",
		Icon:        "hair",
	},
	CategoryBarberServices: {
		Name:        CategoryBarberServices,
		DisplayName: "Barber Services",
		Description: "Men's haircuts, shaves, grooming, beard care",
		Type:        "service",
		Icon:        "barber_service",
	},
	CategorySpaMassage: {
		Name:        CategorySpaMassage,
		DisplayName: "Spa & Massage",
		Description: "Massages, facials, body treatments, wellness services",
		Type:        "service",
		Icon:        "spa_service",
	},
	CategoryNailServices: {
		Name:        CategoryNailServices,
		DisplayName: "Nail Services",
		Description: "Manicure, pedicure, nail art, acrylic nails",
		Type:        "service",
		Icon:        "nails_service",
	},
	CategoryMakeupServices: {
		Name:        CategoryMakeupServices,
		DisplayName: "Makeup Services",
		Description: "Bridal makeup, event makeup, professional makeup",
		Type:        "service",
		Icon:        "makeup_service",
	},

	// ========================================
	// SERVICES - HEALTH
	// ========================================
	CategoryMedicalConsults: {
		Name:        CategoryMedicalConsults,
		DisplayName: "Medical Consultations",
		Description: "General practitioner, specialist consultations",
		Type:        "service",
		Icon:        "consult",
	},
	CategoryDentalServices: {
		Name:        CategoryDentalServices,
		DisplayName: "Dental Services",
		Description: "Check-ups, fillings, braces, teeth whitening",
		Type:        "service",
		Icon:        "dental_service",
	},
	CategoryPhysiotherapy: {
		Name:        CategoryPhysiotherapy,
		DisplayName: "Physiotherapy",
		Description: "Physical therapy, sports injury, rehabilitation",
		Type:        "service",
		Icon:        "physio",
	},
	CategoryMentalHealth: {
		Name:        CategoryMentalHealth,
		DisplayName: "Mental Health",
		Description: "Counseling, therapy, psychological services",
		Type:        "service",
		Icon:        "mental",
	},
	CategoryNutrition: {
		Name:        CategoryNutrition,
		DisplayName: "Nutrition Services",
		Description: "Dietary advice, meal planning, nutrition counseling",
		Type:        "service",
		Icon:        "nutrition",
	},

	// ========================================
	// SERVICES - FITNESS
	// ========================================
	CategoryPersonalTraining: {
		Name:        CategoryPersonalTraining,
		DisplayName: "Personal Training",
		Description: "One-on-one fitness coaching, workout plans",
		Type:        "service",
		Icon:        "pt",
	},
	CategoryYogaPilates: {
		Name:        CategoryYogaPilates,
		DisplayName: "Yoga & Pilates",
		Description: "Yoga classes, pilates, meditation sessions",
		Type:        "service",
		Icon:        "yoga",
	},
	CategorySwimmingLessons: {
		Name:        CategorySwimmingLessons,
		DisplayName: "Swimming Lessons",
		Description: "Swimming classes for kids, adults, beginners",
		Type:        "service",
		Icon:        "swimming",
	},
	CategoryDanceClasses: {
		Name:        CategoryDanceClasses,
		DisplayName: "Dance Classes",
		Description: "Dance classes, choreography, performances",
		Type:        "service",
		Icon:        "dance",
	},

	// ========================================
	// SERVICES - HOME
	// ========================================
	CategoryPlumbing: {
		Name:        CategoryPlumbing,
		DisplayName: "Plumbing Services",
		Description: "Plumbing repairs, installation, drainage",
		Type:        "service",
		Icon:        "plumbing",
	},
	CategoryElectrical: {
		Name:        CategoryElectrical,
		DisplayName: "Electrical Services",
		Description: "Electrical wiring, installation, repairs",
		Type:        "service",
		Icon:        "electrical",
	},
	CategoryCleaning: {
		Name:        CategoryCleaning,
		DisplayName: "Cleaning Services",
		Description: "House cleaning, office cleaning, deep cleaning",
		Type:        "service",
		Icon:        "cleaning",
	},
	CategoryPaintingDecor: {
		Name:        CategoryPaintingDecor,
		DisplayName: "Painting & Decor",
		Description: "Interior painting, exterior painting, decoration",
		Type:        "service",
		Icon:        "painting",
	},
	CategoryCarpentry: {
		Name:        CategoryCarpentry,
		DisplayName: "Carpentry Services",
		Description: "Custom furniture, woodwork, repairs, installations",
		Type:        "service",
		Icon:        "carpentry",
	},

	// ========================================
	// SERVICES - PROFESSIONAL
	// ========================================
	CategoryLegal: {
		Name:        CategoryLegal,
		DisplayName: "Legal Services",
		Description: "Legal advice, contracts, representation",
		Type:        "service",
		Icon:        "legal",
	},
	CategoryAccounting: {
		Name:        CategoryAccounting,
		DisplayName: "Accounting Services",
		Description: "Bookkeeping, tax preparation, audit",
		Type:        "service",
		Icon:        "accounting",
	},
	CategoryITServices: {
		Name:        CategoryITServices,
		DisplayName: "IT Services",
		Description: "Software development, websites, app development",
		Type:        "service",
		Icon:        "it_services",
	},
	CategoryMarketing: {
		Name:        CategoryMarketing,
		DisplayName: "Marketing Services",
		Description: "Digital marketing, SEO, advertising, branding",
		Type:        "service",
		Icon:        "marketing",
	},

	// ========================================
	// SERVICES - EVENTS
	// ========================================
	CategoryEventPlanning: {
		Name:        CategoryEventPlanning,
		DisplayName: "Event Planning",
		Description: "Weddings, corporate events, parties, planning",
		Type:        "service",
		Icon:        "event_planning",
	},
	CategoryCatering: {
		Name:        CategoryCatering,
		DisplayName: "Catering Services",
		Description: "Event catering, office catering, food service",
		Type:        "service",
		Icon:        "catering",
	},
	CategoryPhotography: {
		Name:        CategoryPhotography,
		DisplayName: "Photography Services",
		Description: "Events, portraits, studio, product photography",
		Type:        "service",
		Icon:        "photo_service",
	},
	CategoryVideography: {
		Name:        CategoryVideography,
		DisplayName: "Videography Services",
		Description: "Event videos, corporate videos, production",
		Type:        "service",
		Icon:        "video",
	},

	// ========================================
	// SERVICES - EDUCATION
	// ========================================
	CategoryTutoring: {
		Name:        CategoryTutoring,
		DisplayName: "Tutoring Services",
		Description: "Math, English, science, languages tutoring",
		Type:        "service",
		Icon:        "tutoring",
	},
	CategoryMusicLessons: {
		Name:        CategoryMusicLessons,
		DisplayName: "Music Lessons",
		Description: "Piano, guitar, vocals, drums, music theory",
		Type:        "service",
		Icon:        "music_lessons",
	},
	CategoryDrivingLessons: {
		Name:        CategoryDrivingLessons,
		DisplayName: "Driving Lessons",
		Description: "Manual, automatic, defensive driving lessons",
		Type:        "service",
		Icon:        "driving",
	},
	CategoryCodingClasses: {
		Name:        CategoryCodingClasses,
		DisplayName: "Coding Classes",
		Description: "Programming, web development, app development",
		Type:        "service",
		Icon:        "coding",
	},
	CategoryLanguageClasses: {
		Name:        CategoryLanguageClasses,
		DisplayName: "Language Classes",
		Description: "Language learning, translation, interpretation",
		Type:        "service",
		Icon:        "language",
	},

	// ========================================
	// OTHER SERVICES
	// ========================================
	CategoryCourierDelivery: {
		Name:        CategoryCourierDelivery,
		DisplayName: "Courier & Delivery",
		Description: "Document delivery, package delivery, same-day",
		Type:        "service",
		Icon:        "courier",
	},
	CategoryTranslation: {
		Name:        CategoryTranslation,
		DisplayName: "Translation Services",
		Description: "Document translation, interpreting services",
		Type:        "service",
		Icon:        "translation",
	},
}

// AllProductCategories returns all valid product/service category names
func AllProductCategories() []string {
	categories := make([]string, 0, len(ProductCategories))
	for key := range ProductCategories {
		categories = append(categories, key)
	}
	return categories
}

// IsValidProductCategory checks if a product/service category is valid
func IsValidProductCategory(category string) bool {
	_, exists := ProductCategories[category]
	return exists
}

// GetProductCategoryInfo returns the info for a product/service category
func GetProductCategoryInfo(category string) (ProductCategoryInfo, bool) {
	info, exists := ProductCategories[category]
	return info, exists
}

// GetProductCategoriesByType returns categories filtered by type
func GetProductCategoriesByType(categoryType string) []string {
	var result []string
	for key, info := range ProductCategories {
		if info.Type == categoryType {
			result = append(result, key)
		}
	}
	return result
}