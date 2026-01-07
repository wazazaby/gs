package gs

import (
	"github.com/go4org/hashtriemap"
)

type Map[K comparable, V any] = hashtriemap.HashTrieMap[K, V]
