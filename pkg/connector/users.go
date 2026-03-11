package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-oracle-scm/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client *client.FusionClient
}

func (u *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

// userToResource converts an Oracle Fusion Cloud user account to a Baton resource.
func (u *userBuilder) userToResource(user *client.UserAccount) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"user_id":  user.UserID,
		"username": user.Username,
	}

	if user.PersonID != nil {
		profile["person_id"] = *user.PersonID
	}
	if user.PersonNumber != nil {
		profile["person_number"] = *user.PersonNumber
	}
	if user.GUID != nil {
		profile["guid"] = *user.GUID
	}

	userTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithUserLogin(user.Username),
	}

	// Set email if available.
	if user.EmailAddress != nil && *user.EmailAddress != "" {
		userTraitOptions = append(userTraitOptions, rs.WithEmail(*user.EmailAddress, true))
	}

	// Determine user status based on Suspended field.
	userStatus := v2.UserTrait_Status_STATUS_ENABLED
	if user.Suspended != nil && *user.Suspended {
		userStatus = v2.UserTrait_Status_STATUS_DISABLED
	}
	userTraitOptions = append(userTraitOptions, rs.WithStatus(userStatus))

	// Use display name if available, fall back to username.
	displayName := user.Username
	if user.DisplayName != nil && *user.DisplayName != "" {
		displayName = *user.DisplayName
	}

	// Use UserID as the stable resource identifier.
	resourceID := strconv.FormatInt(user.UserID, 10)
	ret, err := rs.NewUserResource(displayName, userResourceType, resourceID, userTraitOptions)
	if err != nil {
		return nil, fmt.Errorf("baton-oracle-scm: failed to create user resource: %w", err)
	}

	return ret, nil
}

// List returns all user accounts from Oracle Fusion Cloud as resource objects.
func (u *userBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	offset, err := parsePageToken(pToken.Token)
	if err != nil {
		return nil, "", nil, err
	}

	outputAnnotations := annotations.New()
	users, hasMore, rateLimit, err := u.client.ListUserAccounts(ctx, offset)
	outputAnnotations.WithRateLimiting(rateLimit)
	if err != nil {
		return nil, "", outputAnnotations, fmt.Errorf("baton-oracle-scm: failed to list users: %w", err)
	}

	resources := make([]*v2.Resource, 0, len(users))
	for _, user := range users {
		userResource, err := u.userToResource(user)
		if err != nil {
			return nil, "", outputAnnotations, err
		}
		resources = append(resources, userResource)
	}

	nextOffset := client.GetNextOffset(hasMore, offset, client.DefaultPageSize)
	nextToken := formatNextPageToken(nextOffset)

	return resources, nextToken, outputAnnotations, nil
}

// Entitlements always returns an empty slice for users.
func (u *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have entitlements.
func (u *userBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(fusionClient *client.FusionClient) *userBuilder {
	return &userBuilder{
		client: fusionClient,
	}
}
