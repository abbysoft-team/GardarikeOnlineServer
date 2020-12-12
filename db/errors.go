package db

import "errors"

var ErrDuplicatedUniqueKey = errors.New("unique key was duplicated")
