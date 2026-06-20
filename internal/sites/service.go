package sites

import (
	"context"
	repo "uptime-monitor/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListSites(ctx context.Context) ([]repo.Site, error)
	AddSite(ctx context.Context, params createAddParams) (int64, error)
	RemoveSite(ctx context.Context, params createIdParams) (string, error)
	FindSiteByID(ctx context.Context, params createIdParams) (repo.Site, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListSites(ctx context.Context) ([]repo.Site, error) {
	return s.repo.ListSites(ctx)
}

func (s *svc) FindSiteByID(ctx context.Context, params createIdParams) (repo.Site, error) {
	id := params.Id
	return s.repo.FindSiteByID(ctx, id)
}

func (s *svc) AddSite(ctx context.Context, params createAddParams) (int64, error) {
	name := params.Name
	url := params.Url
	return s.repo.AddSite(ctx, repo.AddSiteParams{Name: name, Url: url})
}

func (s *svc) RemoveSite(ctx context.Context, params createIdParams) (string, error) {
	id := params.Id
	return s.repo.RemoveSiteByID(ctx, id)
}
