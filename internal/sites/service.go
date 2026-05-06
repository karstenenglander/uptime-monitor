package sites

import (
	"context"
)

type Service interface {
	ListSites(ctx context.Context) (*struct {Sites []string `json:"sites"`}, error)
	AddSite(ctx context.Context) (*struct {Sites []string `json:"sites"`}, error)
	RemoveSite(ctx context.Context) (*struct {Sites []string `json:"sites"`}, error)
}

type svc struct {
}

func NewService() Service {
	return &svc{}
}

func (s *svc) ListSites(ctx context.Context) (*struct {Sites []string `json:"sites"`}, error) {
	sites := struct {Sites []string `json:"sites"`}{}
	return &sites, nil
}

func (s *svc) AddSite(ctx context.Context) (*struct {Sites []string `json:"sites"`}, error) {
	sites := struct {Sites []string `json:"sites"`}{}
	return &sites, nil
}

func (s *svc) RemoveSite(ctx context.Context) (*struct {Sites []string `json:"sites"`}, error) {
	sites := struct {Sites []string `json:"sites"`}{}
	return &sites, nil
}
