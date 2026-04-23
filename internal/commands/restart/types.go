package restart

type podRestart struct {
	namespace string
	name      string
	restarts  int32
}
