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

type OrganizationService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewOrganizationService(db *bun.DB) *OrganizationService {
	logger := logging.L().With().Str("service", "organizations.svc").Logger()
	return &OrganizationService{db: db, log: logger}
}

func (s *OrganizationService) Create(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	org.ID = ulid.Make().String()
	s.log.Debug().Msg("inserting organization " + "ID=" + org.ID)
	r, err := s.db.NewInsert().Model(org).Exec(ctx, org)
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
		return nil, errors.New("failed to insert organization, rows affected: " + strconv.Itoa(int(rows)))
	}
	s.log.Debug().Msg("organization inserted, rows affected: " + strconv.Itoa(int(rows)))
	return org, nil
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	s.log.Debug().Msg("updating organization " + "ID=" + org.ID)
	r, err := s.db.NewUpdate().Model(org).WherePK().Exec(ctx, org)
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
		return nil, errors.New("failed to update organization, rows affected: " + strconv.Itoa(int(rows)))
	}
	s.log.Debug().Msg("organization updated, rows affected: " + strconv.Itoa(int(rows)))
	return org, nil
}

func (s *OrganizationService) Delete(ctx context.Context, org *models.Organization) error {
	s.log.Debug().Msg("deleting organization " + "ID=" + org.ID)
	r, err := s.db.NewDelete().Model(org).WherePK().Exec(ctx)
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
		return errors.New("failed to delete organization, rows affected: " + strconv.Itoa(int(rows)))
	}
	s.log.Debug().Msg("organization deleted, rows affected: " + strconv.Itoa(int(rows)))
	return nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	org := &models.Organization{}
	s.log.Debug().Msg("getting organization " + "ID=" + id)
	err := s.db.NewSelect().Model(org).Where("id = ?", id).Scan(ctx, org)
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("organization retrieved")
	return org, nil
}

func (s *OrganizationService) GetAll(ctx context.Context, query *dto.ListQuery) ([]*models.Organization, error) {
	orgs := []*models.Organization{}
	s.log.Debug().Msg("getting all organizations")
	q := s.db.NewSelect().Model(&orgs)
	q.Limit(query.PerPage)
	q.Offset((query.Page - 1) * query.PerPage)
	q.Order(query.SortBy + " " + query.SortDir)
	err := q.Scan(ctx, orgs)
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("all organizations retrieved")
	return orgs, nil
}

func (s *OrganizationService) Search(ctx context.Context, query string) ([]*models.Organization, error) {
	orgs := []*models.Organization{}

	s.log.Debug().Msg("searching organizations")
	err := s.db.NewSelect().Model(&orgs).
		Where("name LIKE ? OR array_to_string(tags, ',') LIKE ? OR labels::text LIKE ? OR metadata::text LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Limit(100).
		Scan(ctx, orgs)

	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("all organizations retrieved")

	return orgs, nil
}
