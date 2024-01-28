package models

import (
	"database/sql"
	"fmt"

	"github.com/ivykus/gallery/errors"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB
}

func (gs *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}

	raw := gs.DB.QueryRow(`
		INSERT INTO galleries (title, user_id)
		VALUES ($1, $2) RETURNING id;`,
		gallery.Title,
		gallery.UserID)

	if err := raw.Scan(&gallery.ID); err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByID(galleryID int) (*Gallery, error) {
	gallery := Gallery{
		ID: galleryID,
	}
	row := gs.DB.QueryRow(`
		SELECT title, user_id
		FROM galleries
		WHERE id = $1;
	`, galleryID)

	if err := row.Scan(&gallery.Title, &gallery.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("gallery query by id: %w", err)
	}
	return &gallery, nil
}
