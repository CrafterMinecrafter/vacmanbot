package usermanager

import (
	"github.com/nyan2d/bolteo"
	"github.com/nyan2d/vacmanbot/app/models"
	"gopkg.in/tucnak/telebot.v2"
)

type UserManager struct {
	cache map[int]*models.User
	db    *bolteo.Bolteo
}

func NewUserManager(db *bolteo.Bolteo) *UserManager {
	return &UserManager{
		cache: make(map[int]*models.User),
		db:    db,
	}
}

func (um *UserManager) Check(usr telebot.User) {
	user, ok := um.cache[usr.ID]
	if !ok {
		if !um.cacheUserFromDatabase(usr.ID) {
			u := models.User{
				ID:        usr.ID,
				FirstName: usr.FirstName,
				LastName:  usr.LastName,
				Username:  usr.Username,
				IsAdmin:   false,
				IsIgnored: false,
			}
			um.cache[u.ID] = &u
			um.saveToDatabaseFromCache(u.ID)
			return
		} else {
			user = um.cache[usr.ID]
		}
	}
	if user.FirstName != usr.FirstName || user.LastName != usr.LastName || user.Username != usr.Username {
		u := models.User{
			ID:        usr.ID,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
			Username:  usr.Username,
			IsAdmin:   false,
			IsIgnored: false,
		}
		um.cache[u.ID] = &u
		um.saveToDatabaseFromCache(u.ID)
	}
}

func (um *UserManager) GetUser(id int) models.User {
	user, ok := um.cache[id]
	if !ok {
		if um.cacheUserFromDatabase(id) {
			return *um.cache[id]
		}
		return models.User{
			ID:        id,
			FirstName: "Undefined",
			LastName:  "Undefined",
		}
	}
	return *user
}

func (um *UserManager) SetUser(u models.User) {
	um.db.Bucket("users")
	um.cache[u.ID] = &u
	um.db.Put(u.ID, u)
}

func (um *UserManager) GetNames(id int) string {
	usr := um.GetUser(id)
	if usr.LastName == "" {
		return usr.FirstName
	}
	return usr.FirstName + " " + usr.LastName
}

func (um *UserManager) cacheUserFromDatabase(id int) bool {
	_, ok := um.cache[id]
	if !ok {
		um.db.Bucket("users")
		user := models.User{}
		if um.db.Get(id, &user) != nil {
			return false
		} else {
			um.cache[id] = &user
			return true
		}
	}
	return true
}

func (um *UserManager) saveToDatabaseFromCache(id int) {
	if user, ok := um.cache[id]; ok {
		um.db.Bucket("users")
		um.db.Put(user.ID, *user)
	}
}
