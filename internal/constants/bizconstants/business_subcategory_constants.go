// internal/constants/bizconstants/business_subcategory_constants.go
package bizconstants

// ================================================
// BUSINESS SUBCATEGORY CONSTANTS (Detailed)
// ================================================

const (
	// Retail
	SubcategorySupermarket       = "supermarket"
	SubcategoryMinimart          = "minimart"
	SubcategoryGeneralStore      = "general_store"
	SubcategoryWholesale         = "wholesale"
	SubcategoryHardware          = "hardware"
	SubcategoryBookstore         = "bookstore"
	SubcategoryStationery        = "stationery"
	SubcategoryPharmacy          = "pharmacy"
	SubcategoryOptical           = "optical"
	SubcategoryJewelry           = "jewelry"
	SubcategoryElectronicsStore  = "electronics_store"
	SubcategoryFurnitureStore    = "furniture_store"
	SubcategorySportsStore       = "sports_store"
	SubcategoryToyStore          = "toy_store"
	SubcategoryBabyShop          = "baby_shop"
	SubcategoryPetShop           = "pet_shop"
	SubcategoryGiftShop          = "gift_shop"
	SubcategoryFlorist           = "florist"

	// Fashion
	SubcategoryBoutique          = "boutique"
	SubcategoryTailor            = "tailor"
	SubcategoryMitumba           = "mitumba"
	SubcategoryShoeStore         = "shoe_store"
	SubcategoryAccessories       = "accessories"
	SubcategoryTraditionalWear   = "traditional_wear"
	SubcategoryUniformStore      = "uniform_store"
	SubcategoryWeddingStore      = "wedding_store"
	SubcategoryChildrenClothing  = "children_clothing"

	// Beauty
	SubcategorySalon             = "salon"
	SubcategoryBarber            = "barber"
	SubcategorySpa               = "spa"
	SubcategoryNailSalon         = "nail_salon"
	SubcategoryMakeupStudio      = "makeup_studio"
	SubcategoryBeautySupply      = "beauty_supply"

	// Food
	SubcategoryRestaurant        = "restaurant"
	SubcategoryFastFood          = "fast_food"
	SubcategoryCafe              = "cafe"
	SubcategoryBakery            = "bakery"
	SubcategoryButcher           = "butcher"
	SubcategoryGrocery           = "grocery"
	SubcategoryJuiceBar          = "juice_bar"
	SubcategoryCatering          = "catering"
	SubcategoryFoodTruck         = "food_truck"

	// Health
	SubcategoryClinic            = "clinic"
	SubcategoryDental            = "dental"
	SubcategoryLaboratory        = "laboratory"
	SubcategoryMaternity         = "maternity"
	SubcategoryPhysiotherapy     = "physiotherapy"
	SubcategoryMentalHealth      = "mental_health"
	SubcategoryGym               = "gym"
	SubcategoryYogaStudio        = "yoga_studio"

	// Technology
	SubcategoryPhoneShop         = "phone_shop"
	SubcategoryComputerStore     = "computer_store"
	SubcategoryTechRepair        = "tech_repair"
	SubcategoryCCTV              = "cctv"
	SubcategorySolarStore        = "solar_store"

	// Professional Services
	SubcategoryLawFirm           = "law_firm"
	SubcategoryAccountingFirm    = "accounting_firm"
	SubcategoryConsulting        = "consulting"
	SubcategoryInsuranceAgency   = "insurance_agency"
	SubcategorySacco             = "sacco"
	SubcategoryITConsultancy     = "it_consultancy"
	SubcategoryMarketingAgency   = "marketing_agency"
	SubcategoryEventPlanner      = "event_planner"
	SubcategoryTravelAgency      = "travel_agency"
	SubcategoryTourOperator      = "tour_operator"

	// Education
	SubcategorySchool            = "school"
	SubcategoryCollege           = "college"
	SubcategoryUniversity        = "university"
	SubcategoryDrivingSchool     = "driving_school"
	SubcategoryTutoring          = "tutoring"
	SubcategoryCodingAcademy     = "coding_academy"

	// Construction
	SubcategoryRealEstateAgent   = "real_estate_agent"
	SubcategoryPropertyManager   = "property_manager"
	SubcategoryConstructionCompany = "construction_company"
	SubcategoryArchitect         = "architect"
	SubcategoryInteriorDesigner  = "interior_designer"
	SubcategoryPlumber           = "plumber"
	SubcategoryElectrician       = "electrician"
	SubcategoryCarpenter         = "carpenter"
	SubcategoryPainter           = "painter"

	// Automotive
	SubcategoryCarDealership     = "car_dealership"
	SubcategorySpareParts        = "spare_parts"
	SubcategoryGarage            = "garage"
	SubcategoryCarWash           = "car_wash"
	SubcategoryCarRental         = "car_rental"
	SubcategoryMotorcycleDealer  = "motorcycle_dealer"

	// Hospitality
	SubcategoryHotel             = "hotel"
	SubcategoryGuestHouse        = "guest_house"
	SubcategoryAirbnbHost        = "airbnb_host"
	SubcategoryResort            = "resort"
	SubcategoryCampingSite       = "camping_site"

	// Agriculture
	SubcategoryFarm              = "farm"
	SubcategoryAgriInputs        = "agri_inputs"
	SubcategoryLivestock         = "livestock"
	SubcategoryVetClinic         = "vet_clinic"
	SubcategoryProduceSupply     = "produce_supply"

	// Entertainment
	SubcategoryMusicStudio       = "music_studio"
	SubcategoryArtGallery        = "art_gallery"
	SubcategoryPhotographyStudio = "photography_studio"
	SubcategoryEventVenue        = "event_venue"
)

// AllBusinessSubcategories defines all valid business subcategories
var AllBusinessSubcategories = []string{
	// Retail
	SubcategorySupermarket,
	SubcategoryMinimart,
	SubcategoryGeneralStore,
	SubcategoryWholesale,
	SubcategoryHardware,
	SubcategoryBookstore,
	SubcategoryStationery,
	SubcategoryPharmacy,
	SubcategoryOptical,
	SubcategoryJewelry,
	SubcategoryElectronicsStore,
	SubcategoryFurnitureStore,
	SubcategorySportsStore,
	SubcategoryToyStore,
	SubcategoryBabyShop,
	SubcategoryPetShop,
	SubcategoryGiftShop,
	SubcategoryFlorist,

	// Fashion
	SubcategoryBoutique,
	SubcategoryTailor,
	SubcategoryMitumba,
	SubcategoryShoeStore,
	SubcategoryAccessories,
	SubcategoryTraditionalWear,
	SubcategoryUniformStore,
	SubcategoryWeddingStore,
	SubcategoryChildrenClothing,

	// Beauty
	SubcategorySalon,
	SubcategoryBarber,
	SubcategorySpa,
	SubcategoryNailSalon,
	SubcategoryMakeupStudio,
	SubcategoryBeautySupply,

	// Food
	SubcategoryRestaurant,
	SubcategoryFastFood,
	SubcategoryCafe,
	SubcategoryBakery,
	SubcategoryButcher,
	SubcategoryGrocery,
	SubcategoryJuiceBar,
	SubcategoryCatering,
	SubcategoryFoodTruck,

	// Health
	SubcategoryClinic,
	SubcategoryDental,
	SubcategoryLaboratory,
	SubcategoryMaternity,
	SubcategoryPhysiotherapy,
	SubcategoryMentalHealth,
	SubcategoryGym,
	SubcategoryYogaStudio,

	// Technology
	SubcategoryPhoneShop,
	SubcategoryComputerStore,
	SubcategoryTechRepair,
	SubcategoryCCTV,
	SubcategorySolarStore,

	// Professional Services
	SubcategoryLawFirm,
	SubcategoryAccountingFirm,
	SubcategoryConsulting,
	SubcategoryInsuranceAgency,
	SubcategorySacco,
	SubcategoryITConsultancy,
	SubcategoryMarketingAgency,
	SubcategoryEventPlanner,
	SubcategoryTravelAgency,
	SubcategoryTourOperator,

	// Education
	SubcategorySchool,
	SubcategoryCollege,
	SubcategoryUniversity,
	SubcategoryDrivingSchool,
	SubcategoryTutoring,
	SubcategoryCodingAcademy,

	// Construction
	SubcategoryRealEstateAgent,
	SubcategoryPropertyManager,
	SubcategoryConstructionCompany,
	SubcategoryArchitect,
	SubcategoryInteriorDesigner,
	SubcategoryPlumber,
	SubcategoryElectrician,
	SubcategoryCarpenter,
	SubcategoryPainter,

	// Automotive
	SubcategoryCarDealership,
	SubcategorySpareParts,
	SubcategoryGarage,
	SubcategoryCarWash,
	SubcategoryCarRental,
	SubcategoryMotorcycleDealer,

	// Hospitality
	SubcategoryHotel,
	SubcategoryGuestHouse,
	SubcategoryAirbnbHost,
	SubcategoryResort,
	SubcategoryCampingSite,

	// Agriculture
	SubcategoryFarm,
	SubcategoryAgriInputs,
	SubcategoryLivestock,
	SubcategoryVetClinic,
	SubcategoryProduceSupply,

	// Entertainment
	SubcategoryMusicStudio,
	SubcategoryArtGallery,
	SubcategoryPhotographyStudio,
	SubcategoryEventVenue,
}

// BusinessSubcategoryDisplayNames returns display names for business subcategories
var BusinessSubcategoryDisplayNames = map[string]string{
	// Retail
	SubcategorySupermarket:       "Supermarket",
	SubcategoryMinimart:          "Mini-Mart",
	SubcategoryGeneralStore:      "General Store (Duka)",
	SubcategoryWholesale:         "Wholesale Store",
	SubcategoryHardware:          "Hardware Store",
	SubcategoryBookstore:         "Bookstore",
	SubcategoryStationery:        "Stationery Store",
	SubcategoryPharmacy:          "Pharmacy/Chemist",
	SubcategoryOptical:           "Optical Store",
	SubcategoryJewelry:           "Jewelry Store",
	SubcategoryElectronicsStore:  "Electronics Store",
	SubcategoryFurnitureStore:    "Furniture Store",
	SubcategorySportsStore:       "Sports Store",
	SubcategoryToyStore:          "Toy Store",
	SubcategoryBabyShop:          "Baby Shop",
	SubcategoryPetShop:           "Pet Shop",
	SubcategoryGiftShop:          "Gift Shop",
	SubcategoryFlorist:           "Florist",

	// Fashion
	SubcategoryBoutique:          "Boutique",
	SubcategoryTailor:            "Tailoring Shop",
	SubcategoryMitumba:           "Mitumba Store",
	SubcategoryShoeStore:         "Shoe Store",
	SubcategoryAccessories:       "Accessories Store",
	SubcategoryTraditionalWear:   "Traditional Wear",
	SubcategoryUniformStore:      "Uniform Store",
	SubcategoryWeddingStore:      "Wedding Store",
	SubcategoryChildrenClothing:  "Children's Clothing",

	// Beauty
	SubcategorySalon:             "Salon",
	SubcategoryBarber:            "Barber Shop",
	SubcategorySpa:               "Spa",
	SubcategoryNailSalon:         "Nail Salon",
	SubcategoryMakeupStudio:      "Makeup Studio",
	SubcategoryBeautySupply:      "Beauty Supply Store",

	// Food
	SubcategoryRestaurant:        "Restaurant",
	SubcategoryFastFood:          "Fast Food",
	SubcategoryCafe:              "Cafe",
	SubcategoryBakery:            "Bakery",
	SubcategoryButcher:           "Butcher Shop",
	SubcategoryGrocery:           "Grocery (Mama Mboga)",
	SubcategoryJuiceBar:          "Juice Bar",
	SubcategoryCatering:          "Catering Service",
	SubcategoryFoodTruck:         "Food Truck",

	// Health
	SubcategoryClinic:            "Clinic",
	SubcategoryDental:            "Dental Clinic",
	SubcategoryLaboratory:        "Laboratory",
	SubcategoryMaternity:         "Maternity Clinic",
	SubcategoryPhysiotherapy:     "Physiotherapy",
	SubcategoryMentalHealth:      "Mental Health",
	SubcategoryGym:               "Gym",
	SubcategoryYogaStudio:        "Yoga Studio",

	// Technology
	SubcategoryPhoneShop:         "Phone Shop",
	SubcategoryComputerStore:     "Computer Store",
	SubcategoryTechRepair:        "Tech Repair",
	SubcategoryCCTV:              "CCTV Installation",
	SubcategorySolarStore:        "Solar Store",

	// Professional Services
	SubcategoryLawFirm:           "Law Firm",
	SubcategoryAccountingFirm:    "Accounting Firm",
	SubcategoryConsulting:        "Consulting Firm",
	SubcategoryInsuranceAgency:   "Insurance Agency",
	SubcategorySacco:             "Sacco",
	SubcategoryITConsultancy:     "IT Consultancy",
	SubcategoryMarketingAgency:   "Marketing Agency",
	SubcategoryEventPlanner:      "Event Planner",
	SubcategoryTravelAgency:      "Travel Agency",
	SubcategoryTourOperator:      "Tour Operator",

	// Education
	SubcategorySchool:            "School",
	SubcategoryCollege:           "College",
	SubcategoryUniversity:        "University",
	SubcategoryDrivingSchool:     "Driving School",
	SubcategoryTutoring:          "Tutoring Center",
	SubcategoryCodingAcademy:     "Coding Academy",

	// Construction
	SubcategoryRealEstateAgent:   "Real Estate Agent",
	SubcategoryPropertyManager:   "Property Manager",
	SubcategoryConstructionCompany: "Construction Company",
	SubcategoryArchitect:         "Architect",
	SubcategoryInteriorDesigner:  "Interior Designer",
	SubcategoryPlumber:           "Plumber",
	SubcategoryElectrician:       "Electrician",
	SubcategoryCarpenter:         "Carpenter",
	SubcategoryPainter:           "Painter",

	// Automotive
	SubcategoryCarDealership:     "Car Dealership",
	SubcategorySpareParts:        "Spare Parts Shop",
	SubcategoryGarage:            "Garage",
	SubcategoryCarWash:           "Car Wash",
	SubcategoryCarRental:         "Car Rental",
	SubcategoryMotorcycleDealer:  "Motorcycle Dealer",

	// Hospitality
	SubcategoryHotel:             "Hotel",
	SubcategoryGuestHouse:        "Guest House",
	SubcategoryAirbnbHost:        "Airbnb Host",
	SubcategoryResort:            "Resort",
	SubcategoryCampingSite:       "Camping Site",

	// Agriculture
	SubcategoryFarm:              "Farm",
	SubcategoryAgriInputs:        "Agri-Inputs",
	SubcategoryLivestock:         "Livestock",
	SubcategoryVetClinic:         "Vet Clinic",
	SubcategoryProduceSupply:     "Produce Supply",

	// Entertainment
	SubcategoryMusicStudio:       "Music Studio",
	SubcategoryArtGallery:        "Art Gallery",
	SubcategoryPhotographyStudio: "Photography Studio",
	SubcategoryEventVenue:        "Event Venue",
}

// BusinessSubcategoryDescriptions returns descriptions for business subcategories
var BusinessSubcategoryDescriptions = map[string]string{
	// Retail
	SubcategorySupermarket:       "Full grocery and household goods store",
	SubcategoryMinimart:          "Convenience store, quick shopping",
	SubcategoryGeneralStore:      "Small retail shop selling everyday items",
	SubcategoryWholesale:         "Bulk selling to retailers and businesses",
	SubcategoryHardware:          "Building materials, tools, and supplies",
	SubcategoryBookstore:         "Books, stationery, and office supplies",
	SubcategoryStationery:        "Office and school stationery supplies",
	SubcategoryPharmacy:          "Medicine, health products, and cosmetics",
	SubcategoryOptical:           "Glasses, contact lenses, and eye care",
	SubcategoryJewelry:           "Gold, silver, and custom jewelry",
	SubcategoryElectronicsStore:  "TVs, gadgets, and electronics",
	SubcategoryFurnitureStore:    "Home and office furniture",
	SubcategorySportsStore:       "Sporting goods, equipment, and apparel",
	SubcategoryToyStore:          "Toys, games, and educational items",
	SubcategoryBabyShop:          "Baby items, diapers, and prams",
	SubcategoryPetShop:           "Pet food, accessories, and supplies",
	SubcategoryGiftShop:          "Gifts, cards, and souvenirs",
	SubcategoryFlorist:           "Fresh flowers, bouquets, and arrangements",

	// Fashion
	SubcategoryBoutique:          "Curated fashion collections",
	SubcategoryTailor:            "Custom clothing and alterations",
	SubcategoryMitumba:           "Quality second-hand clothing",
	SubcategoryShoeStore:         "All types of footwear",
	SubcategoryAccessories:       "Bags, belts, hats, and scarves",
	SubcategoryTraditionalWear:   "Kitenge, kanga, khanga, and cultural attire",
	SubcategoryUniformStore:      "School, office, and medical uniforms",
	SubcategoryWeddingStore:      "Wedding dresses, suits, and decor",
	SubcategoryChildrenClothing:  "Kids' fashion and school uniforms",

	// Beauty
	SubcategorySalon:             "Hair, makeup, and beauty services",
	SubcategoryBarber:            "Men's haircuts and grooming",
	SubcategorySpa:               "Massage, facials, and wellness",
	SubcategoryNailSalon:         "Manicure, pedicure, and nail art",
	SubcategoryMakeupStudio:      "Professional makeup services",
	SubcategoryBeautySupply:      "Hair products and cosmetics",

	// Food
	SubcategoryRestaurant:        "Full-service dining",
	SubcategoryFastFood:          "Quick service and takeaway",
	SubcategoryCafe:              "Coffee, tea, and light meals",
	SubcategoryBakery:            "Bread, cakes, and pastries",
	SubcategoryButcher:           "Meat, chicken, and sausages",
	SubcategoryGrocery:           "Fresh vegetables and fruits (Mama Mboga)",
	SubcategoryJuiceBar:          "Fresh juices and smoothies",
	SubcategoryCatering:          "Event and office catering",
	SubcategoryFoodTruck:         "Mobile food vending",

	// Health
	SubcategoryClinic:            "General medical consultations",
	SubcategoryDental:            "Dentistry and oral health",
	SubcategoryLaboratory:        "Medical tests and blood work",
	SubcategoryMaternity:         "Prenatal and postnatal care",
	SubcategoryPhysiotherapy:     "Physical therapy and rehabilitation",
	SubcategoryMentalHealth:      "Counseling and therapy services",
	SubcategoryGym:               "Fitness center and training",
	SubcategoryYogaStudio:        "Yoga and meditation classes",

	// Technology
	SubcategoryPhoneShop:         "Phones and accessories",
	SubcategoryComputerStore:     "Laptops and computers",
	SubcategoryTechRepair:        "Phone, laptop, and electronics repair",
	SubcategoryCCTV:              "Security and surveillance systems",
	SubcategorySolarStore:        "Solar panels and products",

	// Professional Services
	SubcategoryLawFirm:           "Legal services and representation",
	SubcategoryAccountingFirm:    "Bookkeeping, tax, and audit",
	SubcategoryConsulting:        "Management and strategy consulting",
	SubcategoryInsuranceAgency:   "Insurance sales and claims",
	SubcategorySacco:             "Savings and loans (member-owned)",
	SubcategoryITConsultancy:     "Tech consulting and software",
	SubcategoryMarketingAgency:   "Marketing and advertising",
	SubcategoryEventPlanner:      "Wedding and corporate events",
	SubcategoryTravelAgency:      "Travel bookings and tours",
	SubcategoryTourOperator:      "Safari and tourism services",

	// Education
	SubcategorySchool:            "Primary and secondary education",
	SubcategoryCollege:           "Tertiary education",
	SubcategoryUniversity:        "Higher education",
	SubcategoryDrivingSchool:     "Driving lessons",
	SubcategoryTutoring:          "Private lessons and homework help",
	SubcategoryCodingAcademy:     "Programming and tech skills",

	// Construction
	SubcategoryRealEstateAgent:   "Property sales and rentals",
	SubcategoryPropertyManager:   "Property management",
	SubcategoryConstructionCompany: "Building and construction",
	SubcategoryArchitect:         "Building design and planning",
	SubcategoryInteriorDesigner:  "Interior decoration",
	SubcategoryPlumber:           "Plumbing services",
	SubcategoryElectrician:       "Electrical installations",
	SubcategoryCarpenter:         "Woodwork and furniture",
	SubcategoryPainter:           "Painting and decorating",

	// Automotive
	SubcategoryCarDealership:     "New and used car sales",
	SubcategorySpareParts:        "Auto parts and accessories",
	SubcategoryGarage:            "Car repairs and service",
	SubcategoryCarWash:           "Vehicle cleaning and detailing",
	SubcategoryCarRental:         "Vehicle hire services",
	SubcategoryMotorcycleDealer:  "Motorcycles and boda boda",

	// Hospitality
	SubcategoryHotel:             "Full-service accommodation",
	SubcategoryGuestHouse:        "Bed and breakfast",
	SubcategoryAirbnbHost:        "Short-term rentals",
	SubcategoryResort:            "Vacation resort",
	SubcategoryCampingSite:       "Outdoor accommodation",

	// Agriculture
	SubcategoryFarm:              "Crop and livestock farming",
	SubcategoryAgriInputs:        "Seeds, fertilizer, and pesticides",
	SubcategoryLivestock:         "Animals, feed, and supplies",
	SubcategoryVetClinic:         "Animal health services",
	SubcategoryProduceSupply:     "Fresh farm produce",

	// Entertainment
	SubcategoryMusicStudio:       "Recording and production",
	SubcategoryArtGallery:        "Art exhibitions and sales",
	SubcategoryPhotographyStudio: "Professional photography",
	SubcategoryEventVenue:        "Wedding and party venues",
}


// BusinessSubcategorySectorMap maps each subcategory to its parent sector display name
var BusinessSubcategorySectorMap = map[string]string{
	// Retail - Use the display name that matches the database
	SubcategorySupermarket:       "Retail & Consumer Goods",
	SubcategoryMinimart:          "Retail & Consumer Goods",
	SubcategoryGeneralStore:      "Retail & Consumer Goods",
	SubcategoryWholesale:         "Retail & Consumer Goods",
	SubcategoryHardware:          "Retail & Consumer Goods",
	SubcategoryBookstore:         "Retail & Consumer Goods",
	SubcategoryStationery:        "Retail & Consumer Goods",
	SubcategoryPharmacy:          "Retail & Consumer Goods",
	SubcategoryOptical:           "Retail & Consumer Goods",
	SubcategoryJewelry:           "Retail & Consumer Goods",
	SubcategoryElectronicsStore:  "Retail & Consumer Goods",
	SubcategoryFurnitureStore:    "Retail & Consumer Goods",
	SubcategorySportsStore:       "Retail & Consumer Goods",
	SubcategoryToyStore:          "Retail & Consumer Goods",
	SubcategoryBabyShop:          "Retail & Consumer Goods",
	SubcategoryPetShop:           "Retail & Consumer Goods",
	SubcategoryGiftShop:          "Retail & Consumer Goods",
	SubcategoryFlorist:           "Retail & Consumer Goods",

	// Fashion
	SubcategoryBoutique:          "Fashion & Apparel",
	SubcategoryTailor:            "Fashion & Apparel",
	SubcategoryMitumba:           "Fashion & Apparel",
	SubcategoryShoeStore:         "Fashion & Apparel",
	SubcategoryAccessories:       "Fashion & Apparel",
	SubcategoryTraditionalWear:   "Fashion & Apparel",
	SubcategoryUniformStore:      "Fashion & Apparel",
	SubcategoryWeddingStore:      "Fashion & Apparel",
	SubcategoryChildrenClothing:  "Fashion & Apparel",

	// Beauty
	SubcategorySalon:             "Beauty & Personal Care",
	SubcategoryBarber:            "Beauty & Personal Care",
	SubcategorySpa:               "Beauty & Personal Care",
	SubcategoryNailSalon:         "Beauty & Personal Care",
	SubcategoryMakeupStudio:      "Beauty & Personal Care",
	SubcategoryBeautySupply:      "Beauty & Personal Care",

	// Food
	SubcategoryRestaurant:        "Food & Beverage",
	SubcategoryFastFood:          "Food & Beverage",
	SubcategoryCafe:              "Food & Beverage",
	SubcategoryBakery:            "Food & Beverage",
	SubcategoryButcher:           "Food & Beverage",
	SubcategoryGrocery:           "Food & Beverage",
	SubcategoryJuiceBar:          "Food & Beverage",
	SubcategoryCatering:          "Food & Beverage",
	SubcategoryFoodTruck:         "Food & Beverage",

	// Health
	SubcategoryClinic:            "Health & Wellness",
	SubcategoryDental:            "Health & Wellness",
	SubcategoryLaboratory:        "Health & Wellness",
	SubcategoryMaternity:         "Health & Wellness",
	SubcategoryPhysiotherapy:     "Health & Wellness",
	SubcategoryMentalHealth:      "Health & Wellness",
	SubcategoryGym:               "Health & Wellness",
	SubcategoryYogaStudio:        "Health & Wellness",

	// Technology
	SubcategoryPhoneShop:         "Technology & Digital",
	SubcategoryComputerStore:     "Technology & Digital",
	SubcategoryTechRepair:        "Technology & Digital",
	SubcategoryCCTV:              "Technology & Digital",
	SubcategorySolarStore:        "Technology & Digital",

	// Professional Services
	SubcategoryLawFirm:           "Professional Services",
	SubcategoryAccountingFirm:    "Professional Services",
	SubcategoryConsulting:        "Professional Services",
	SubcategoryInsuranceAgency:   "Professional Services",
	SubcategorySacco:             "Professional Services",
	SubcategoryITConsultancy:     "Professional Services",
	SubcategoryMarketingAgency:   "Professional Services",
	SubcategoryEventPlanner:      "Professional Services",
	SubcategoryTravelAgency:      "Professional Services",
	SubcategoryTourOperator:      "Professional Services",

	// Education
	SubcategorySchool:            "Education & Training",
	SubcategoryCollege:           "Education & Training",
	SubcategoryUniversity:        "Education & Training",
	SubcategoryDrivingSchool:     "Education & Training",
	SubcategoryTutoring:          "Education & Training",
	SubcategoryCodingAcademy:     "Education & Training",

	// Construction
	SubcategoryRealEstateAgent:   "Construction & Building",
	SubcategoryPropertyManager:   "Construction & Building",
	SubcategoryConstructionCompany: "Construction & Building",
	SubcategoryArchitect:         "Construction & Building",
	SubcategoryInteriorDesigner:  "Construction & Building",
	SubcategoryPlumber:           "Construction & Building",
	SubcategoryElectrician:       "Construction & Building",
	SubcategoryCarpenter:         "Construction & Building",
	SubcategoryPainter:           "Construction & Building",

	// Automotive
	SubcategoryCarDealership:     "Automotive",
	SubcategorySpareParts:        "Automotive",
	SubcategoryGarage:            "Automotive",
	SubcategoryCarWash:           "Automotive",
	SubcategoryCarRental:         "Automotive",
	SubcategoryMotorcycleDealer:  "Automotive",

	// Hospitality
	SubcategoryHotel:             "Hospitality",
	SubcategoryGuestHouse:        "Hospitality",
	SubcategoryAirbnbHost:        "Hospitality",
	SubcategoryResort:            "Hospitality",
	SubcategoryCampingSite:       "Hospitality",

	// Agriculture
	SubcategoryFarm:              "Agriculture & Farming",
	SubcategoryAgriInputs:        "Agriculture & Farming",
	SubcategoryLivestock:         "Agriculture & Farming",
	SubcategoryVetClinic:         "Agriculture & Farming",
	SubcategoryProduceSupply:     "Agriculture & Farming",

	// Entertainment
	SubcategoryMusicStudio:       "Entertainment",
	SubcategoryArtGallery:        "Entertainment",
	SubcategoryPhotographyStudio: "Entertainment",
	SubcategoryEventVenue:        "Entertainment",
}

// BusinessSubcategoryIcons returns icons for business subcategories
var BusinessSubcategoryIcons = map[string]string{
	// Retail
	SubcategorySupermarket:       "supermarket",
	SubcategoryMinimart:          "minimart",
	SubcategoryGeneralStore:      "store",
	SubcategoryWholesale:         "warehouse",
	SubcategoryHardware:          "hardware",
	SubcategoryBookstore:         "book",
	SubcategoryStationery:        "stationery",
	SubcategoryPharmacy:          "pharmacy",
	SubcategoryOptical:           "glasses",
	SubcategoryJewelry:           "diamond",
	SubcategoryElectronicsStore:  "electronics",
	SubcategoryFurnitureStore:    "furniture",
	SubcategorySportsStore:       "sports",
	SubcategoryToyStore:          "toys",
	SubcategoryBabyShop:          "baby",
	SubcategoryPetShop:           "pet",
	SubcategoryGiftShop:          "gift",
	SubcategoryFlorist:           "flower",

	// Fashion
	SubcategoryBoutique:          "boutique",
	SubcategoryTailor:            "tailor",
	SubcategoryMitumba:           "clothes",
	SubcategoryShoeStore:         "shoes",
	SubcategoryAccessories:       "accessories",
	SubcategoryTraditionalWear:   "traditional",
	SubcategoryUniformStore:      "uniform",
	SubcategoryWeddingStore:      "wedding",
	SubcategoryChildrenClothing:  "children",

	// Beauty
	SubcategorySalon:             "salon",
	SubcategoryBarber:            "barber",
	SubcategorySpa:               "spa",
	SubcategoryNailSalon:         "nails",
	SubcategoryMakeupStudio:      "makeup",
	SubcategoryBeautySupply:      "beauty",

	// Food
	SubcategoryRestaurant:        "restaurant",
	SubcategoryFastFood:          "fastfood",
	SubcategoryCafe:              "cafe",
	SubcategoryBakery:            "bakery",
	SubcategoryButcher:           "butcher",
	SubcategoryGrocery:           "grocery",
	SubcategoryJuiceBar:          "juice",
	SubcategoryCatering:          "catering",
	SubcategoryFoodTruck:         "foodtruck",

	// Health
	SubcategoryClinic:            "clinic",
	SubcategoryDental:            "dental",
	SubcategoryLaboratory:        "lab",
	SubcategoryMaternity:         "maternity",
	SubcategoryPhysiotherapy:     "physio",
	SubcategoryMentalHealth:      "mental",
	SubcategoryGym:               "gym",
	SubcategoryYogaStudio:        "yoga",

	// Technology
	SubcategoryPhoneShop:         "phone",
	SubcategoryComputerStore:     "computer",
	SubcategoryTechRepair:        "repair",
	SubcategoryCCTV:              "cctv",
	SubcategorySolarStore:        "solar",

	// Professional Services
	SubcategoryLawFirm:           "law",
	SubcategoryAccountingFirm:    "accounting",
	SubcategoryConsulting:        "consulting",
	SubcategoryInsuranceAgency:   "insurance",
	SubcategorySacco:             "sacco",
	SubcategoryITConsultancy:     "it",
	SubcategoryMarketingAgency:   "marketing",
	SubcategoryEventPlanner:      "event",
	SubcategoryTravelAgency:      "travel",
	SubcategoryTourOperator:      "tour",

	// Education
	SubcategorySchool:            "school",
	SubcategoryCollege:           "college",
	SubcategoryUniversity:        "university",
	SubcategoryDrivingSchool:     "driving",
	SubcategoryTutoring:          "tutor",
	SubcategoryCodingAcademy:     "coding",

	// Construction
	SubcategoryRealEstateAgent:   "realestate",
	SubcategoryPropertyManager:   "property",
	SubcategoryConstructionCompany: "construction",
	SubcategoryArchitect:         "architect",
	SubcategoryInteriorDesigner:  "interior",
	SubcategoryPlumber:           "plumber",
	SubcategoryElectrician:       "electrician",
	SubcategoryCarpenter:         "carpenter",
	SubcategoryPainter:           "painter",

	// Automotive
	SubcategoryCarDealership:     "cardealer",
	SubcategorySpareParts:        "spareparts",
	SubcategoryGarage:            "garage",
	SubcategoryCarWash:           "carwash",
	SubcategoryCarRental:         "carrental",
	SubcategoryMotorcycleDealer:  "motorcycle",

	// Hospitality
	SubcategoryHotel:             "hotel",
	SubcategoryGuestHouse:        "guesthouse",
	SubcategoryAirbnbHost:        "airbnb",
	SubcategoryResort:            "resort",
	SubcategoryCampingSite:       "camping",

	// Agriculture
	SubcategoryFarm:              "farm",
	SubcategoryAgriInputs:        "agri",
	SubcategoryLivestock:         "livestock",
	SubcategoryVetClinic:         "vet",
	SubcategoryProduceSupply:     "produce",

	// Entertainment
	SubcategoryMusicStudio:       "music",
	SubcategoryArtGallery:        "art",
	SubcategoryPhotographyStudio: "photo",
	SubcategoryEventVenue:        "event",
}