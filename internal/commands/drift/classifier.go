package drift

import (
	"fmt"
	"strings"

	"github.com/freekrpark/kubectl-rok/pkg/image"
)

func classifyDrift(info *driftInfo) {
	if len(info.actualImages) == 0 {
		info.driftType = driftNoPods
		if info.replicas == 0 {
			info.driftDetail = "replicas=0"
		} else {
			info.driftDetail = "no pods running"
		}
		return
	}

	if reason, count := findStuckPods(info.podStates); reason != "" {
		info.driftType = driftStuck
		info.driftDetail = fmt.Sprintf("%d pods %s", count, reason)
		return
	}

	if len(info.actualImages) == 1 {
		classifySingleImage(info)
		return
	}

	classifyMultipleImages(info)

}

func findStuckPods(states map[string]int) (reason string, count int) {
	stuckReasons := map[string]bool{
		"Waiting:ImagePullBackOff":     true,
		"Waiting:ErrImagePull":         true,
		"Waiting:CrashLoopBackOff":     true,
		"Waiting:CreateContainerError": true,
		"Terminated:Error":             true,
		"Terminated:OOMKilled":         true,
	}

	for state, cnt := range states {
		if stuckReasons[state] {
			return state, cnt
		}
	}
	return "", 0
}

func classifySingleImage(info *driftInfo) {
	var runningImage string
	var runningCount int

	for img, cnt := range info.actualImages {
		runningImage = img
		runningCount = cnt
	}

	isLatest := strings.HasSuffix(runningImage, ":latest") ||
		!strings.Contains(image.ImageBasename(runningImage), ":")

	if isLatest {
		if len(info.actualDigests) > 1 {
			info.driftType = driftLatestDrift
			info.driftDetail = fmt.Sprintf("latest tag, %d different digests", len(info.actualDigests))
			return
		}

		info.driftType = driftLatestTag
		info.driftDetail = "using latest tag"
		return
	}

	if runningImage == info.desiredImage {
		info.driftType = driftSynced
		info.driftDetail = fmt.Sprintf("%d pods match", runningCount)
		return
	}

	info.driftType = driftManual
	info.driftDetail = fmt.Sprintf("running %s (spec: %s)", image.ShortImage(runningImage), image.ShortImage(info.desiredImage))
}

func classifyMultipleImages(info *driftInfo) {
	desiredCount := info.actualImages[info.desiredImage]
	total := 0
	for _, cnt := range info.actualImages {
		total += cnt
	}

	if desiredCount > 0 {
		info.driftType = driftRolling
		info.driftDetail = fmt.Sprintf("%d/%d pods updated", desiredCount, total)
		return
	}

	info.driftType = driftMixed
	images := []string{}
	for img, cnt := range info.actualImages {
		images = append(images, fmt.Sprintf("%s(%d)", image.ShortImage(img), cnt))
	}

	info.driftDetail = fmt.Sprintf("mixed: %s (spec: %s)", strings.Join(images, ","), image.ShortImage(info.desiredImage))
}

func driftIcon(dt driftType) string {
	switch dt {
	case driftSynced:
		return "✓"
	case driftRolling, driftLatestTag:
		return "⚠"
	case driftStuck, driftManual,
		driftLatestDrift, driftMixed:
		return "✗"
	default:
		return "?"
	}
}
