package repository

import "github.com/sergeychur/avito_auto/internal/models"

type Repository interface {
	InsertLink(link models.Link) (int, models.Link)
	GetLink() (int, models.Link)
	Start(maxConn, acquireTimeout int) error
	Close()
}
