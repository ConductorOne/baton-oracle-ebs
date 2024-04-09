package ebs

import (
	"context"

	ora "github.com/sijms/go-ora/v2"
)

const (
	UsersTable = "APPS.FND_USER"

	LimitToOneRecord = "ROWNUM = 1"
)

var (
	UsersAttributes = []string{
		"USER_ID", "USER_NAME", "EMAIL_ADDRESS", "DESCRIPTION", "EMPLOYEE_ID", "LAST_LOGON_DATE", "CREATION_DATE", "START_DATE", "END_DATE",
	}
)

type Client struct {
	Conn *ora.Connection
}

func NewEBSClient(c *ora.Connection) *Client {
	return &Client{Conn: c}
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

func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	// prepare the SQL statement
	query := ComposeSQLQuery(UsersAttributes, UsersTable, "")
	stmt := ora.NewStmt(query, c.Conn)
	defer stmt.Close()

	// execute the SQL statement
	rows, err := stmt.Query_(nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows and parse the data
	var users []User
	for rows.Next_() {
		var user User

		err := rows.Scan(&user.ID, &user.UserName, &user.EmailAddress, &user.Description, &user.EmployeeID, &user.LastLogonDate, &user.CreatedAt, &user.StartDate, &user.EndDate)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
