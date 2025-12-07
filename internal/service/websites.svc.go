package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/users/internal/dto"
	"github.com/ICan-TC/users/internal/models"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type WebsiteService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewWebsiteService(db *bun.DB) *WebsiteService {
	logger := logging.L().With().Str("service", "websites.svc").Logger()
	return &WebsiteService{db: db, log: logger}
}

func (s *WebsiteService) Create(ctx context.Context, org *models.Organization, website *models.Website) (*models.Website, error) {
	org.ID = ulid.Make().String()
	s.log.Debug().Msg("inserting website " + "ID=" + website.ID)
	r, err := s.db.NewInsert().Model(website).Exec(ctx, website)
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("getting rows affected")
	rows, err := r.RowsAffected()
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("checking rows affected")
	if rows != 1 {
		return nil, errors.New("failed to insert website, rows affected: " + strconv.Itoa(int(rows)))
	}
	s.log.Debug().Msg("website inserted, rows affected: " + strconv.Itoa(int(rows)))
	return website, nil
}

func (s *WebsiteService) Update(ctx context.Context, org *models.Organization, website *models.Website) (*models.Website, error) {
	s.log.Debug().Msg("updating website " + "ID=" + website.ID)
	r, err := s.db.NewUpdate().Model(website).WherePK().Exec(ctx, website)
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("getting rows affected")
	rows, err := r.RowsAffected()
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("checking rows affected")
	if rows != 1 {
		return nil, errors.New("failed to update website, rows affected: " + strconv.Itoa(int(rows)))
	}
	s.log.Debug().Msg("website updated, rows affected: " + strconv.Itoa(int(rows)))
	return website, nil
}

func (s *WebsiteService) Delete(ctx context.Context, org *models.Organization, website *models.Website) error {
	s.log.Debug().Msg("deleting website " + "ID=" + website.ID)
	r, err := s.db.NewDelete().Model(website).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	s.log.Debug().Msg("getting rows affected")
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	s.log.Debug().Msg("checking rows affected")
	if rows != 1 {
		return errors.New("failed to delete website, rows affected: " + strconv.Itoa(int(rows)))
	}
	s.log.Debug().Msg("website deleted, rows affected: " + strconv.Itoa(int(rows)))
	return nil
}

func (s *WebsiteService) GetByID(ctx context.Context, id string) (*models.Website, error) {
	website := &models.Website{}
	s.log.Debug().Msg("getting website " + "ID=" + id)
	err := s.db.NewSelect().Model(website).Where("id = ?", id).Scan(ctx, website)
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("website retrieved")
	return website, nil
}

func (s *WebsiteService) GetAll(ctx context.Context, query *dto.ListQuery) ([]*models.Website, error) {
	websites := []*models.Website{}
	s.log.Debug().Msg("getting all websites")
	q := s.db.NewSelect().Model(&websites)
	q.Limit(query.PerPage)
	q.Offset((query.Page - 1) * query.PerPage)
	q.Order(query.SortBy + " " + query.SortDir)
	err := q.Scan(ctx, websites)
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("all websites retrieved")
	return websites, nil
}

func (s *WebsiteService) Search(ctx context.Context, query string) ([]*models.Website, error) {
	websites := []*models.Website{}

	s.log.Debug().Msg("searching websites")
	err := s.db.NewSelect().Model(&websites).
		Where("name LIKE ? OR array_to_string(tags, ',') LIKE ? OR labels::text LIKE ? OR metadata::text LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Limit(100).
		Scan(ctx, websites)

	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("all websites retrieved")

	return websites, nil
}
