package server

import (
	"github.com/sergeychur/avito_auto/internal/models"
)

type IValidator interface {
	ValidateLink(link models.Link) error
}
