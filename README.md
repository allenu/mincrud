# mincrud

This is a minimal CRUD REST API written in go-lang


# TODOs

- [ ] Sketch out the client API
    - [ ] Clients should be able to provide their own authorizer -- entity_authorizer.go shouldn't need to be in this module
        - [ ] Define the interface for the authorizer
        - [ ] Provide it to NewEntityController as a param
    - [ ] See if there's a way to provide the Entity struct as a param or interface. It should be possible to
          use all the code in core without having a fixed Entity struct that everyone must use.

- [ ] Basic Create/Edit
    - [ ] Create page with form for creating story

- [ ] List all entities
- [ ] Make it possible to pass in an entity name
- [ ] Support more than one entity type

