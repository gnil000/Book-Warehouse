package repositories

import (
	"gin_main/src/database/entities"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthorRepositoryInterface interface {
	Create(author *entities.Author) error
	FindById(id uuid.UUID) (entities.Author, error)
	FindBySurname(surname string) ([]entities.Author, error)
	Update(author entities.Author) error
}

type authorRepository struct {
	database *gorm.DB
}

func NewAuthorRepository(database *gorm.DB) AuthorRepositoryInterface {
	return &authorRepository{database: database}
}

func (r *authorRepository) Create(author *entities.Author) error {
	author.FirstName = strings.ToTitle(author.FirstName)
	author.SecondName = strings.ToTitle(author.SecondName)
	author.Surname = strings.ToTitle(author.Surname)
	if result := r.database.Create(&author); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *authorRepository) FindById(id uuid.UUID) (entities.Author, error) {
	var author entities.Author
	if resutl := r.database.First(&author, id); resutl.Error != nil {
		return entities.Author{}, resutl.Error
	}
	return author, nil
}

func (r *authorRepository) FindBySurname(surname string) ([]entities.Author, error) {
	var authors []entities.Author
	if resutl := r.database.Where("surname ilike ?", surname).Find(&authors); resutl.Error != nil {
		return nil, resutl.Error
	}
	return authors, nil
}

func (r *authorRepository) Update(author entities.Author) error {
	author.FirstName = strings.ToTitle(author.FirstName)
	author.SecondName = strings.ToTitle(author.SecondName)
	author.Surname = strings.ToTitle(author.Surname)
	var result *gorm.DB
	if result = r.database.Updates(&author); result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
