package auth

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ChewX3D/wbcli/internal/app/ports"
)

// ListedProfile is safe auth profile row for output.
type ListedProfile struct {
	Name       string
	Backend    string
	APIKeyHint string
	UpdatedAt  time.Time
	Active     bool
}

// ListResult contains metadata-only profile listing.
type ListResult struct {
	Profiles []ListedProfile
}

// ListService returns non-secret profile metadata.
type ListService struct {
	profileStore ports.ProfileStore
}

// NewListService constructs ListService.
func NewListService(profileStore ports.ProfileStore) *ListService {
	return &ListService{profileStore: profileStore}
}

// Execute returns all profile metadata.
func (service *ListService) Execute(ctx context.Context) (ListResult, error) {
	profiles, err := service.profileStore.ListProfiles(ctx)
	if err != nil {
		return ListResult{}, fmt.Errorf("list profile metadata: %w", err)
	}
	activeProfile, err := service.profileStore.GetActiveProfile(ctx)
	if err != nil {
		return ListResult{}, fmt.Errorf("get active profile: %w", err)
	}

	rows := make([]ListedProfile, 0, len(profiles))
	for _, profile := range profiles {
		rows = append(rows, ListedProfile{
			Name:       profile.Name,
			Backend:    profile.Backend,
			APIKeyHint: profile.APIKeyHint,
			UpdatedAt:  profile.UpdatedAt,
			Active:     profile.Name == activeProfile,
		})
	}

	sort.Slice(rows, func(indexA, indexB int) bool {
		return rows[indexA].Name < rows[indexB].Name
	})

	return ListResult{Profiles: rows}, nil
}
