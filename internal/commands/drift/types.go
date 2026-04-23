package drift

type driftType string

const (
	driftSynced      driftType = "synced"
	driftRolling     driftType = "rolling"
	driftLatestTag   driftType = "latest-tag"
	driftLatestDrift driftType = "latest-drift"
	driftStuck       driftType = "stuck"
	driftManual      driftType = "manual"
	driftMixed       driftType = "mixed"
	driftNoPods      driftType = "no-pods"
)

type driftInfo struct {
	namespace     string
	deployment    string
	container     string
	desiredImage  string
	actualImages  map[string]int
	actualDigests map[string]int
	replicas      int32
	readyReplicas int32
	podStates     map[string]int
	driftType     driftType
	driftDetail   string
	selector      string
}
