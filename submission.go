package main

import (
	"log"
	"time"
)

type Submission struct {
	Id      int64
	Created int64
	Title   string
	Author  string
	Code    string
}

func (s *Submission) Insert() {
	err := dbmap.Insert(&s)
	if err != nil {
		log.Fatalln("Insert failed", err)
	}
}

func (s *Submission) Update() int64 {
	count, err := dbmap.Update(&s)
	if err != nil {
		log.Fatalln("Insert failed", err)
	}
	return count
}

func newSubmission(title, author, code string) *Submission {
	return &Submission{
		Created: time.Now().UnixNano(),
		Title:   title,
		Author:  author,
		Code:    code,
	}
}

func loadSubmission(id int64) (*Submission, error) {
	obj, err := dbmap.Get(Submission{}, id)
	if err != nil {
		return nil, err
	}
	return obj.(*Submission), nil
}
