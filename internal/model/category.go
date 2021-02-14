package model

// TabCategoryRaceMap ...
type TabCategoryRaceMap struct {
	MapID      int    `gorm:"column:map_id" json:"map_id"`
	CategoryID string `gorm:"column:category_id" json:"category_id"`
	RaceID     string `gorm:"column:race_id" json:"race_id"`
}

// TableName ...
func (TabCategoryRaceMap) TableName() string {
	return "tab_category_race_map"
}

// GetCatInfo ...
func GetCatInfo(raceID string) (data TabCategoryRaceMap, err error) {
	err = DBCat.Model(&TabCategoryRaceMap{}).
		Where("race_id=?", raceID).
		First(&data).Error
	return data, err
}

// GetCatRaces ...
func GetCatRaces(cID string) (data []TabCategoryRaceMap, err error) {
	err = DBCat.Model(&TabCategoryRaceMap{}).
		Where("category_id=?", cID).
		Find(&data).Error
	return data, err
}
