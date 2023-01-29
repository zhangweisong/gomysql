package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"model"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/recordings")
	if err != nil {
		log.Fatal(err)
	}

	list, err := AlbumByID(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(list)

}

func AddAlbum(alb model.Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist,price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}
	return id, nil
}

func canPurchase(id int, price int) (bool, error) {
	var enough bool
	// 基于单行查询值. 因为查询的只有一个字段，最终结果是1*1 的 放入一个变量即可。
	if err := db.QueryRow("SELECT (price >= ?) from album where id = ?", price, id).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("canPurchase %d: unknown album", id)
		}
		return false, fmt.Errorf("canPurchase %d: %v", id)
	}
	return enough, nil
}

func albumsByArtist(artist string) ([]model.Album, error) {
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 用于保存返回行中的数据的专辑切片
	var albums []model.Album

	// 遍历行，使用Scan将列数据分配给结构字段。
	for rows.Next() {
		var alb model.Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist,
			&alb.Price); err != nil {
			return albums, err
		}
		albums = append(albums, alb)
	}
	if err = rows.Err(); err != nil {
		return albums, err
	}
	return albums, nil
}

// AlbumByID 检索指定的专辑.
func AlbumByID(id int) (model.Album, error) {
	// 定义预处理语句。您通常会在其他位置定义语句，
	// 并将其保存以用于诸如此函数的函数中
	stmt, err := db.Prepare("SELECT * FROM album WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	var album model.Album

	// 执行预准备语句，为占位符为 ？ 的
	// 参数传入 id 值
	err = stmt.QueryRow(id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			// 处理未返回行的情况.
		}
		return album, err
	}
	return album, nil
}
