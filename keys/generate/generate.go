package generate

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

func GenerateKeys(numKeys int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = uuid.New().String()
	}

	_, err = file.WriteString(strings.Join(keys, ","))
	if err != nil {
		return err
	}

	return nil
}
