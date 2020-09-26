package repository

import (
	"github.com/sergeychur/avito_auto/internal/models"
)

func (repo *PostgresRepository) InsertLink(link models.Link) (int, models.Link){
	panic("implement me")
}

func (repo *PostgresRepository) GetLink() (int, models.Link) {
	panic("implement me")
}
