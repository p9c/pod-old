package app

import (
	"strings"

	"github.com/tucnak/climax"
)

// GetFlags reads out the flags in a climax.Command and reads the default value stored there into a searchable map
func GetFlags(cmd climax.Command) (out map[string]string) {
	out = make(map[string]string)
	for i := range cmd.Flags {
		usage := strings.Split(cmd.Flags[i].Usage, " ")
		if cmd.Flags[i].Usage == `""` ||
			len(cmd.Flags[i].Usage) < 2 ||
			len(usage) < 2 {
			out[cmd.Flags[i].Name] = ""
		}
		if len(usage) > 1 {
			u := usage[1][1 : len(usage)-2]
			out[cmd.Flags[i].Name] = u
		}
	}
	return
}
