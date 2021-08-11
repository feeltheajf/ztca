package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	Validate *validator.Validate
)

func init() {
	Validate = validator.New()
	Validate.SetTagName("binding")
}

// Model is a base model for db interaction
type Model struct {
	ID        uuid.UUID  `json:"-" gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

// BeforeCreate implements UUID generation for new resources
func (*Model) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New())
}

// service is a convenience struct for grouping methods
type service struct{}
