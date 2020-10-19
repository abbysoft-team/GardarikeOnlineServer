package game

import (
	"crypto/md5"
	"fmt"
	"projectx-server/model"
	rpc "projectx-server/rpc/generated"
)

func (s *SimpleLogic) Login(request *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	s.log.WithField("login", request.GetUsername()).Debug("Login request")

	acc, err := s.db.GetAccount(request.GetUsername())
	if err != nil {
		return nil, fmt.Errorf("invalid username/password combination")
	}

	hashedPass := md5.Sum([]byte(request.Password))
	saltedHash := fmt.Sprintf("%s%x%s", acc.Salt, string(hashedPass[:]), acc.Salt)
	finalPass := md5.Sum([]byte(saltedHash))

	if fmt.Sprintf("%x", string(finalPass[:])) != acc.Password {
		return nil, fmt.Errorf("invalid username/password combination")
	}

	chars, err := s.db.GetCharacters(acc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account characters: %w", err)
	}

	var rpcChars []*rpc.Character
	for _, char := range chars {
		rpcChars = append(rpcChars, model.ToRPCCharacter(char))
	}

	return &rpc.LoginResponse{
		Characters: rpcChars,
	}, nil
}
