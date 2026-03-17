package repository

import (
	"context"
	"nethttp/entity"

	"github.com/jmoiron/sqlx"
)

type PersonRepository interface {
	CreatePerson(ctx context.Context, person *entity.Person) error
	GetPerson(ctx context.Context, id int) (*entity.Person, error)
	GetAllPersons(ctx context.Context) ([]entity.Person, error)
}

type personRepositoryImpl struct {
	db *sqlx.DB
}

func NewPersonRepository(db *sqlx.DB) PersonRepository {
	return &personRepositoryImpl{db: db}
}

func (r *personRepositoryImpl) CreatePerson(ctx context.Context, person *entity.Person) error {
	query := `INSERT INTO persons (name, email) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, person.Name, person.Email).Scan(&person.ID)
}

func (r *personRepositoryImpl) GetPerson(ctx context.Context, id int) (*entity.Person, error) {
	var person entity.Person
	query := `SELECT id, name, email FROM persons WHERE id = $1`
	err := r.db.GetContext(ctx, &person, query, id)
	return &person, err
}

func (r *personRepositoryImpl) GetAllPersons(ctx context.Context) ([]entity.Person, error) {
	var persons []entity.Person
	query := `SELECT id, name, email FROM persons`
	err := r.db.SelectContext(ctx, &persons, query)
	return persons, err
}
