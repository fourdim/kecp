package kecpws

import "errors"

var (
	ErrCanNotCreateTheRoom = errors.New("cannot create the room")
	ErrCanNotJoinTheRoom   = errors.New("cannot join the room")
	ErrNameIsAlreadyInUse  = errors.New("name is already in use")
	ErrNotAValidName       = errors.New("not a valid name")
)
