package image

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func StateString(state corev1.ContainerState) string {
	switch {
	case state.Running != nil:
		return "Running"
	case state.Waiting != nil:
		if state.Waiting.Reason != "" {
			return "Waiting:" + state.Waiting.Reason
		}
		return "Waiting"
	case state.Terminated != nil:
		t := state.Terminated
		if t.Reason != "" {
			return fmt.Sprintf("Terminated:%s(exit=%d)", t.Reason, t.ExitCode)
		}
		return fmt.Sprintf("Terminated(exit=%d)", t.ExitCode)
	default:
		return "Unknown"
	}
}
