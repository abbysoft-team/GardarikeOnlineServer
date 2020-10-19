package model

import rpc "projectx-server/rpc/generated"

type Account struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Salt     string `db:"salt"`
}

type Character struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Gold uint64 `db:"gold"`
}

func ToRPCCharacter(char Character) *rpc.Character {
	return &rpc.Character{
		Id:   int32(char.ID),
		Name: char.Name,
		Gold: char.Gold,
	}
}
