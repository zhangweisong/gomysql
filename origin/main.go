package origin

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"model"
)

type DB struct {
	*sql.DB
}

func Connect(dbname string) (*DB, error) {
	// 捕获连接属性.
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               dbname,
		AllowNativePasswords: true,
	}
	// 获取数据库句柄.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	return &DB{db}, nil
}

func (db DB) Select(sql string, variables ...any) ([]model.Album, error) {
	var albums []model.Album

	rows, err := db.Query(sql, variables...)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist: %v", err)
	}
	defer rows.Close()
	// 遍历行，使用"扫描"将列数据分配给结构字段.
	for rows.Next() {
		var alb model.Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist: %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist: %v", err)
	}
	return albums, nil
}

// albumByID查询指定ID的专辑。
func (db DB) ItemByID(id int64) (model.Album, error) {
	// 一个专辑将保存返回行的数据。
	var alb model.Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

func (db DB) Add(alb model.Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
