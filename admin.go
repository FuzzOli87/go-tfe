package tfe

import (
	"context"
	"fmt"
)

var _ AdminWorkspaces = (*adminWorkspaces)(nil)

type AdminWorkspaces interface {
	List(ctx context.Context, options AdminWorkspaceListOptions) (*AdminWorkspaceList, error)
}

type adminWorkspaces struct {
	client *Client
}

// AdminWorkspaceListOptions represents the options for listing workspaces.
type AdminWorkspaceListOptions struct {
	ListOptions
}

type AdminWorkspaceList struct {
	*Pagination
	Items []*AdminWorkspace
}

type AdminWorkspace struct {
	ID      string   `jsonapi:"primary,workspaces"`
	Locked  bool     `jsonapi:"attr,locked"`
	Name    string   `jsonapi:"attr,name"`
	VCSRepo *VCSRepo `jsonapi:"attr,vcs-repo"`

	// Relations
	CurrentRun   *Run          `jsonapi:"relation,current-run"`
	Organization *Organization `jsonapi:"relation,organization"`
}

type Admin struct {
	Workspaces AdminWorkspaces
}

func (s *adminWorkspaces) List(ctx context.Context, options AdminWorkspaceListOptions) (*AdminWorkspaceList, error) {
	u := fmt.Sprintf("admin/workspaces")

	req, err := s.client.newRequest("GET", u, &options)

	if err != nil {
		return nil, err
	}

	wl := &AdminWorkspaceList{}

	err = s.client.do(ctx, req, wl)

	if err != nil {
		return nil, err
	}

	return wl, nil
}
