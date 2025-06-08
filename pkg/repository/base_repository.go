package repository

import "gorm.io/gorm"

type BaseRepository[T any] interface {
	Create(model *T) error
	GetByID(id string) (*T, error)
	FindBy(condition map[string]interface{}) (*T, error)
	FindAllBy(condition map[string]interface{}) ([]T, error)
	UpdateByID(id string, updates map[string]interface{}) error
	DeleteByID(id string) error
	SoftDeleteByID(id string) error
	HardDeleteByID(id string) error
	Exists(condition map[string]interface{}) (bool, error)
	Count(condition map[string]interface{}) (int64, error)
	RawQuery(query string, args ...interface{}) *gorm.DB
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{db: db}
}

func (r *baseRepository[T]) Create(model *T) error {
	return r.db.Create(model).Error
}

func (r *baseRepository[T]) GetByID(id string) (*T, error) {
	var model T
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *baseRepository[T]) FindBy(condition map[string]interface{}) (*T, error) {
	var model T
	if err := r.db.Where(condition).First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *baseRepository[T]) FindAllBy(condition map[string]interface{}) ([]T, error) {
	var models []T
	if err := r.db.Where(condition).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *baseRepository[T]) UpdateByID(id string, updates map[string]interface{}) error {
	return r.db.Model(new(T)).Where("id = ?", id).Updates(updates).Error
}

func (r *baseRepository[T]) DeleteByID(id string) error {
	return r.db.Delete(new(T), "id = ?", id).Error
}

func (r *baseRepository[T]) SoftDeleteByID(id string) error {
	return r.db.Where("id = ?", id).Delete(new(T)).Error
}

func (r *baseRepository[T]) HardDeleteByID(id string) error {
	return r.db.Unscoped().Where("id = ?", id).Delete(new(T)).Error
}

func (r *baseRepository[T]) Exists(condition map[string]interface{}) (bool, error) {
	var count int64
	if err := r.db.Model(new(T)).Where(condition).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *baseRepository[T]) Count(condition map[string]interface{}) (int64, error) {
	var count int64
	if err := r.db.Model(new(T)).Where(condition).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *baseRepository[T]) RawQuery(query string, args ...interface{}) *gorm.DB {
	return r.db.Raw(query, args...)
}
