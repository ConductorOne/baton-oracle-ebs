package connector

import (
	"context"
	"time"

	"github.com/conductorone/baton-oracle-ebs/pkg/ebs"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client       *ebs.Client
	resourceType *v2.ResourceType
}

func (u *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func userResource(user *ebs.User) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"user_id":        user.ID,
		"start_date":     user.StartDate.Format(time.RFC3339),
		"description":    user.Description,
		"employee_id":    user.EmployeeID,
		"security_group": user.Group,
	}

	if user.EndDate != nil {
		profile["end_date"] = user.EndDate.Format(time.RFC3339)
	}

	status := v2.UserTrait_Status_STATUS_ENABLED
	if user.EndDate != nil {
		status = v2.UserTrait_Status_STATUS_DISABLED
	}

	options := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(status),
		rs.WithEmail(user.EmailAddress, true),
		rs.WithUserLogin(user.UserName),
	}

	if user.LastLogonDate != nil {
		options = append(options, rs.WithLastLogin(*user.LastLogonDate))
	}

	if user.CreatedAt != nil {
		options = append(options, rs.WithCreatedAt(*user.CreatedAt))
	}

	res, err := rs.NewUserResource(
		user.UserName,
		userResourceType,
		user.ID,
		options,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (u *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	err := u.client.Conn.Open()
	if err != nil {
		return nil, "", nil, err
	}

	defer u.client.Conn.Close()

	users, err := u.client.ListUsers(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	var rv []*v2.Resource
	for _, user := range users {
		ur, err := userResource(&user) // #nosec G601
		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, ur)
	}

	return rv, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (u *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (u *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *ebs.Client) *userBuilder {
	return &userBuilder{
		client:       client,
		resourceType: userResourceType,
	}
}
