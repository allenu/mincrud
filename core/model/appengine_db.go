package model

import (
    "golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type GoogleEntityDatabase struct {
    appEngineContext context.Context
}

const entityName = "MinCrudEntry"

// The App Engine code only cares about the following fields:
// - AuthorId: UUID of the user that created this entry
// - Created: creation date

func NewGoogleEntityDatabase(appEngineContext context.Context) EntityDatabase {
    return GoogleEntityDatabase{appEngineContext: appEngineContext}
}

func (db GoogleEntityDatabase) Read(entityId string) (Entity, error) {
    key := datastore.NewKey(db.appEngineContext, entityName, entityId, 0, nil)

    var entity Entity
    err := datastore.Get(db.appEngineContext, key, &entity)
    if err != nil {
        return entity, err
    }
    return entity, nil
}

func (db GoogleEntityDatabase) Update(entity Entity) error {
    key := datastore.NewKey(db.appEngineContext, entityName, entity.EntityId, 0, nil)
    _, err := datastore.Put(db.appEngineContext, key, &entity)

    return err
}

// --------------------------------------------------------------------------------
// Other

func (db GoogleEntityDatabase) ReadEntities() ([]Entity, error) {
    const maxChildren = 50
    q := datastore.NewQuery(entityName).Order("Created").Limit(maxChildren)

    children := make([]Entity, 0, maxChildren)
    if _, err := q.GetAll(db.appEngineContext, &children); err == nil {
        return children, nil
    } else {
        return children, err
    }
}

func (db GoogleEntityDatabase) ReadEntitiesForUser(userId string) ([]Entity, error) {
    const maxChildren = 50
    q := datastore.NewQuery(entityName).Order("Created").Filter("AuthorId =", userId).Limit(maxChildren)

    children := make([]Entity, 0, maxChildren)
    if _, err := q.GetAll(db.appEngineContext, &children); err == nil {
        return children, nil
    } else {
        return children, err
    }
}

