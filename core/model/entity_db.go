package model

type EntityDatabase interface {
    Read(entityId string) (Entity, error)
    ReadEntities() ([]Entity, error)
    ReadEntitiesForUser(userId string) ([]Entity, error)
    Update(entity Entity) error
}

