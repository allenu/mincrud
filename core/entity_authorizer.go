package core

import (
    "github.com/allenu/mincrud/core/model"
)

type EntityAuthorizer struct {
}

func (t EntityAuthorizer) CanRead(entity model.Entity, user User, adminMode bool) bool {
    if user.GetUserId() == "banned" {
        return false
    }

    return !entity.IsPrivate || user.GetUserId() == entity.AuthorId || user.GetUserId() == "reader"
}

func (t EntityAuthorizer) CanWrite(user User, adminMode bool, entity model.Entity) bool {
    /*
    if user.GetUserId() == "anonymous" || user.GetUserId() == "banned" {
        return false
    }
    */

    return user.GetUserId() == entity.AuthorId
}

func (t EntityAuthorizer) UpdatableFields(entity model.Entity, user User, adminMode bool) []string {
    if user.GetUserId() == entity.AuthorId {
        return []string{"Title", "Body", "IsPrivate"}
    } else {
        return []string{}
    }
}

func (t EntityAuthorizer) ApplyUpdates(entity model.Entity, updatedFields []string, updates model.Entity) model.Entity {
    if Contains(updatedFields, "IsPrivate") {
        entity.IsPrivate = updates.IsPrivate
    }
    if Contains(updatedFields, "Title") {
        entity.Title = updates.Title
    }
    if Contains(updatedFields, "Body") {
        entity.Body = updates.Body
    }

    return entity
}

