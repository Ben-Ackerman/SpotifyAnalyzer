package app

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// Database represents a database connection used to store music tracks
type Database interface {
	Close()
	InsertTrack(*Track) (string, error)
	GetTrackID(string, string) (string, error)
	GetTrack(string) (*Track, error)
}

// PostgresDB is a concret type for the interface Database
// Which uses postgres as a DB.  To use call InitPostgresDB()
// and when fininish must call close()
type PostgresDB struct {
	db *sql.DB
}

// Close closes the *sqlDB connection
func (p *PostgresDB) Close() {
	p.db.Close()
}

// InitPostgresDB creates a PostgresDB struct and initializes the
// nessisary fields.
func InitPostgresDB() (*PostgresDB, error) {
	host := "localhost"
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DB")
	password := os.Getenv("POSTGRES_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

// InsertTrack inserts a track into the database
func (p *PostgresDB) InsertTrack(track *Track) (string, error) {
	//Error check inputs and sanitize inputs
	if track == nil {
		return "", fmt.Errorf("Input to InsertTrack cannot be nil")
	}
	if track.Artist == "" {
		return "", fmt.Errorf("InsertTrack: track.artist cannot be empty")
	}
	if track.Name == "" {
		return "", fmt.Errorf("InsertTrack: track.name cannot be empty")
	}

	track.Artist = strings.ToLower(track.Artist)
	track.Name = strings.ToLower(track.Name)
	sqlStatement := `INSERT INTO info.tracks (track, artist, lyrics, geniusURI)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	id := ""
	err := p.db.QueryRow(sqlStatement,
		track.Name,
		track.Artist,
		track.Lyrics,
		track.GeniusURI).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetTrackID searches the database for a track id that relates to the
// given name and artist.  Must be an exact match.
func (p *PostgresDB) GetTrackID(nameIn string, artistIn string) (string, error) {
	//Error checking and input sanitation
	if artistIn == "" {
		return "", fmt.Errorf("GetTrack: artistIn cannot be empty")
	}
	if nameIn == "" {
		return "", fmt.Errorf("GetTrack: nameIn cannot be empty")
	}

	nameIn = strings.ToLower(nameIn)
	artistIn = strings.ToLower(artistIn)

	sqlStatement := `SELECT id FROM info.tracks
	WHERE artist = $1 AND track = $2
	LIMIT 1`
	var id string
	row := p.db.QueryRow(sqlStatement, artistIn, nameIn)
	err := row.Scan(&id)
	switch err {
	case sql.ErrNoRows:
		return "", nil
	case nil:
		return id, nil
	default:
		return "", err
	}

}

// GetTrack returns a track object whose fields come from the
// database entry matching the input id.
func (p *PostgresDB) GetTrack(id string) (*Track, error) {
	//Error checking and input sanitation
	if id == "" {
		return nil, fmt.Errorf("GetTrack: id cannot be empty")
	}

	sqlStatement := `SELECT track, artist, lyrics, geniusURI FROM info.tracks
	WHERE id = $1
	LIMIT 1`
	var name, artist, lyrics, geniusURI string
	row := p.db.QueryRow(sqlStatement, id)
	err := row.Scan(&name, &artist, &lyrics, &geniusURI)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &Track{
			ID:        id,
			Artist:    artist,
			Name:      name,
			Lyrics:    lyrics,
			GeniusURI: geniusURI,
		}, nil
	default:
		return nil, err
	}

}
