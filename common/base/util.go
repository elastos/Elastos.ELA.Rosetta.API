package base

import (
	"strconv"

	"github.com/elastos/Elastos.ELA/common"
)

func GetSelaString(value common.Fixed64) string {
	return strconv.Itoa(int(value))
}

func GetCoinIdentifier(hash common.Uint256, index uint16) string {
	return hash.String() + ":" + strconv.Itoa(int(index))
}
