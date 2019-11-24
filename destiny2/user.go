package destiny2

import (
	"path"
	"time"
)

type UserService struct {
	c *Client
}

// UserMembershipData ...
// https://bungie-net.github.io/multi/schema_User-UserMembershipData.html#schema_User-UserMembershipData
type UserMembershipData struct {
	DestinyMemberships []GroupUserInfoCard
	BungieNetUser      GeneralUser
}

// GeneralUser ...
// https://bungie-net.github.io/multi/schema_User-GeneralUser.html#schema_User-GeneralUser
type GeneralUser struct {
	MembershipID           int64             `json:"membershipId,string"`
	UniqueName             string            `json:"uniqueName"`
	NormalizedName         string            `json:"normalizedName"`
	DisplayName            string            `json:"displayName"`
	ProfilePicture         int               `json:"profilePicture"`
	ProfileTheme           int               `json:"profileTheme"`
	UserTitle              int               `json:"userTitle"`
	SuccessMessageFlags    int64             `json:"successMessageFlags,string"`
	IsDeleted              bool              `json:"isDeleted"`
	About                  string            `json:"about"`
	FirstAccess            *time.Time        `json:"firstAccess"`
	LastUpdate             *time.Time        `json:"lastUpdate"`
	LegacyPortalUID        *int64            `json:"legacyPortalUID"`
	Context                UserToUserContext `json:"context"`
	PSNDisplayName         *string           `json:"psnDisplayName"`
	XboxDisplayName        *string           `json:"xboxDisplayName"`
	FBDisplayName          *string           `json:"fbDisplayName"`
	BlizzardDisplayName    *string           `json:"blizzardDisplayName"`
	SteamDisplayName       *string           `json:"steamDisplayName"`
	StadiaDisplayName      *string           `json:"stadiaDisplayName"`
	ShowActivity           *bool             `json:"showActivity"`
	Locale                 string            `json:"locale"`
	LocaleInheritDefault   bool              `json:"localeInheritDefault"`
	LastBanReportID        *int64            `json:"lastBanReportId"`
	ShowGroupMessaging     bool              `json:"showGroupMessaging"`
	ProfilePicturePath     string            `json:"profilePicturePath"`
	ProfilePictureWidePath string            `json:"profilePictureWidePath"`
	ProfileThemeName       string            `json:"profileThemeName"`
	UserTitleDisplay       string            `json:"userTitleDisplay"`
	StatusText             string            `json:"statusText"`
	StatusDate             time.Time         `json:"statusDate"`
	ProfileBanExpire       *time.Time        `json:"profileBanExpire"`
}

// UserToUserContext ...
// https://bungie-net.github.io/multi/schema_User-UserToUserContext.html#schema_User-UserToUserContext
type UserToUserContext struct {
	IsFollowing         bool       `json:"isFollowing"`
	GlobalIgnoreEndDate *time.Time `json:"globalIgnoreEndDate"`
}

// IgnoreResponse ...
// https://bungie-net.github.io/multi/schema_Ignores-IgnoreResponse.html#schema_Ignores-IgnoreResponse
type IgnoreResponse struct {
	IsIgnored   bool `json:"isIgnored"`
	IgnoreFlags int  `json:"ignoreFlags"`
}

// GetMembershipDataForCurrentUser returns a list of accounts associated with signed in user.
func (us *UserService) GetMembershipDataForCurrentUser(opts ...RequestOption) (UserMembershipData, error) {
	r := UserMembershipData{}
	err := us.do("GET", "/GetMembershipsForCurrentUser", &r, opts...)
	return r, err
}

func (us *UserService) do(method, endpoint string, dst interface{}, opts ...RequestOption) error {
	endpoint = path.Join("/User", endpoint)
	return us.c.do(method, endpoint, dst, opts...)
}
