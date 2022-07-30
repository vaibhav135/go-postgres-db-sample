package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string, db *sql.DB) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album
	query := fmt.Sprintf("SELECT * FROM album WHERE artist='%s'", name)
	rows, err := db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func albumByID(id int64, db *sql.DB) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album
	fmt.Println(id)
	row := db.QueryRow("SELECT * FROM album WHERE id=$1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// albumByID queries for the album with the specified ID.
func albumDeletedByID(id int64, db *sql.DB) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album
	fmt.Println(id)
	row := db.QueryRow("DELETE FROM album WHERE id >= $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsDeletedById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsDeletedById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album, db *sql.DB) (Album, error) {
	//var album Album
	album := db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id, title, artist, price", alb.Title, alb.Artist, alb.Price)
	var result Album
	if err := album.Scan(&result.ID, &result.Title, &result.Artist, &result.Price); err != nil {
		return result, fmt.Errorf("addAlbum: %v", err)
	}
	fmt.Println(result)
	return result, nil
}

func main() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Error loading env")
	}

	cfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"),
	)
	db, err := sql.Open("pgx", cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	albums, err := albumsByArtist("John Coltrane", db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	album, err := albumByID(albums[0].ID, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", album)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)

	albumDeleted, err := albumDeletedByID(6, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album deleted successfully: %v\n", albumDeleted)

	//defer db.Close()

}
