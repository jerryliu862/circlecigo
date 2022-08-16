package model

const (
	UserStatusInit     int = 0
	UserStatusActive   int = 1
	UserStatusInactive int = 2

	UserAuthTypeBonus  int = 1
	UserAuthTypePayout int = 2
	UserAuthTypeReport int = 3
	UserAuthTypeTax    int = 4
	UserAuthTypeSystem int = 99

	UserAuthLevelView    int = 1
	UserAuthLevelEdit    int = 2
	UserAuthLevelApprove int = 3

	RegionAll string = "ALL"

	CampaignLangAR string = "ar"
	CampaignLangEN string = "en"
	CampaignLangUS string = "en_US"
	CampaignLangJP string = "ja"
	CampaignLangCN string = "zh_CN"
	CampaignLangHK string = "zh_HK"
	CampaignLangTW string = "zh_TW"

	BonusTypeFixed          int = 0
	BonusTypeVariable       int = 1
	BonusTypeTransportation int = 2
	BonusTypeAddon          int = 3
	BonusTypeDeduction      int = 4

	PayTypeForeign string = "F"
	PayTypeLocal   string = "L"
	PayTypeAgency  string = "A"
	PayTypeMissing string = "M"

	OwnerTypeStreamer int = 1
	OwnerTypeAgency   int = 2

	AccountTypeLocal           int = 1
	AccountTypeForeign         int = 2
	AccountTypeVerifiedRevenue int = 4
	AccountTypeCompany         int = 5
	AccountTypeOffline         int = 6
)
