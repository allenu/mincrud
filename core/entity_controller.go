package core

import (
    "errors"
    "log"
    "github.com/allenu/mincrud/core/model"
)

type EntityController struct {
    authorizer EntityAuthorizer
    db model.EntityDatabase
}

func NewEntityController(db model.EntityDatabase) EntityController {
    return EntityController{authorizer: EntityAuthorizer{}, db: db}
}

func (sc EntityController) Update(user User, adminMode bool, modifiedEntity model.Entity, modifiedFields []string) (model.Entity, error) {
    var storedEntity model.Entity

    // Read entity from database, if it exists, if it doesn't exist, treat it as an insert operation...
    if modifiedEntity.EntityId != "" {
        var err error
        storedEntity, err = sc.db.Read(modifiedEntity.EntityId)
        if err != nil {
            log.Printf("Entity doesn't exist at that id %v err: %v", modifiedEntity.EntityId, err)

            // Cannot update it because couldn't read entity
            return storedEntity, err
        }

        if !sc.authorizer.CanRead(storedEntity, user, adminMode) {
            return storedEntity, errors.New("Unauthorized read")
        }
    } else {
        // Entity is a new one, so create a blank one
        log.Print("Creating new entry")
        storedEntity = model.NewEntity(user.GetUserId())

        // See if user is allowed to write at all
        if !sc.authorizer.CanWrite(user, adminMode, storedEntity) {
            return storedEntity, errors.New("User has no write access for new stories")
        }
    }

    updatableFields := sc.authorizer.UpdatableFields(storedEntity, user, adminMode)

    log.Printf("You are allowed to edit these fields: %v", updatableFields)
    log.Printf("Requested fields: %v", modifiedFields)
    // Find intersection of updatableFields and the ones requested
    //actualFieldsToUpdate := intersection(updatableFields, modifiedFields)
    actualFieldsToUpdate := Intersect(updatableFields, modifiedFields)

    log.Printf("The fields that actually will be edited are %v", actualFieldsToUpdate)

    updatedEntity := sc.authorizer.ApplyUpdates(storedEntity, actualFieldsToUpdate, modifiedEntity)
    err := sc.db.Update(updatedEntity)

    if err != nil {
        log.Printf("Error UpdateEntity: %v", err)
        return storedEntity, err
    }

    return updatedEntity, nil
}

func (sc EntityController) Read(user User, adminMode bool, entityId string) (model.Entity, error) {
    entity, err := sc.db.Read(entityId)
    if err != nil {
        return entity, err
    }

    if !sc.authorizer.CanRead(entity, user, adminMode) {
        return entity, errors.New("Not authorized to read object")
    }

    return entity, nil
}

func (sc EntityController) ReadChildren(user User, adminMode bool) ([]model.Entity, error) {
    var filteredStories []model.Entity = []model.Entity{}
    var children []model.Entity
    var err error

    children, err = sc.db.ReadChildren()
    if err == nil {
        for _, entity := range children {
            if sc.authorizer.CanRead(entity, user, adminMode) {
                filteredStories = append(filteredStories, entity)
            }
        }
    }
    return filteredStories, err
}

func (sc EntityController) ReadChildrenForUser(user User, adminMode bool, userId string) ([]model.Entity, error) {
    var filteredStories []model.Entity = []model.Entity{}
    var children []model.Entity
    var err error

    children, err = sc.db.ReadChildrenForUser(userId)
    if err == nil {
        for _, entity := range children {
            if sc.authorizer.CanRead(entity, user, adminMode) {
                filteredStories = append(filteredStories, entity)
            }
        }
    }
    return filteredStories, err
}
