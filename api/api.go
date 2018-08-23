package api

import (
    "io/ioutil"
    "bytes"
    "encoding/json"
    "log"
    "net/http"

    "github.com/allenu/mincrud/core"
    "github.com/allenu/mincrud/core/model"
    userpkg "github.com/allenu/minauth/user"
    "github.com/allenu/minauth/auth"

	"google.golang.org/appengine"
)

const (
    ResponseSuccess = "OK"
    ResponseNoAccess = "No Access"
    ResponseDatabaseError = "Database Error"
    ResponseMissingTitle = "Missing Title"
    ResponseMalformedRequest = "Malformed Request"
    ResponseInvalidParent = "Invalid Parent"
    ResponseNoEntity = "No Entity Found"
)
type ReadEntityRequest struct {
    EntityId string
    AdminMode bool
}

type UpdateEntityRequest struct {
    Entity model.Entity
    Options []string
    AdminMode bool
    RedirectId string
}

type EntityResponse struct {
    ResponseCode string
    Entity model.Entity
}

type ReadUserRequest struct {
    UserId string
}
type ReadUserResponse struct {
    ResponseCode string
    User User
}
type User struct {
    UserId string
    Username string
}

var userStoreProvider userpkg.UserStoreProvider

func init() {
    usp := userpkg.NewUserStoreProvider()
    Setup(usp)
}

func Setup(usp userpkg.UserStoreProvider) {
    userStoreProvider = usp

    log.Println("api.Setup")

    http.HandleFunc("/api/entity/update", apiEntityUpdateHandler)
    http.HandleFunc("/api/entity/read", apiEntityReadHandler)
}

func newEntityController(r *http.Request) core.EntityController {
    ctx := appengine.NewContext(r)
    db := model.NewGoogleEntityDatabase(ctx)
    return core.NewEntityController(db)
}

func apiEntityUpdateHandler(w http.ResponseWriter, r *http.Request) {
    // Read all of the body into a string. This is ugly, but we decode it into the UpdateEntityRequest
    // and into a map[] so that we can read out the properties that were provided.
    defer r.Body.Close()
    body, readErr := ioutil.ReadAll(r.Body)
    if readErr != nil {
        w.Write([]byte("Error reading body"))
        return
    }

    decoder := json.NewDecoder(bytes.NewReader(body))

    var updateEntityRequest UpdateEntityRequest
    err := decoder.Decode(&updateEntityRequest)
    var response EntityResponse
    if err != nil {
        log.Print(err)
        response.ResponseCode = ResponseMalformedRequest
    } else {
        // Read body as an interface{} so we can inspect which properties were provided
        // in Entity.
        var bodyAsMap map[string]interface{}
        unmarshalErr := json.Unmarshal(body, &bodyAsMap)
        if unmarshalErr != nil {
            log.Print(unmarshalErr)
            w.Write([]byte("Error unmarshaling"))
            return
        }
        // Figure out which fields were requested to change in Entity
        var updatedFields []string = []string{}
        if entityMap, ok := bodyAsMap["Entity"].(map[string]interface{}); ok {
            for k := range entityMap {
                updatedFields = append(updatedFields, k)
            }
            log.Printf("updatedFields is %v\n", updatedFields)
        }

        userInfo := auth.GetUserInfo(r)

        sc := newEntityController(r)
        adminMode := false

        updatedEntity, err := sc.Update(userInfo, adminMode, updateEntityRequest.Entity, updatedFields)

        if err == nil {
            response.ResponseCode = ResponseSuccess
            response.Entity = updatedEntity
        } else {
            // Parse err to see which return code
            response.ResponseCode = ResponseNoAccess
        }
    }
    js, _ := json.Marshal(response)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func apiEntityReadHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("api.apiEntityReadHandler")

    decoder := json.NewDecoder(r.Body)

    var readEntityRequest ReadEntityRequest
    err := decoder.Decode(&readEntityRequest)
    var response EntityResponse
    if err != nil {
        response.ResponseCode = ResponseMalformedRequest
    } else {
        userInfo := auth.GetUserInfo(r)
        sc := newEntityController(r)
        adminMode := false

        entity, err := sc.Read(userInfo, adminMode, readEntityRequest.EntityId)
        if err == nil {
            response.ResponseCode = ResponseSuccess
            response.Entity = entity
        } else {
            response.ResponseCode = ResponseNoAccess
        }
    }
    js, _ := json.Marshal(response)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

