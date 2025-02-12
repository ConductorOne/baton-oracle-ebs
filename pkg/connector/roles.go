package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/conductorone/baton-oracle-ebs/pkg/ebs"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleBuilder struct {
	client       *ebs.Client
	resourceType *v2.ResourceType
}

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

func roleResource(role *ebs.Role) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_id":           role.ID,
		"role_type":         role.Type,
		"business_group_id": role.BusinessGroupID,
	}

	if role.CreatedAt != nil {
		profile["created_at"] = role.CreatedAt.Format(time.RFC3339)
	}

	options := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	res, err := rs.NewRoleResource(
		role.Name,
		roleResourceType,
		role.ID,
		options,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	err := r.client.Conn.Open()
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to open connection: %w", err)
	}

	defer r.client.Conn.Close()

	bag, offset, err := parsePageToken(pToken.Token, &v2.ResourceId{ResourceType: roleResourceType.Id})
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to parse page token: %w", err)
	}

	pgVars := ebs.NewPaginationVars(offset, ResourcesPageSize)
	roles, pageTotal, err := r.client.ListRoles(ctx, pgVars)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list roles: %w", err)
	}

	var rv []*v2.Resource
	for _, role := range roles {
		ur, err := roleResource(&role) // #nosec G601
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to create role resource: %w", err)
		}

		rv = append(rv, ur)
	}

	next := prepareNextToken(offset, pageTotal)
	nextToken, err := bag.NextToken(next)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to prepare next token: %w", err)
	}

	return rv, nextToken, nil, nil
}

func (r *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newRoleBuilder(client *ebs.Client) *roleBuilder {
	return &roleBuilder{
		client:       client,
		resourceType: roleResourceType,
	}
}
