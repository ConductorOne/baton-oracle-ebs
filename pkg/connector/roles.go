package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-oracle-scm/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const (
	roleEntitlementMember = "member"
)

type roleBuilder struct {
	client *client.FusionClient
}

func (r *roleBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return roleResourceType
}

// roleToResource converts an Oracle Fusion Cloud role to a Baton resource.
func (r *roleBuilder) roleToResource(role *client.Role) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_id":   role.RoleID,
		"role_code": role.RoleCode,
	}

	if role.RoleDescription != nil {
		profile["description"] = *role.RoleDescription
	}
	if role.RoleCategory != nil {
		profile["category"] = *role.RoleCategory
	}
	if role.RoleCategoryCode != nil {
		profile["category_code"] = *role.RoleCategoryCode
	}
	if role.AbstractRole != nil {
		profile["abstract_role"] = *role.AbstractRole
	}

	roleTraitOptions := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	// Use RoleID as the stable resource identifier.
	resourceID := strconv.FormatInt(role.RoleID, 10)
	ret, err := rs.NewRoleResource(role.RoleName, roleResourceType, resourceID, roleTraitOptions)
	if err != nil {
		return nil, fmt.Errorf("baton-oracle-scm: failed to create role resource: %w", err)
	}

	return ret, nil
}

// List returns all roles from Oracle Fusion Cloud as resource objects.
func (r *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	offset, err := parsePageToken(pToken.Token)
	if err != nil {
		return nil, "", nil, err
	}

	outputAnnotations := annotations.New()
	roles, hasMore, rateLimit, err := r.client.ListRoles(ctx, offset)
	outputAnnotations.WithRateLimiting(rateLimit)
	if err != nil {
		return nil, "", outputAnnotations, fmt.Errorf("baton-oracle-scm: failed to list roles: %w", err)
	}

	resources := make([]*v2.Resource, 0, len(roles))
	for _, role := range roles {
		roleResource, err := r.roleToResource(role)
		if err != nil {
			return nil, "", outputAnnotations, err
		}
		resources = append(resources, roleResource)
	}

	nextOffset := client.GetNextOffset(hasMore, offset, client.DefaultPageSize)
	nextToken := formatNextPageToken(nextOffset)

	return resources, nextToken, outputAnnotations, nil
}

// Entitlements returns a "member" entitlement for each role resource.
func (r *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	entitlementResource := ent.NewAssignmentEntitlement(
		resource,
		roleEntitlementMember,
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Role Member", resource.DisplayName)),
		ent.WithDescription(fmt.Sprintf("Member of the %s role in Oracle Fusion Cloud", resource.DisplayName)),
	)

	return []*v2.Entitlement{entitlementResource}, "", nil, nil
}

// Grants returns the user-role assignments for a given role.
// This queries all user accounts and checks their role assignments.
func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	// Parse the role's resource ID to get the RoleID we're looking for.
	roleIDStr := resource.Id.Resource
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-oracle-scm: failed to parse role ID %q: %w", roleIDStr, err)
	}

	// We use page token to track pagination of user accounts.
	offset, err := parsePageToken(pToken.Token)
	if err != nil {
		return nil, "", nil, err
	}

	outputAnnotations := annotations.New()

	// List user accounts for this page.
	users, hasMore, rateLimit, err := r.client.ListUserAccounts(ctx, offset)
	outputAnnotations.WithRateLimiting(rateLimit)
	if err != nil {
		return nil, "", outputAnnotations, fmt.Errorf("baton-oracle-scm: failed to list users for grants: %w", err)
	}

	var grants []*v2.Grant

	// For each user, check if they have this role assigned.
	for _, user := range users {
		userRoles, _, userRateLimit, err := r.client.ListUserRoles(ctx, user.UserID, 0)
		outputAnnotations.WithRateLimiting(userRateLimit)
		if err != nil {
			// Log but continue - some users may not have accessible role data.
			continue
		}

		for _, userRole := range userRoles {
			if userRole.RoleID == roleID {
				userResourceID := strconv.FormatInt(user.UserID, 10)
				principalID := &v2.ResourceId{
					ResourceType: userResourceType.Id,
					Resource:     userResourceID,
				}
				g := grant.NewGrant(resource, roleEntitlementMember, principalID)
				grants = append(grants, g)
				break
			}
		}
	}

	nextOffset := client.GetNextOffset(hasMore, offset, client.DefaultPageSize)
	nextToken := formatNextPageToken(nextOffset)

	return grants, nextToken, outputAnnotations, nil
}

func newRoleBuilder(fusionClient *client.FusionClient) *roleBuilder {
	return &roleBuilder{
		client: fusionClient,
	}
}
