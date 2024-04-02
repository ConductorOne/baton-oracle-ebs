package ebs

import (
	"context"
	"database/sql"
)

const (
	UsersTable = "FND_USER"

	LimitToOneRecord = "ROWNUM = 1"
)

var (
	UsersAttributes = []string{
		"USER_ID", "USER_NAME", "EMAIL_ADDRESS", "DESCRIPTION", "EMPLOYEE_ID", "LAST_LOGON_DATE", "CREATION_DATE", "START_DATE", "END_DATE",
	}
)

type Client struct {
	db *sql.DB
}

func ComposeSQLQuery(attributes []string, table string, filter string) string {
	if len(attributes) == 0 {
		return ""
	}

	// add the attributes to the query
	query := "SELECT " + attributes[0]

	if len(attributes) > 1 {
		for i := 1; i < len(attributes); i++ {
			query += ", " + attributes[i]
		}
	}

	// add the table to the query
	query += " FROM " + table

	// add the filter to the query
	if filter != "" {
		query += " WHERE " + filter
	}

	return query
}

func NewEBSClient(db *sql.DB) *Client {
	return &Client{db: db}
}

func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := c.db.QueryContext(ctx, ComposeSQLQuery(UsersAttributes, UsersTable, ""))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.UserName, &user.EmailAddress, &user.Description, &user.EmployeeID, &user.LastLogonDate, &user.CreatedAt, &user.StartDate, &user.EndDate)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
