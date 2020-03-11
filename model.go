package main

import (
	"fmt"
	"log"
)

const DSN = "root:1234@tcp(localhost:3306)/task_list_app?charset=utf8"

type Post struct {
	Id      int
	Title   string
	Date    string
	Link    string
	Comment string
}

/*
CREATE TABLE `posts`.`habr_posts` (
`id` INT NOT NULL AUTO_INCREMENT,
`title` TEXT NOT NULL,
`date` TEXT NULL,
`link` TEXT NULL,
`comment` TEXT NULL,
PRIMARY KEY (`id`),
UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE);
*/

const (
	databaseName = "posts"
	tableName    = "habr_posts"
)

func (s Server) insertDefault() {
	for _, post := range createPosts() {
		query := fmt.Sprintf(
			`insert into %s.%s (id,title,date,link,comment) values (?,?,?,?,?);`, databaseName, tableName)
		_, err := s.db.Exec(query, post.Id, post.Title, post.Date, post.Link, post.Comment)
		if err != nil {
			log.Println(err)
			return
		}
	}
	log.Print("inserted default")
}

func createPosts() []Post {
	return []Post{
		{
			Id:      1,
			Title:   "Hello world! Or Habr in English, v1.0",
			Date:    "15.01.2019 in 14:15",
			Link:    "https://habr.com/ru/company/habr/blog/435764/",
			Comment: "Nice one!",
		},
		{
			Id:      2,
			Title:   "Common Errors in English Usage",
			Date:    "29.06.2010 in 19:51",
			Link:    "https://habr.com/ru/post/97778/",
			Comment: "test comment",
		},
	}
}

func (s Server) truncate() {
	query := fmt.Sprintf("truncate %s.%s;", databaseName, tableName)
	_, err := s.db.Exec(query)
	if err != nil {
		log.Println(err)
	}
	log.Print("truncate table")
}

/*
insert into posts.habr_posts (title,date,link,comment)
 values ("First List", "First Description"), ("Second List", "Second Description");
*/
