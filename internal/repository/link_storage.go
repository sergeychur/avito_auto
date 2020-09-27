package repository

import (
	"github.com/sergeychur/avito_auto/internal/models"
	"gopkg.in/jackc/pgx.v2"
	"log"
)

const (
	InsertLink = "INSERT INTO shortcuts (long_url, short_code) VALUES ($1,$2)"
	CheckIfNameExists = "SELECT EXISTS(SELECT short_code FROM shortcuts WHERE short_code = $1)"
	GetLongLink = "SELECT long_url FROM shortcuts WHERE short_code = $1"
)

func (repo *PostgresRepository) InsertLink(link models.Link) (int, models.Link) {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Println("tx error: ", err)
		return DB_ERROR, models.Link{}
	}
	defer func() {
		_ = tx.Rollback()
	}()
	if link.Shortcut == "" {
		link.Shortcut = GenerateShortcut(link)
	} else {
		exists := false
		err := tx.QueryRow(CheckIfNameExists, link.Shortcut).Scan(&exists)
		if err != nil {
			log.Println("CheckNameExists error: ", err)
			return DB_ERROR, models.Link{}
		}

		if exists {
			return FORBIDDEN, models.Link{}
		}
	}
	_, err = tx.Exec(InsertLink, link.RealURL, link.Shortcut)
	if err != nil {
		log.Println("InsertLink error: ", err)
	}
	err = tx.Commit()
	if err != nil {
		log.Println("insert tx commit error: ", err)
	}
	return CREATED, link
}

func (repo *PostgresRepository) GetLink(shortcut string) (int, string) {
	longLink := ""
	err := repo.db.QueryRow(GetLongLink, shortcut).Scan(&longLink)
	if err == pgx.ErrNoRows {
		return EMPTY_RESULT, ""
	}
	return OK, longLink
}
