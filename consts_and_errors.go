package main

import "errors"

// TODO : ADD MORE ERRORS (mainly for repositories and where you see 'id must be > 0' or smth like that, just add it here)
// TODO : ADD MORE CONSTS (such as DB_URL, PORT, etc)

const (
	DB_URL_KEY = "DB_URL" // Key for DB_URL env var, value is being set in database.env
)

var (
	ErrDBisNotSet            = errors.New(DB_URL_KEY + " is not set")                                // Error returned when DB_URL is not set, check env vars
	ErrIdMustBeGtZero        = errors.New("id must be greater than 0")                               // Error returned when id is not greater than 0
	ErrLenNameIsZero         = errors.New("the length of name must be greater than 0")               // Error returned when len(name) is 0
	ErrPasswordMustBeGt6     = errors.New("the length of a password must be greater than 6 symbols") //
	ErrOldPasswordIsWrong    = errors.New("old password is incorrect")                               // When the old password is incorrect
	ErrNewPasswordIsSame     = errors.New("new password must be different from old password")        // When the new password is the same as the old password
	ErrUserNotFound          = errors.New("user not found")                                          // When user with this id does not exist
	ErrNoUsers               = errors.New("no users found")                                          // When no users exist
	ErrNoTasks               = errors.New("no tasks found")                                          // When no tasks exist
	ErrEmptyTitle            = errors.New("title must be not empty")                                 // When a title is empty
	ErrNoUserWithThisId      = errors.New("user with this id does not exist")                        // When user with this id does not exist
	ErrTaskDescNotUpdated    = errors.New("task's description was not updated")                      // when description was not updated due to a 'no rows affected' error
	ErrTaskStatusNotSwitched = errors.New("task's status was not switched")                          // when task status was not switched due to a 'no rows affected' error
	ErrTaskTitleNotUpdated   = errors.New("task's title was not updated")                            // when a task title was not updated due to a 'no rows affected' error
)
