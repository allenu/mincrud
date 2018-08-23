package model

import (
    "time"
    "github.com/google/uuid"
)

// --------------------------------------------------------------------------------
// Entity

type Entity struct {
    Created         time.Time
    LastModified    time.Time

    EntityId     string
    AuthorId    string

    IsPrivate   bool
    Title       string
    Body        string
}

// authdb.DatabaseObject

func (s Entity) DbObjectId() string {
    return s.EntityId
}

func NewEntity(authorId string) Entity {
    return Entity{
        Created: time.Now(),
        LastModified: time.Now(),
        EntityId: uuid.New().String(),
        AuthorId: authorId,
        IsPrivate: false,
        Title: "",
        Body: "",
    }
}

