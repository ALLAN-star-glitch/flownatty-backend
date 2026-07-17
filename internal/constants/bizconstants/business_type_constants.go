// internal/constants/bizconstants/business__type_constants.go
package bizconstants

// ================================================
// BUSINESS TYPE CONSTANTS (Legal Structures)
// ================================================

const (
	// Individual/Personal Businesses
	BusinessTypeSoleProprietorship = "sole_proprietorship"
	BusinessTypeSoleTrader         = "sole_trader"
	
	// Partnership
	BusinessTypePartnership        = "partnership"
	BusinessTypeLimitedPartnership = "limited_partnership"
	
	// Companies
	BusinessTypePrivateCompany     = "private_company"      // Ltd
	BusinessTypePublicCompany      = "public_company"       // PLC
	BusinessTypeLimitedByGuarantee = "limited_by_guarantee" // Ltd by Guarantee
	
	// Cooperatives & Community
	BusinessTypeCooperative        = "cooperative"
	BusinessTypeSACCO              = "sacco"                // Savings and Credit Cooperative
	
	// Non-Profit
	BusinessTypeNGO                = "ngo"
	BusinessTypeCBO                = "cbo"                  // Community Based Organization
	BusinessTypeTrust              = "trust"
	BusinessTypeFoundation         = "foundation"
	BusinessTypeFaithBased         = "faith_based"
	
	// Special
	BusinessTypeFranchise          = "franchise"
	BusinessTypeEPZ                = "epz"                  // Export Processing Zone
	BusinessTypeSpecialEconomic    = "special_economic"     // Special Economic Zone
	BusinessTypeStateCorporation   = "state_corporation"    // Parastatal
	BusinessTypeGovernmentAgency   = "government_agency"
)

// AllBusinessTypes defines all valid business types
var AllBusinessTypes = []string{
	// Individual
	BusinessTypeSoleProprietorship,
	BusinessTypeSoleTrader,
	
	// Partnership
	BusinessTypePartnership,
	BusinessTypeLimitedPartnership,
	
	// Companies
	BusinessTypePrivateCompany,
	BusinessTypePublicCompany,
	BusinessTypeLimitedByGuarantee,
	
	// Cooperatives
	BusinessTypeCooperative,
	BusinessTypeSACCO,
	
	// Non-Profit
	BusinessTypeNGO,
	BusinessTypeCBO,
	BusinessTypeTrust,
	BusinessTypeFoundation,
	BusinessTypeFaithBased,
	
	// Special
	BusinessTypeFranchise,
	BusinessTypeEPZ,
	BusinessTypeSpecialEconomic,
	BusinessTypeStateCorporation,
	BusinessTypeGovernmentAgency,
}

// BusinessTypeDisplayNames returns display names for business types
var BusinessTypeDisplayNames = map[string]string{
	// Individual
	BusinessTypeSoleProprietorship: "Sole Proprietorship",
	BusinessTypeSoleTrader:         "Sole Trader",
	
	// Partnership
	BusinessTypePartnership:        "Partnership",
	BusinessTypeLimitedPartnership: "Limited Partnership",
	
	// Companies
	BusinessTypePrivateCompany:     "Private Limited Company (Ltd)",
	BusinessTypePublicCompany:      "Public Limited Company (PLC)",
	BusinessTypeLimitedByGuarantee: "Company Limited by Guarantee",
	
	// Cooperatives
	BusinessTypeCooperative:        "Cooperative Society",
	BusinessTypeSACCO:              "SACCO (Savings & Credit)",
	
	// Non-Profit
	BusinessTypeNGO:                "Non-Governmental Organization (NGO)",
	BusinessTypeCBO:                "Community Based Organization (CBO)",
	BusinessTypeTrust:              "Trust",
	BusinessTypeFoundation:         "Foundation",
	BusinessTypeFaithBased:         "Faith-Based Organization",
	
	// Special
	BusinessTypeFranchise:          "Franchise",
	BusinessTypeEPZ:                "Export Processing Zone (EPZ)",
	BusinessTypeSpecialEconomic:    "Special Economic Zone",
	BusinessTypeStateCorporation:   "State Corporation (Parastatal)",
	BusinessTypeGovernmentAgency:   "Government Agency",
}

// BusinessTypeDescriptions returns descriptions for business types
var BusinessTypeDescriptions = map[string]string{
	BusinessTypeSoleProprietorship: "Single owner, unlimited liability. Common for small businesses (kiosks, mama mboga, tailors)",
	BusinessTypeSoleTrader:         "Unregistered individual trading under own name. Common for hawkers, artisans",
	BusinessTypePartnership:        "2-20 partners sharing profits and liabilities. Common for law firms, clinics, accounting firms",
	BusinessTypeLimitedPartnership: "Partners with limited liability. Common for investment partnerships",
	BusinessTypePrivateCompany:     "Limited liability, 1-50 shareholders. Most common for medium-large businesses",
	BusinessTypePublicCompany:      "Listed on stock exchange, unlimited shareholders. Large corporations (Safaricom, KCB)",
	BusinessTypeLimitedByGuarantee: "No share capital, members guarantee debts. Common for non-profits, clubs",
	BusinessTypeCooperative:        "Member-owned, democratic control. Common for farmer cooperatives, housing",
	BusinessTypeSACCO:              "Savings and Credit Cooperative. Common for employee savings (Mwalimu SACCO)",
	BusinessTypeNGO:                "Non-profit working on social issues. Common for charities, advocacy (Red Cross, AMREF)",
	BusinessTypeCBO:                "Community-based organization. Common for self-help groups, community projects",
	BusinessTypeTrust:              "Assets held for beneficiaries. Common for family trusts, community trusts",
	BusinessTypeFoundation:         "Non-profit foundation. Common for philanthropic organizations",
	BusinessTypeFaithBased:         "Religious organization. Common for churches, mosques, temples",
	BusinessTypeFranchise:          "Licensed operation of established brand. Common for fast food (KFC, Chicken Inn)",
	BusinessTypeEPZ:                "Business operating in Export Processing Zone. Common for textile manufacturing",
	BusinessTypeSpecialEconomic:    "Business in Special Economic Zone. Common for manufacturing, logistics",
	BusinessTypeStateCorporation:   "State-owned corporation. Common for utilities, infrastructure (KPLC, KPA)",
	BusinessTypeGovernmentAgency:   "Government department or agency. Common for regulatory bodies (KRA, KEBS)",
}

// BusinessTypeIcons returns icons for business types
var BusinessTypeIcons = map[string]string{
	// Individual
	BusinessTypeSoleProprietorship: "person",
	BusinessTypeSoleTrader:         "trader",
	
	// Partnership
	BusinessTypePartnership:        "people",
	BusinessTypeLimitedPartnership: "partnership",
	
	// Companies
	BusinessTypePrivateCompany:     "business",
	BusinessTypePublicCompany:      "public",
	BusinessTypeLimitedByGuarantee: "guarantee",
	
	// Cooperatives
	BusinessTypeCooperative:        "cooperative",
	BusinessTypeSACCO:              "sacco",
	
	// Non-Profit
	BusinessTypeNGO:                "volunteer",
	BusinessTypeCBO:                "community",
	BusinessTypeTrust:              "trust",
	BusinessTypeFoundation:         "foundation",
	BusinessTypeFaithBased:         "faith",
	
	// Special
	BusinessTypeFranchise:          "franchise",
	BusinessTypeEPZ:                "epz",
	BusinessTypeSpecialEconomic:    "special_economic",
	BusinessTypeStateCorporation:   "state",
	BusinessTypeGovernmentAgency:   "government",
}