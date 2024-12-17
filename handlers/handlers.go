package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"music_library/models"
	"net/http"
	"strconv"
	"strings"
)


// GetSongs godoc
// @Summary      Get songs with filtering and pagination
// @Description  Retrieves a list of songs based on filters and pagination parameters.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        group_name query string false "Filter by group name"
// @Param        song_name query string false "Filter by song name"
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Number of items per page" default(10)
// @Success      200  {object}  []models.Song
// @Failure      400  {string}  string "Invalid request parameters"
// @Failure      500  {string}  string "Internal server error"
// @Router       /songs [get]
func GetList(rw http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	query := "SELECT id, group_name, song_name, release_date, text, link FROM Songs WHERE 1=1"
	if group != "" {
		query += " AND group_name = '" + group + "'"
	}
	query += " LIMIT " + limit + " OFFSET " + offset

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=music sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Error executing query:", err)
	}

	defer rows.Close()

	var songs []models.Song

	for rows.Next() {
		var song models.Song
		err = rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); 
		
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		songs = append(songs, song)
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(songs)
}

// AddSong godoc
// @Summary      Add a new song
// @Description  Adds a new song to the library.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        song body models.Song true "Song data"
// @Success      201  {string}  string "Song added successfully"
// @Failure      400  {string}  string "Invalid JSON request"
// @Failure      500  {string}  string "Internal server error"
// @Router       /song [post]
func AddSong(rw http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song); 
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=music sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO songs (group_name, song_name) VALUES ($1, $2)", song.Group, song.Song)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(song)
}


// GetSongTextWithPagination godoc
// @Summary      Get song text with pagination
// @Description  Retrieves the text of a song split into verses, with pagination.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id query int true "ID of the song"
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Verses per page" default(2)
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {string}  string "Invalid request parameters"
// @Failure      404  {string}  string "Song not found"
// @Failure      500  {string}  string "Internal server error"
// @Router       /song-text [get]
func GetSongTextWithPagination(rw http.ResponseWriter, r *http.Request) {
	songID := r.URL.Query().Get("id") // ID песни
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1 
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 2 
	}

	if songID == "" {
		http.Error(rw, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=music sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var songText string
	err = db.QueryRow("SELECT text FROM songs WHERE id = $1", songID).Scan(&songText)
	if err != nil {
		http.Error(rw, "Song not found", http.StatusNotFound)
		return
	}

	verses := strings.Split(songText, "\n\n")
	totalVerses := len(verses)                                

	totalPages := int(math.Ceil(float64(totalVerses) / float64(limit)))
	start := (page - 1) * limit                              // Индекс первого куплета
	end := start + limit                                     // Индекс последнего куплета

	if start >= totalVerses {
		http.Error(rw, "Page number out of range", http.StatusBadRequest)
		return
	}

	if end > totalVerses {
		end = totalVerses
	}

	response := map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"verses":      verses[start:end],
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

// DeleteSong godoc
// @Summary      Delete a song
// @Description  Deletes a song by its ID.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id query int true "ID of the song to delete"
// @Success      200  {string}  string "Song deleted successfully"
// @Failure      400  {string}  string "Invalid 'id' parameter"
// @Failure      404  {string}  string "Song not found"
// @Failure      500  {string}  string "Internal server error"
// @Router       /song [delete]
func DeleteSong(rw http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=musiс sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	songID := r.URL.Query().Get("id")
	if songID == "" {
		http.Error(rw, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(songID)
	if err != nil {
		http.Error(rw, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		http.Error(rw, "Failed to delete song", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(rw, "Error checking rows affected", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		// Если песня с таким ID не найдена
		http.Error(rw, "Song not found", http.StatusNotFound)
		return
	}

	// Формируем ответ
	response := map[string]string{"message": "Song deleted successfully"}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}



// UpdateSong godoc
// @Summary      Update a song
// @Description  Updates details of a song by its ID. Fields are optional.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id query int true "ID of the song to update"
// @Param        song body models.SongUpdateRequest true "Fields to update"
// @Success      200  {string}  string "Song updated successfully"
// @Failure      400  {string}  string "Invalid request parameters"
// @Failure      404  {string}  string "Song not found"
// @Failure      500  {string}  string "Internal server error"
// @Router       /song [put]
func UpdateSong(rw http.ResponseWriter, r *http.Request) {
	songID := r.URL.Query().Get("id")
	if songID == "" {
		http.Error(rw, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(songID)
	if err != nil {
		http.Error(rw, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}
	var req models.SongUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(rw, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	query := "UPDATE songs SET"
	params := []interface{}{}
	paramCount := 1

	if req.GroupName != nil {
		query += " group_name = $" + strconv.Itoa(paramCount) + ","
		params = append(params, *req.GroupName)
		paramCount++
	}
	if req.SongName != nil {
		query += " song_name = $" + strconv.Itoa(paramCount) + ","
		params = append(params, *req.SongName)
		paramCount++
	}
	if req.ReleaseDate != nil {
		query += " release_date = $" + strconv.Itoa(paramCount) + ","
		params = append(params, *req.ReleaseDate)
		paramCount++
	}
	if req.Text != nil {
		query += " text = $" + strconv.Itoa(paramCount) + ","
		params = append(params, *req.Text)
		paramCount++
	}
	if req.Link != nil {
		query += " link = $" + strconv.Itoa(paramCount) + ","
		params = append(params, *req.Link)
		paramCount++
	}

	if len(params) == 0 {
		http.Error(rw, "No fields to update", http.StatusBadRequest)
		return
	}

	query = query[:len(query)-1] // Удаляем последнюю запятую
	query += " WHERE id = $" + strconv.Itoa(paramCount)
	params = append(params, id)


	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=music sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	result, err := db.Exec(query, params...)
	if err != nil {
		http.Error(rw, "Failed to update song", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(rw, "Error checking rows affected", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(rw, "Song not found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"message": "Song updated successfully"})
}
