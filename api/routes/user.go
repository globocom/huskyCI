package routes

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/security"
	"github.com/globocom/huskyCI/api/types"
	"github.com/labstack/echo"
)

// UpdateUser edits an user
func UpdateUser(c echo.Context) error {

	// step 1.1: valid JSON?
	attemptUser := types.User{}
	err := c.Bind(&attemptUser)
	if err != nil {
		log.Error("EditUser", "USER", 1024, err)
		reply := map[string]interface{}{"success": false, "error": "invalid user JSON"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	// step 1.2: check mongoDB injection

	// step 2.1: password/user is empty?
	if attemptUser.Password == "" || attemptUser.Name == "" || attemptUser.NewPassword == "" {
		reply := map[string]interface{}{"success": false, "error": "passwords/username can not be empty"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	// step 2.2: passwords match?
	if attemptUser.NewPassword != attemptUser.ConfirmNewPassword {
		reply := map[string]interface{}{"success": false, "error": "passwords do not match"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	// step 3: user exists?
	userQuery := map[string]interface{}{"username": attemptUser.Name}
	user, err := db.FindOneDBUser(userQuery)
	if err != nil {
		if err == mgo.ErrNotFound {
			reply := map[string]interface{}{"success": false, "error": "user not found"}
			return c.JSON(http.StatusNotFound, reply)
		}
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}

	// step 4: password is correct?
	if ok := security.CheckPasswordHash(attemptUser.Password, user.HashedPassword); !ok {
		reply := map[string]interface{}{"success": false, "error": "unauthorized"}
		return c.JSON(http.StatusUnauthorized, reply)
	}

	// step 5.1: prepare new user struct to be updated
	newHashedPassword, err := security.BcryptPassword(attemptUser.NewPassword)
	if err != nil {
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}

	updatedUser := types.User{
		Name:           attemptUser.Name,
		HashedPassword: newHashedPassword,
	}

	// step 5.2: update user
	if err := db.UpdateOneDBUser(userQuery, updatedUser); err != nil {
		if err == mgo.ErrNotFound {
			reply := map[string]interface{}{"success": false, "error": "user not found"}
			return c.JSON(http.StatusNotFound, reply)
		}
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}

	reply := map[string]interface{}{"success": true, "error": ""}
	return c.JSON(http.StatusCreated, reply)
}
