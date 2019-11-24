package destiny2

import (
	"fmt"
	"path"
	"time"
)

// GroupV2Service is an interface for interfacing with the groupv2 endpoints
// of the Bungie API.
// https://bungie-net.github.io/multi/operation_get_GroupV2-GetAvailableAvatars.html#operation_get_GroupV2-GetAvailableAvatars
type GroupV2Service struct {
	c *Client
}

// SearchResultOfGroupMember ...
// https://bungie-net.github.io/multi/schema_SearchResultOfGroupMember.html#schema_SearchResultOfGroupMember
type SearchResultOfGroupMember struct {
	Results                      []GroupMember `json:"results"`
	TotalResults                 int           `json:"totalResults"`
	HasMore                      bool          `json:"hasMore"`
	Query                        PagedQuery    `json:"query"`
	ReplacementContinuationToken string        `json:"replacementContinuationToken"`
}

// GroupMember ...
// https://bungie-net.github.io/multi/schema_GroupsV2-GroupMember.html#schema_GroupsV2-GroupMember
type GroupMember struct {
	MemberType             int               `json:"memberType"`
	IsOnline               bool              `json:"isOnline"`
	LastOnlineStatusChange int64             `json:"lastOnlineStatusChange"`
	GroupID                int64             `json:"groudId"`
	DestinyUserInfo        GroupUserInfoCard `json:"destinyUserInfo"`
	JoinDate               time.Time         `json:"joinDate"`
}

// GroupUserInfoCard ...
// https://bungie-net.github.io/multi/schema_GroupsV2-GroupUserInfoCard.html#schema_GroupsV2-GroupUserInfoCard
type GroupUserInfoCard struct {
	LastSeenDisplayName       string  `json:"LastSeenDisplayName"`
	LastSeenDisplayNameType   int     `json:"LastSeenDisplayNameType"`
	SupplementalDisplayName   string  `json:"supplementalDisplayName"`
	IconPath                  string  `json:"iconPath"`
	CrossSaveOverride         int     `json:"crossSaveOverride"`
	ApplicableMembershipTypes []int32 `json:"applicableMembershipTypes"`
	IsPublic                  bool    `json:"isPublic"`
	MembershipType            int     `json:"membershipType"`
	MembershipID              int64   `json:"membershipId,string"`
	DisplayName               string  `json:"displayName"`
}

// GetMembersOfGroup gets a list of members in a given group
func (gs *GroupV2Service) GetMembersOfGroup(gid int64, opts ...RequestOption) (SearchResultOfGroupMember, error) {
	r := SearchResultOfGroupMember{}
	endpoint := fmt.Sprintf("/%d/Members", gid)
	err := gs.do("GET", endpoint, &r, opts...)
	return r, err
}

// GetAllMembersOfGroup gets all the members of a given group and paginates through pages if needed.
func (gs *GroupV2Service) GetAllMembersOfGroup(gid int64, opts ...RequestOption) ([]GroupMember, error) {
	members := []GroupMember{}

	// Looping until no more pages or unrecoverable error
	for page := 1; ; page++ {
		// Adding/overriding page query param
		opts = append(opts, OptionQuery("currentPage", page))

		// Getting page of memeber
		resp, err := gs.GetMembersOfGroup(gid, opts...)
		if err != nil {
			return nil, err
		}

		// Adding members to array
		members = append(members, resp.Results...)

		// Checking if there are more members to get
		if !resp.HasMore {
			break
		}
	}

	return members, nil
}

func (gs *GroupV2Service) do(method, endpoint string, dst interface{}, opts ...RequestOption) error {
	endpoint = path.Join("/GroupV2", endpoint)
	return gs.c.do(method, endpoint, dst, opts...)
}
