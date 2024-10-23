package generate

import (
	topics "etcd-tester/topics/generate"
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

	var sb strings.Builder
	for i := 0; i < numKeys; i++ {
		sb.WriteString(uuid.New().String())
		sb.WriteString(topics.Delimiter)
	}

	_, err = file.WriteString(sb.String())
	if err != nil {
		return err
	}

	return nil
}
