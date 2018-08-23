package model

type EntityDatabase interface {
    Read(entityId string) (Entity, error)
    ReadChildren() ([]Entity, error)
    ReadChildrenForUser(userId string) ([]Entity, error)
    Update(entity Entity) error
}

