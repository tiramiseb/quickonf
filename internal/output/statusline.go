package output

import (
	"fmt"
	"strings"
)

func statusline() string {
	status := fmt.Sprintf(" %d remaining -> %d running -> %d finished / %d failed ", nbWaitingGroups, nbRunningGroups, nbFinishedGroups, nbFailedGroups)
	if len(status) > width {
		status = fmt.Sprintf("%d rem -> %d run -> %d fin / %d fail", nbWaitingGroups, nbRunningGroups, nbFinishedGroups, nbFailedGroups)
		if len(status) > width {
			status = fmt.Sprintf("%d->%d->%d->%d", nbWaitingGroups, nbRunningGroups, nbFinishedGroups, nbFailedGroups)
		}
	}
	spaces := strings.Repeat(" ", width-len(status))
	return bgBlack + fgWhite + italic + status + spaces + reset
}
