package ebs

import "time"

type User struct {
	ID            int       `json:"user_id" sql:"USER_ID"`
	UserName      string    `json:"user_name" sql:"USER_NAME"`
	EmailAddress  string    `json:"email_address" sql:"EMAIL_ADDRESS"`
	Description   string    `json:"description" sql:"DESCRIPTION"`
	EmployeeID    int       `json:"employee_id" sql:"EMPLOYEE_ID"`
	LastLogonDate time.Time `json:"last_logon_date" sql:"LAST_LOGON_DATE"`
	CreatedAt     time.Time `json:"creation_date" sql:"CREATION_DATE"`
	StartDate     time.Time `json:"start_date" sql:"START_DATE"`
	EndDate       time.Time `json:"end_date" sql:"END_DATE"`
}

type Role struct {
	ID              int       `json:"role_id" sql:"ROLE_ID"`
	Name            string    `json:"role_name" sql:"ROLE_NAME"`
	Type            string    `json:"role_type_cd" sql:"ROLE_TYPE_CD"`
	BusinessGroupID int       `json:"business_group_id" sql:"BUSINESS_GROUP_ID"`
	CreatedAt       time.Time `json:"creation_date" sql:"CREATION_DATE"`
}
