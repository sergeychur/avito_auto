package repository

import (
	"github.com/rs/xid"
	"github.com/sergeychur/avito_auto/internal/models"
)

func GenerateShortcut(link models.Link) string {
	id := xid.New()
	return id.String()
}
