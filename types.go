package main

import "html/template"

type BlogPost struct {
	Content template.HTML
}

type BlogPosts struct {
	Posts []string
}
