package setting

var (
	instance *globalSetting
)

type globalSetting struct {
	Values *settingValues
}

func New() *globalSetting {
	if instance != nil {
		return instance
	}

	instance = &globalSetting{
		Values: &settingValues{
			SearchBehavior: FUZZY_MATCH,
		},
	}
	return instance
}

func (g *globalSetting) SetValues(value *settingValues) {
	if instance == nil {
		return
	}

	g.Values = value
}
