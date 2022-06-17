package kecpws

import "errors"

var (
	ErrCanNotCreateTheRoom = errors.New("cannot creat the room")
	ErrCanNotJoinTheRoom   = errors.New("cannot join the room")
)
