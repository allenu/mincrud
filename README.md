# mincrud

This is a minimal CRUD app written in go-lang. It uses [minauth](https://github.com/allenu/minauth) for Twitter auth.

The purpose is to illustrate the basics of a CRUD app that has a simple REST API for reading/writing entries and a simple frontend for the UI.

See a demo of it at https://mincrud.appspot.com

# TODOs

- [ ] Sketch out the client API
    - [ ] Clients should be able to provide their own authorizer -- entity_authorizer.go shouldn't need to be in this module
        - [ ] Define the interface for the authorizer
        - [ ] Provide it to NewEntityController as a param
    - [ ] See if there's a way to provide the Entity struct as a param or interface. It should be possible to
          use all the code in core without having a fixed Entity struct that everyone must use.
- [ ] Make it possible to pass in an entity name
- [ ] Support more than one entity type

- [x] Basic Create/Edit
    - [x] Create page with form for creating story

- [x] List all entities
- [ ] Delete entries
- [ ] Users
    - [ ] Show all users
    - [ ] Show all stories per user

