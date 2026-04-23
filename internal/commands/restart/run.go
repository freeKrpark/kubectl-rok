package restart

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/freekrpark/kubectl-rok/internal/kube"
	"github.com/freekrpark/kubectl-rok/internal/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func run(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) error {
	clientset, err := kube.NewClientset(configFlags)
	if err != nil {
		return err
	}

	namespace := kube.GetNamespace(configFlags, allNamespaces)

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("pod list failed: %w", err)
	}

	var results []podRestart
	for _, p := range pods.Items {
		var total int32
		for _, cs := range p.Status.ContainerStatuses {
			total += cs.RestartCount
		}
		results = append(results, podRestart{p.Namespace, p.Name, total})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].restarts > results[j].restarts
	})

	w := output.NewTable()
	defer w.Flush()

	if allNamespaces {
		fmt.Fprintln(w, "NAMESPACE\tPOD\tRESTARTS")
		for _, r := range results {
			fmt.Fprintf(w, "%s\t%s\t%d\n", r.namespace, r.name, r.restarts)
		}
	} else {
		fmt.Fprintln(w, "POD\tRESTARTS")
		for _, r := range results {
			fmt.Fprintf(w, "%s\t%d\n", r.name, r.restarts)
		}
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "No pods found.")
	}

	return nil
}
