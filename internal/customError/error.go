package customError

import "17live_wso_be/util"

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(errorCode int) CustomError {
	return CustomError{
		Code:    errorCode,
		Message: errorMessage[errorCode],
	}
}

func (e CustomError) Error() string {
	return string(util.MapToJson(e))
}

const (
	FetchAccessTokenFail      int = 455
	FetchCampaignFail         int = 456
	FetchLeaderboardFail      int = 457
	FetchStreamerFail         int = 458
	FetchStreamerContractFail int = 459

	UserGoogleAuthFail     int = 460
	UserEmailDomainInvalid int = 461
	UserNotExist           int = 462
	UserNotActive          int = 463
	UserDuplicated         int = 464
	RegionNotExist         int = 465

	CampaignNotExist           int = 470
	CampaignBonusApproved      int = 471
	CampaignNoRegion           int = 472
	CampaignRegionSet          int = 473
	CampaignRegionNotInList    int = 474
	CampaignBonusNotExist      int = 475
	CampaignFixedBonusNotExist int = 476
	CampaignBonusPaid          int = 477
	LeaderboardNotInCampaign   int = 478
	RankNotInLeaderboard       int = 479

	PayoutNotExist               int = 480
	PayoutAdjustmentNonDeletable int = 481
	PayoutGrouped                int = 482

	MailDeliveryFailed  int = 991
	RecordNotFound      int = 992
	InvalidRequestQuery int = 993
	InvalidRequestId    int = 994
	InvalidRequestData  int = 995
	InvalidUserToken    int = 996
	PermissionDenied    int = 997
	DatabaseError       int = 998
	UnknownError        int = 999
)

var errorMessage = map[int]string{
	FetchAccessTokenFail:      "fail to fetch access token from 17 media",
	FetchCampaignFail:         "fail to fetch campaign data from 17 media",
	FetchLeaderboardFail:      "fail to fetch leaderboard data from 17 media",
	FetchStreamerFail:         "fail to fetch streamer data from 17 media",
	FetchStreamerContractFail: "fail to fetch streamer contract data from 17 media",

	UserGoogleAuthFail:     "user google auth fail",
	UserEmailDomainInvalid: "user email doamin not in admitted list",
	UserNotExist:           "user does not exist",
	UserNotActive:          "user is not active",
	UserDuplicated:         "user duplicated",
	RegionNotExist:         "region does not exist",

	CampaignNotExist:           "campaign does not exist",
	CampaignBonusApproved:      "campaign bonus has approved",
	CampaignNoRegion:           "campaign has no region",
	CampaignRegionSet:          "campaign region has set",
	CampaignRegionNotInList:    "campaign region does not in assigned region list",
	CampaignBonusNotExist:      "campaign bonus does not exist",
	CampaignFixedBonusNotExist: "campaign does not have fixed bonus record",
	CampaignBonusPaid:          "campaign bonus has paid",
	LeaderboardNotInCampaign:   "leaderboard does not in the campaign",
	RankNotInLeaderboard:       "rank does not in the leaderboard",

	PayoutNotExist:               "payout does not exist",
	PayoutAdjustmentNonDeletable: "payout adjustment cannot be deleted",
	PayoutGrouped:                "payout has grouped",

	MailDeliveryFailed:  "fail to deliver notification email",
	RecordNotFound:      "record not found",
	InvalidRequestQuery: "invalid request query",
	InvalidRequestId:    "invalid request id",
	InvalidRequestData:  "invalid request data",
	InvalidUserToken:    "invalid user token",
	PermissionDenied:    "permission denied",
	DatabaseError:       "database operation fail",
	UnknownError:        "unknown error",
}
