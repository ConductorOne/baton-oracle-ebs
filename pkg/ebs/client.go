package ebs

import (
	"context"
	"fmt"

	ora "github.com/sijms/go-ora/v2"
)

const (
	UsersTable = "APPS.FND_USER"
	RolesTable = "APPS.PQH_ROLES"

	LimitToOneRecord = "ROWNUM = 1"
)

var (
	UsersAttributes = []string{
		"USER_ID", "USER_NAME", "EMAIL_ADDRESS", "DESCRIPTION", "EMPLOYEE_ID", "LAST_LOGON_DATE", "CREATION_DATE", "START_DATE", "END_DATE",
	}
	RolesAttributes = []string{
		"ROLE_ID", "ROLE_NAME", "ROLE_TYPE_CD", "BUSINESS_GROUP_ID", "CREATION_DATE",
	}
)

type Client struct {
	Conn *ora.Connection
}

func NewEBSClient(c *ora.Connection) *Client {
	return &Client{Conn: c}
}

func ComposeSQLQuery(attributes []string, table string, pgVars *PaginationVars) string {
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

	// add the pagination to the query
	wrapper := fmt.Sprintf(
		"SELECT * FROM (SELECT a.*, ROWNUM rnum FROM (%s) a WHERE ROWNUM <= %d) WHERE rnum > %d",
		query,
		pgVars.Offset+pgVars.Limit,
		pgVars.Offset,
	)

	return wrapper
}

type PaginationVars struct {
	Offset uint
	Limit  uint
}

func NewPaginationVars(offset, limit uint) *PaginationVars {
	return &PaginationVars{Offset: offset, Limit: limit}
}

func (c *Client) ListUsers(ctx context.Context, pgVars *PaginationVars) ([]User, uint, error) {
	// prepare the SQL statement
	query := ComposeSQLQuery(UsersAttributes, UsersTable, pgVars)
	stmt := ora.NewStmt(query, c.Conn)
	defer stmt.Close()

	// execute the SQL statement
	rows, err := stmt.Query_(nil)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// iterate over the rows and parse the data
	var users []User
	for rows.Next_() {
		var user User
		var rnum int

		err := rows.Scan(&user.ID, &user.UserName, &user.EmailAddress, &user.Description, &user.EmployeeID, &user.LastLogonDate, &user.CreatedAt, &user.StartDate, &user.EndDate, &rnum)
		if err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}

	// stop paginating if the number of records is less than the page size
	if len(users) < int(pgVars.Limit) {
		return users, 0, nil
	}

	return users, uint(len(users)), nil
}

func (c *Client) ListRoles(ctx context.Context, pgVars *PaginationVars) ([]Role, uint, error) {
	// prepare the SQL statement
	query := ComposeSQLQuery(RolesAttributes, RolesTable, pgVars)
	stmt := ora.NewStmt(query, c.Conn)
	defer stmt.Close()

	// execute the SQL statement
	rows, err := stmt.Query_(nil)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// iterate over the rows and parse the data
	var roles []Role
	for rows.Next_() {
		var role Role
		var rnum int

		err := rows.Scan(&role.ID, &role.Name, &role.Type, &role.BusinessGroupID, &role.CreatedAt, &rnum)
		if err != nil {
			return nil, 0, err
		}

		roles = append(roles, role)
	}

	// stop paginating if the number of records is less than the page size
	if len(roles) < int(pgVars.Limit) {
		return roles, 0, nil
	}

	return roles, uint(len(roles)), nil
}
