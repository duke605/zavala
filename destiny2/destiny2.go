package destiny2

import (
	"path"
	"time"
)

type Destiny2Service struct {
	c *Client
}

// Component ...
// https://bungie-net.github.io/multi/schema_Destiny-DestinyComponentType.html#schema_Destiny-DestinyComponentType
type Component int

const (
	// Profiles is the most basic component, only relevant when calling GetProfile. This returns basic information
	// about the profile, which is almost nothing: a list of characterIds, some information about the last time
	// you logged in, and that most sobering statistic: how long you've played.
	Profiles Component = 100

	// Characters gets summary info about each of the characters in the profile.
	Characters Component = 200
)

// DestinyProfileResponse ...
// https://bungie-net.github.io/multi/schema_Destiny-Responses-DestinyProfileResponse.html#schema_Destiny-Responses-DestinyProfileResponse
type DestinyProfileResponse struct {
	Profile    *SingleComponentResponseOfDestinyProfileComponent               `json:"profile"`
	Characters *DictionaryComponentResponseOfint64AndDestinyCharacterComponent `json:"characters"`
}

// SingleComponentResponseOfDestinyProfileComponent ...
// https://bungie-net.github.io/multi/schema_SingleComponentResponseOfDestinyProfileComponent.html#schema_SingleComponentResponseOfDestinyProfileComponent
type SingleComponentResponseOfDestinyProfileComponent struct {
	Data    DestinyProfileComponent
	Privacy int `json:"privacy"`
}

// DestinyProfileComponent ...
// https://bungie-net.github.io/multi/schema_Destiny-Entities-Profiles-DestinyProfileComponent.html#schema_Destiny-Entities-Profiles-DestinyProfileComponent
type DestinyProfileComponent struct {
	UserInfo       struct{}  `json:"userInfo"`
	DateLastPlayed time.Time `json:"dateLastPlayed"`
	VersionsOwned  int       `json:"versionsOwned"`
	CharacterIds   []int64   `json:"characterIds"`
	SeasonHashes   []uint    `json:"seasonHashes"`
}

// DictionaryComponentResponseOfint64AndDestinyCharacterComponent ...
// https://bungie-net.github.io/multi/schema_DictionaryComponentResponseOfint64AndDestinyCharacterComponent.html#schema_DictionaryComponentResponseOfint64AndDestinyCharacterComponent
type DictionaryComponentResponseOfint64AndDestinyCharacterComponent struct {
	Data    map[int64]DestinyCharacterComponent
	Privacy int `json:"privacy"`
}

// DestinyCharacterComponent ...
// https://bungie-net.github.io/multi/schema_Destiny-Entities-Characters-DestinyCharacterComponent.html#schema_Destiny-Entities-Characters-DestinyCharacterComponent
type DestinyCharacterComponent struct {
	MembershipID             int64              `json:"membershipId,string"`
	MembershipType           int                `json:"membershipType"`
	CharacterID              int64              `json:"characterId,string"`
	DateLastPlayed           time.Time          `json:"dateLastPlayed"`
	MinutesPlayedThisSession int64              `json:"minutesPlayedThisSession,string"`
	MinutesPlayedTotal       int64              `json:"minutesPlayedTotal,string"`
	Light                    int                `json:"light"`
	Stats                    map[uint]int       `json:"stats"`
	RaceHash                 uint               `json:"raceHash"`
	GenderHash               uint               `json:"genderHash"`
	ClassHash                uint               `json:"classHash"`
	EmblemPath               string             `json:"emblemPath"`
	EmblemBackgroundPath     string             `json:"emblemBackgroundPath"`
	EmblemHash               uint               `json:"string"`
	EmblemColor              DestinyColor       `json:"emblemColor"`
	LevelProgression         DestinyProgression `json:"levelProgression"`
	BaseCharacterLevel       int                `json:"baseCharacterLevel"`
	PercentToNextLevel       float32            `json:"percentToNextLevel"`
	TitleRecordHash          *uint              `json:"titleRecordHash"`
}

// DestinyColor ...
// https://bungie-net.github.io/multi/schema_Destiny-Misc-DestinyColor.html#schema_Destiny-Misc-DestinyColor
type DestinyColor struct {
	Red   byte `json:"red"`
	Green byte `json:"green"`
	Blue  byte `json:"blue"`
	Alpha byte `json:"alpha"`
}

// DestinyProgression ...
// https://bungie-net.github.io/multi/schema_Destiny-DestinyProgression.html#schema_Destiny-DestinyProgression
type DestinyProgression struct {
	ProgressionHash     uint                           `json:"progressionHash"`
	DailyProgress       int                            `json:"dailyProgress"`
	DailyLimit          int                            `json:"dailyLimit"`
	WeeklyProgress      int                            `json:"weeklyProgress"`
	WeeklyLimit         int                            `json:"weeklyLimit"`
	CurrentProgress     int                            `json:"currentProgress"`
	Level               int                            `json:"level"`
	LevelCap            int                            `json:"levelCap"`
	StepIndex           int                            `json:"stepIndex"`
	ProgressToNextLevel int                            `json:"progressToNextLevel"`
	NextLevelAt         int                            `json:"nextLevelAt"`
	CurrentResetCount   *int                           `json:"currentResetCount"`
	SeasonResets        []DestinyProgressionResetEntry `json:"seasonResets"`
	RewardItemStates    []int                          `json:"rewardItemStates"`
}

// DestinyProgressionResetEntry ...
// https://bungie-net.github.io/multi/schema_Destiny-DestinyProgressionResetEntry.html#schema_Destiny-DestinyProgressionResetEntry
type DestinyProgressionResetEntry struct {
	Season int `json:"season"`
	Resets int `json:"resets"`
}

// func (ds *Destiny2Service) GetProfile(membershipType int, membershipID int64, opts ...RequestOption) error {

// }

func (gs *Destiny2Service) do(method, endpoint string, dst interface{}, opts ...RequestOption) error {
	endpoint = path.Join("/Destiny2", endpoint)
	return gs.c.do(method, endpoint, dst, opts...)
}
