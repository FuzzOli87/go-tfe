package tfe

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func convertToAdminWorkspace(ws *Workspace) *AdminWorkspace {
	return &AdminWorkspace{
		ID:           ws.ID,
		Locked:       ws.Locked,
		Name:         ws.Name,
		VCSRepo:      ws.VCSRepo,
		CurrentRun:   ws.CurrentRun,
		Organization: ws.Organization,
	}
}

func TestAdminWorkspaceList(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createOrganization(t, client)
	defer orgTestCleanup()

	wTest1, wTest1Cleanup := createWorkspace(t, client, orgTest)
	defer wTest1Cleanup()

	wTest2, wTest2Cleanup := createWorkspace(t, client, orgTest)
	defer wTest2Cleanup()

	t.Run("without list options", func(t *testing.T) {
		wl, err := client.Admin.Workspaces.List(ctx, AdminWorkspaceListOptions{})

		require.NoError(t, err)
		assert.Contains(t, wl.Items, convertToAdminWorkspace(wTest1))
		assert.Contains(t, wl.Items, convertToAdminWorkspace(wTest2))
		assert.Equal(t, 1, wl.CurrentPage)
		assert.Equal(t, 2, wl.TotalCount)
	})

	t.Run("with list options", func(t *testing.T) {
		// Request a page number which is out of range. The result should
		// be successful, but return no results if the paging options are
		// properly passed along.
		wl, err := client.Admin.Workspaces.List(ctx, AdminWorkspaceListOptions{
			ListOptions: ListOptions{
				PageNumber: 999,
				PageSize:   100,
			},
		})
		require.NoError(t, err)
		assert.Empty(t, wl.Items)
		assert.Equal(t, 999, wl.CurrentPage)
		assert.Equal(t, 2, wl.TotalCount)
	})
}
