package util

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

func JsonShaObject(obj interface{}) (string, error) {

	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	hash := sha1.New()
	hash.Write(bytes)

	bs := hash.Sum(nil)
	shastr := fmt.Sprintf("%x", bs)

	return shastr, nil
}
