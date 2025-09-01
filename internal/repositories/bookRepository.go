package repositories

import (
	"fmt"
	"gin_main/internal/repositories/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookRepositoryInterface interface {
	Create(book entities.Book) (entities.Book, error)                                                      // создаёт книгу и возвращает ID или созданный объект книги
	Update(book entities.Book) error                                                                       // изменяет конкретную книгу
	FindById(id uuid.UUID) (entities.Book, error)                                                          // найдёт книгу по конкретному id
	FindByParameters(title, author string, yearOfWriting, yearOfBirth *time.Time) ([]entities.Book, error) // найдёт по параметрам (автор, название, год) | мне могут передать ФИО полностью, ФИО с инициалами, только фамилию или год рождения или год написания
	GetAll() ([]entities.Book, error)                                                                      // возвращает все книги, должен возвращать потоком
	ChangeQuantity(id uuid.UUID, quantity int) (int, error)                                                // изменяет количество остатка для книги по id
}

type bookRepository struct {
	database *gorm.DB
}

func NewBookRepository(database *gorm.DB) BookRepositoryInterface {
	return &bookRepository{database: database}
}

func (r *bookRepository) Create(book entities.Book) (entities.Book, error) {
	if result := r.database.Create(&book); result.Error != nil {
		return entities.Book{}, result.Error
	}
	return book, nil
}

func (r *bookRepository) Update(book entities.Book) error {
	if result := r.database.Updates(&book); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *bookRepository) FindById(id uuid.UUID) (entities.Book, error) {
	var book entities.Book
	if result := r.database.Find(&book, id); result != nil {
		return entities.Book{}, result.Error
	}
	return book, nil
}

func (r *bookRepository) FindByParameters(title, author string, yearOfWriting, yearOfBirth *time.Time) ([]entities.Book, error) {
	var books []entities.Book
	if title != "" {
		r.database.Where("title = ?", title)
	}
	if author != "" {
		r.database.Where("author_id = (select id from author where full_name like ?)", author) //TODO: вернуться и решить как правильно искать автора разными способами
	}
	if yearOfWriting != nil {
		r.database.Where("yearOfWriting = ?", *yearOfWriting)
	}
	if yearOfBirth != nil {
		r.database.Where("yearOfBirth = ?", *yearOfBirth)
	}
	if results := r.database.Find(&books); results != nil {
		return nil, results.Error
	}
	return books, nil
}

func (r *bookRepository) GetAll() ([]entities.Book, error) {
	var books []entities.Book
	if results := r.database.Find(&books); results.Error != nil {
		return nil, results.Error
	}
	return books, nil
}

func (r *bookRepository) ChangeQuantity(id uuid.UUID, quantity int) (int, error) {
	var newQuantity int
	err := r.database.Transaction(func(tx *gorm.DB) error {
		var current int
		if err := tx.Table("book").Select("quantity").Where("id = ?", id).Scan(&current).Error; err != nil {
			return err
		}
		if current+quantity < 0 {
			return fmt.Errorf("quantity cannot be negative") // TODO: создать отдельный файл с ошибками
		}
		if err := tx.Exec("update book set quantity = quantity + ? where id = ?", quantity, id).Error; err != nil {
			return err
		}
		newQuantity = current + quantity
		return nil
	})
	if err != nil {
		return 0, err
	}
	return newQuantity, nil
}
