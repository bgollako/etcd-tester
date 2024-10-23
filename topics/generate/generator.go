package topics

import (
	"etcd-tester/utils"
	"os"
	"strings"

	"github.com/google/uuid"
)

const (
	Delimiter  = ","
	Prefix     = "/device-"
	TopicsFile = "topics.txt"
	Topics     = "topics/" + TopicsFile
)

// generateTopics generates 'num' number of topics
// with given prefix, with a unique uuid appended to it
// All prefixes will be written to the given filename
func GenerateTopics(num int, prefix, filename string) {
	data := strings.Builder{}
	for i := 0; i < num; i++ {
		utils.FatalAnyErr(data.WriteString(prefix))
		utils.FatalAnyErr(data.WriteString(uuid.NewString()))
		utils.FatalAnyErr(data.WriteString(Delimiter))
	}
	utils.Fatal(os.WriteFile(filename, []byte(data.String()), 0777))
}
