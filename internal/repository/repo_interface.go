package repository

import "github.com/sergeychur/avito_auto/internal/models"

type Repository interface {
	InsertLink(link models.Link) (int, models.Link)
	GetLink(shortcut string) (int, string)
	Start(maxConn, acquireTimeout int) error
	Close()
}
