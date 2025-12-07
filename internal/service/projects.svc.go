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

type ProjectService struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewProjectService(db *bun.DB) *ProjectService {
	logger := logging.L().With().Str("service", "projects.svc").Logger()
	return &ProjectService{db: db, log: logger}
}

func (s *ProjectService) Create(ctx context.Context, org *models.Organization, project *models.Project) (*models.Project, error) {
	org.ID = ulid.Make().String()
	s.log.Debug().Msg("inserting project " + "ID=" + project.ID)
	r, err := s.db.NewInsert().Model(project).Exec(ctx, project)
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
		return nil, errors.New("failed to insert project, rows affected: " + strconv.Itoa(int(rows)))
	}
	s.log.Debug().Msg("project inserted, rows affected: " + strconv.Itoa(int(rows)))
	return project, nil
}

func (s *ProjectService) Update(ctx context.Context, project *models.Project) (*models.Project, error) {
	s.log.Debug().Msg("updating project " + "ID=" + project.ID)
	r, err := s.db.NewUpdate().Model(project).WherePK().Exec(ctx, project)
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
		return nil, errors.New("failed to update project, rows affected: " + strconv.Itoa(int(rows)))
	}
	s.log.Debug().Msg("project updated, rows affected: " + strconv.Itoa(int(rows)))
	return project, nil
}

func (s *ProjectService) Delete(ctx context.Context, project *models.Project) error {
	s.log.Debug().Msg("deleting project " + "ID=" + project.ID)
	r, err := s.db.NewDelete().Model(project).WherePK().Exec(ctx)
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
		return errors.New("failed to delete project, rows affected: " + strconv.Itoa(int(rows)))
	}
	s.log.Debug().Msg("project deleted, rows affected: " + strconv.Itoa(int(rows)))
	return nil
}

func (s *ProjectService) GetByID(ctx context.Context, id string) (*models.Project, error) {
	project := &models.Project{}
	s.log.Debug().Msg("getting project " + "ID=" + id)
	err := s.db.NewSelect().Model(project).Where("id = ?", id).Scan(ctx, project)
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("project retrieved")
	return project, nil
}

func (s *ProjectService) GetAll(ctx context.Context, query *dto.ListQuery) ([]*models.Project, error) {
	projects := []*models.Project{}
	s.log.Debug().Msg("getting all projects")
	q := s.db.NewSelect().Model(&projects)
	q.Limit(query.PerPage)
	q.Offset((query.Page - 1) * query.PerPage)
	q.Order(query.SortBy + " " + query.SortDir)
	err := q.Scan(ctx, projects)
	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("all projects retrieved")
	return projects, nil
}

func (s *ProjectService) Search(ctx context.Context, query string) ([]*models.Project, error) {
	projects := []*models.Project{}

	s.log.Debug().Msg("searching projects")
	err := s.db.NewSelect().Model(&projects).
		Where("name LIKE ? OR array_to_string(tags, ',') LIKE ? OR labels::text LIKE ? OR metadata::text LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Limit(100).
		Scan(ctx, projects)

	if err != nil {
		return nil, err
	}
	s.log.Debug().Msg("all projects retrieved")

	return projects, nil
}
