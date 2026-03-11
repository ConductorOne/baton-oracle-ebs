package client

import "fmt"

// ErrorResponse represents an error response from Oracle Fusion Cloud REST APIs.
type ErrorResponse struct {
	Type    *string `json:"type,omitempty"`
	Title   *string `json:"title,omitempty"`
	Detail  *string `json:"detail,omitempty"`
	Status  *int    `json:"status,omitempty"`
	ErrCode *string `json:"o:errorCode,omitempty"`
	ErrPath *string `json:"o:errorPath,omitempty"`
}

// Message implements the error interface for uhttp.WithErrorResponse.
func (e *ErrorResponse) Message() string {
	title := "unknown"
	detail := "no detail"
	if e.Title != nil {
		title = *e.Title
	}
	if e.Detail != nil {
		detail = *e.Detail
	}
	if e.ErrCode != nil {
		return fmt.Sprintf("code: %s, title: %s, detail: %s", *e.ErrCode, title, detail)
	}
	return fmt.Sprintf("title: %s, detail: %s", title, detail)
}

// FusionListResponse represents a paginated list response from Oracle Fusion Cloud REST APIs.
type FusionListResponse[T any] struct {
	Items   []*T   `json:"items"`
	Count   int    `json:"count"`
	HasMore bool   `json:"hasMore"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Links   []Link `json:"links,omitempty"`
}

// Link represents a HATEOAS link in Oracle Fusion Cloud REST API responses.
type Link struct {
	Rel    string `json:"rel"`
	Href   string `json:"href"`
	Name   string `json:"name,omitempty"`
	Kind   string `json:"kind,omitempty"`
	Method string `json:"method,omitempty"`
}

// UserAccount represents a user account from the HCM User Accounts REST API.
// Endpoint: /hcmRestApi/resources/latest/userAccounts
type UserAccount struct {
	UserID          int64   `json:"UserId"`
	Username        string  `json:"Username"`
	PersonID        *int64  `json:"PersonId,omitempty"`
	PersonNumber    *string `json:"PersonNumber,omitempty"`
	DisplayName     *string `json:"DisplayName,omitempty"`
	FirstName       *string `json:"FirstName,omitempty"`
	LastName        *string `json:"LastName,omitempty"`
	EmailAddress    *string `json:"WorkEmailAddress,omitempty"`
	Suspended       *bool   `json:"Suspended,omitempty"`
	CredentialsType *string `json:"CredentialsEmailSentFlag,omitempty"`
	GUID            *string `json:"GUID,omitempty"`
	Links           []Link  `json:"links,omitempty"`
}

// Role represents a role from the Common Features Roles REST API.
// Endpoint: /hcmRestApi/resources/latest/roles
type Role struct {
	RoleID            int64   `json:"RoleId"`
	RoleCode          string  `json:"RoleCode"`
	RoleName          string  `json:"RoleName"`
	RoleDescription   *string `json:"RoleDescription,omitempty"`
	RoleCategory      *string `json:"RoleCategory,omitempty"`
	AbstractRole      *bool   `json:"AbstractRole,omitempty"`
	RoleCategoryCode  *string `json:"RoleCategoryCode,omitempty"`
	ExternallyManaged *bool   `json:"ExternallyManaged,omitempty"`
	Links             []Link  `json:"links,omitempty"`
}

// UserRole represents a role assignment for a user.
// Endpoint: /hcmRestApi/resources/latest/userAccounts/{UserId}/child/userAccountRoles
type UserRole struct {
	RoleID            int64   `json:"RoleId"`
	RoleMappingID     int64   `json:"RoleMappingId"`
	RoleCode          string  `json:"RoleCode"`
	RoleName          string  `json:"RoleName"`
	RequestID         *int64  `json:"RequestId,omitempty"`
	AssignmentType    *string `json:"AssignmentType,omitempty"`
	ProvisioningState *string `json:"RoleProvisioningState,omitempty"`
	Links             []Link  `json:"links,omitempty"`
}
