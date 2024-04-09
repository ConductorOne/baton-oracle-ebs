package ebs

import "time"

// Table name: FND_USER.
type User struct {
	ID            int        `json:"user_id" sql:"USER_ID"`
	UserName      string     `json:"user_name" sql:"USER_NAME"`
	EmailAddress  string     `json:"email_address" sql:"EMAIL_ADDRESS"`
	Description   string     `json:"description" sql:"DESCRIPTION"`
	EmployeeID    int        `json:"employee_id" sql:"EMPLOYEE_ID"`
	LastLogonDate *time.Time `json:"last_logon_date" sql:"LAST_LOGON_DATE"`
	CreatedAt     *time.Time `json:"creation_date" sql:"CREATION_DATE"`
	StartDate     *time.Time `json:"start_date" sql:"START_DATE"`
	EndDate       *time.Time `json:"end_date" sql:"END_DATE"`
	Group         string     `json:"security_grou_id" sql:"SECURITY_GROUP_ID"`
}

// Table name: PQH_ROLES.
type Role struct {
	ID              int        `json:"role_id" sql:"ROLE_ID"`
	Name            string     `json:"role_name" sql:"ROLE_NAME"`
	Type            string     `json:"role_type_cd" sql:"ROLE_TYPE_CD"`
	BusinessGroupID int        `json:"business_group_id" sql:"BUSINESS_GROUP_ID"`
	CreatedAt       *time.Time `json:"creation_date" sql:"CREATION_DATE"`
}

// Table name: FND_RESPONSIBILITY.
type Responsibility struct {
	ID           int    `json:"responsibility_id" sql:"RESPONSIBILITY_ID"`
	AppID        int    `json:"application_id" sql:"APPLICATION_ID"`
	MenuID       int    `json:"menu_id" sql:"MENU_ID"`
	DataGroupApp int    `json:"data_group_application_id" sql:"DATA_GROUP_APPLICATION_ID"`
	DataGroup    int    `json:"data_group_id" sql:"DATA_GROUP_ID"`
	GroupApp     int    `json:"group_application_id" sql:"GROUP_APPLICATION_ID"`
	RequestGroup int    `json:"request_group_id" sql:"REQUEST_GROUP_ID"`
	Key          string `json:"responsibility_key" sql:"RESPONSIBILITY_KEY"`
	Group        string `json:"security_group_id" sql:"SECURITY_GROUP_ID"`

	CreatedAt *time.Time `json:"creation_date" sql:"CREATION_DATE"`
	StartDate *time.Time `json:"start_date" sql:"START_DATE"`
	EndDate   *time.Time `json:"end_date" sql:"END_DATE"`
}

type Config struct {
	Port    int    `mapstructure:"port"`
	Server  string `mapstructure:"server"`
	Service string `mapstructure:"service"`

	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
