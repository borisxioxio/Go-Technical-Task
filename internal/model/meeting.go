package model

// TabMeeting ...
type TabMeeting struct {
	MeetingID   string `gorm:"column:meeting_id" json:"meeting_id"`
	MeetingName string `gorm:"column:meeting_name" json:"meeting_name"`
}

// TableName ...
func (TabMeeting) TableName() string {
	return "tab_meeting"
}

// GetMetInfo ...
func GetMetInfo(mID string) (data TabMeeting, err error) {
	err = DBMeeting.Model(&TabMeeting{}).
		Where("meeting_id=?", mID).
		First(&data).Error
	return data, err
}
