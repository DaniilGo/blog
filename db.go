package main

import (
	"database/sql"
	"fmt"
	"log"
)

func getPosts(db *sql.DB) ([]Post, error) {
	res := make([]Post, 0, 1)

	rows, err := db.Query("select * from posts.habr_posts")
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		post := Post{}

		if err := rows.Scan(&post.Id, &post.Title, &post.Date, &post.Link, &post.Comment); err != nil {
			log.Println(err)
			continue
		}

		res = append(res, post)
	}

	return res, nil
}

// getList — получение списка по id
func getPost(db *sql.DB, id string) (Post, error) {
	row := db.QueryRow(fmt.Sprintf("select * from posts.habr_posts WHERE id = %v", id))

	post := Post{}
	if err := row.Scan(&post.Id, &post.Title, &post.Date, &post.Link, &post.Comment); err != nil {
		return Post{}, err
	}

	return post, nil
}

func editPost(db *sql.DB, post Post, id string) error {
	query := fmt.Sprintf(`UPDATE posts.habr_posts SET title="%s", date="%s", link="%s", comment="%s"  where id=?;`,
		post.Title, post.Date, post.Link, post.Comment)
	_, err := db.Exec(query, id)

	return err
}
