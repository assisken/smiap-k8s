package differ

import (
	"encoding/hex"
	"encoding/json"

	"golang.org/x/crypto/blake2b"
)

type reason string

const (
	NotCompared  reason = "Not compared yet"
	IsWanted     reason = "Object is wanted but not present"
	NotWanted    reason = "Object is not wanted but present"
	HashMismatch reason = "Hash mismatch"
)

type DiffItem struct {
	Id     string
	Object DiffObject
	Reason reason

	hash string
}

func NewDiffItem(id string, obj DiffObject) (DiffItem, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return DiffItem{}, err
	}

	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	hash, ok := annotations["hash"]
	if !ok {
		hashBytes := blake2b.Sum512(data)
		hash = hex.EncodeToString(hashBytes[:])
	}

	annotations["hash"] = hash
	obj.SetAnnotations(annotations)

	return DiffItem{
		Id:     id,
		Object: obj,
		hash:   hash,
		Reason: NotCompared,
	}, nil
}
