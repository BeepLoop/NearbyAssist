package setting

type SearchBehavior string

const (
	EXACT_MATCH SearchBehavior = "exact_match"
	FUZZY_MATCH SearchBehavior = "fuzzy_match"
)

var validSearchBehaviors = map[SearchBehavior]struct{}{
	EXACT_MATCH: {},
	FUZZY_MATCH: {},
}

func IsValidSearchBehavior(s string) bool {
	_, ok := validSearchBehaviors[SearchBehavior(s)]
	return ok
}

type settingValues struct {
	SearchBehavior SearchBehavior
}

func NewSettingValues(searchBehavior SearchBehavior) settingValues {
	return settingValues{
		SearchBehavior: searchBehavior,
	}
}
