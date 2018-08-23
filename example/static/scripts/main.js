
function reduce(oldAppContext, action, payload) {
    var newState = Object.assign({}, oldAppContext.state)

    switch (action) {
        case "load_error":
            newState.mode = "error"
            newState.message = payload
            break;

        case "load_initial_page":
            newState.entityInfo = payload
            newState.mode = oldAppContext.props.initialMode
            newState.history = [newState.entityInfo.entityId]
            newState.readIndex = 0
            break;

        case "cancel_edits":
            newState.message = "Editing cancelled."
            newState.mode = "reading"
            if (oldAppContext.props.redirectOnSave) {
                window.location.replace("/") 
            }
            break

        case "save_edits":
            newState.entityInfo.entity = payload
            newState.entityInfo.entityId = payload.EntityId
            newState.mode = "reading"
            if (oldAppContext.props.redirectOnSave) {
                window.location.replace("/read/" + newState.entityInfo.entityId) 
            } else {
                newState.message = "Saved."
            }
            break

        case "message":
            newState.message = payload
            break;
    }

    const newAppContext = {
        props: oldAppContext.props,
        state: newState
    }
    appContext = newAppContext
    render(appContext)
}

function loadEntityInfo(entityId, userInfo) {
    if (entityId == "") {
        const entityInfo = {
            entityId: "",
            entity: {Title: "", Body: ""},
            contributor: userInfo
        }
        return Promise.resolve(entityInfo)
    } else {
        const url = "/api/entity/read"
        const readRequest = {"EntityId": entityId}
        var fetchEntityPromise = fetch(url, { method: "POST", body: JSON.stringify(readRequest), credentials: "include" })
            .then( response => response.json() )
            .then( apiResponse => {
                if (apiResponse.ResponseCode == "OK") {
                    return Promise.resolve(apiResponse.Entity)
                } else {
                    return Promise.reject(new Error(errorLoadingEntity))
                }
            }).then( entity => {
                var authorId = entity.AuthorId
                const userUrl = "/api/user/" + authorId
                return fetch(userUrl, { method: "GET", credentials: "include" })
                    .then( response => response.json() )
                    .then( apiResponse => {
                        if (apiResponse.ResponseCode == "OK") {
                            return Promise.resolve({entity: entity, userInfo: { userId: apiResponse.User.UserId, username: apiResponse.User.Username }})
                        } else {
                            return Promise.resolve({entity: entity, userInfo: { userId: "anonymous", username: "anonymous" }})
                            // return Promise.reject(new Error(errorLoadingUser))
                        }
                    })
            })

        // TODO: Fetch userInfo 
        //
        return fetchEntityPromise.then( (entityAndUserInfo) => {
            const entityInfo = {
                entityId: entityId,
                entity: entityAndUserInfo.entity,
                contributor: entityAndUserInfo.userInfo
            }
            console.log(entityInfo)
            return Promise.resolve(entityInfo)
        })
    }
}

function initializePage(oldAppContext) {
    loadEntityInfo(oldAppContext.props.initialEntityId, oldAppContext.props.userInfo)
        .then( (entityInfo) => {
            if (oldAppContext.props.initialMode == "editing" && entityInfo.entityId != "" && !canEditEntity(entityInfo.entity, oldAppContext.props.userInfo)) {
                reduce(oldAppContext, "load_error", "Sorry, you don't have access to edit this entity.")
            } else {
                reduce(oldAppContext, "load_initial_page", entityInfo)
            }
        }).catch( (err) => {
            reduce(oldAppContext, "load_error", "Error loading initial entity: " + err)
        })
}

function canEditEntity(entity, userInfo) {
    const entityHasNoBody = entity.Body == ""
    return (entity.AuthorId == userInfo.userId)
}

function editEntity() {
    if (appContext.state.entityInfo.entityId != "" && !canEditEntity(appContext.state.entityInfo.entity, appContext.props.userInfo)) {
        reduce(appContext, "message", "Sorry, you don't have access to edit this entity.")
        return
    }

    var newState = Object.assign({}, appContext.state)
    newState.message = ""
    newState.mode = "editing"
    const newAppContext = {
        props: appContext.props,
        state: newState
    }
    appContext = newAppContext
    render(appContext)
}

function cancelEdits() {
    reduce(appContext, "cancel_edits", null)
}

function saveEdits() {
    const title = document.getElementById("Title").value;
    const body = document.getElementById("Body").value;
    if (title == "" || (body == "" && appContext.state.entityInfo.Body != "")) {
        reduce(appContext, "message", "Cannot save with no body or no title")
        return
    }

    console.log("saveEdits()")

    var updatedEntity = {"EntityId": appContext.state.entityInfo.entityId, "Title": title, "Body": body}
    const updateRequest = {"Entity":updatedEntity}
    const updateUrl = "/api/entity/update"
    fetch(updateUrl, {
        method: "POST",
        body: JSON.stringify(updateRequest), 
        credentials: "include"
        })
        .then( response => response.json())
        .then( serverResponse => {
            if (serverResponse.ResponseCode == "OK") {
                reduce(appContext, "save_edits", serverResponse.Entity)
            } else {
                reduce(appContext, "message", "Error saving")
            }
        })
}

// -------------------------------------------------------------------------------- 

function FriendlyTimeAgo(date) {
    let today = new Date()
    let post_date = new Date(date)
    let ms_per_day = 1000 * 60 * 60 * 24
    let ms_ago = Math.abs(today - post_date)
    let days_ago = Math.floor(ms_ago / ms_per_day)
    if (days_ago <= 0) {
        if (ms_ago > 1000*60*60) {
            let h_ago = Math.floor(ms_ago / (1000 * 60*60))
            days_ago = h_ago + "h ago"
        } else if (ms_ago > 1000*60*5) {
            let mins_ago = Math.floor(ms_ago / (1000 * 60))
            days_ago = mins_ago + "m ago"
        } else {
            days_ago = "moments ago"
        }
    } else {
        days_ago = days_ago + "d ago"
    }

    return days_ago
}

// Via https://stackoverflow.com/questions/6234773/can-i-escape-html-special-chars-in-javascript
function escapeHtml(unsafe) {
    return unsafe
         .replace(/&/g, "&amp;")
         .replace(/</g, "&lt;")
         .replace(/>/g, "&gt;")
         .replace(/"/g, "&quot;")
         .replace(/'/g, "&#039;");
 }

