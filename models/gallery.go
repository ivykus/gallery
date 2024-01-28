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

func (gs *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := gs.DB.Query(`
		SELECT id, title
		FROM galleries
		WHERE user_id = $1;
	`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("gallery query by user id: %w", err)
	}
	defer rows.Close()

	var galleries []Gallery
	for rows.Next() {
		var gallery Gallery
		if err := rows.Scan(&gallery.ID, &gallery.Title); err != nil {
			return nil, fmt.Errorf("gallery query by user_id: %w", err)
		}
		galleries = append(galleries, gallery)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gallery query by user_id: %w", err)
	}
	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	_, err := gs.DB.Exec(`
		UPDATE galleries
		SET title = $1
		WHERE id = $2;
	`, gallery.Title, gallery.ID)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Delete(galleryID int) error {
	_, err := gs.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;
	`, galleryID)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}
	return nil
}
