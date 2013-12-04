package main

import (
	"time"
)

type Submission struct {
	Id      int64
	Created int64
	Title   string
	Author  string
	Code    string
}

func (s *Submission) insert() error {
	return dbmap.Insert(s)
}

func (s *Submission) update() (int64, error) {
	return dbmap.Update(s)
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
	if err != nil || obj == nil {
		return nil, err
	}
	return obj.(*Submission), nil
}
