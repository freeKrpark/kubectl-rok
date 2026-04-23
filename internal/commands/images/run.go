package images

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/freekrpark/kubectl-rok/internal/kube"
	"github.com/freekrpark/kubectl-rok/internal/output"
	"github.com/freekrpark/kubectl-rok/pkg/image"
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

	var results []podImage
	for _, p := range pods.Items {
		for _, cs := range p.Status.ContainerStatuses {
			repo, tag := image.ParseImage(cs.Image)
			results = append(results, podImage{p.Namespace, p.Name, cs.Name, repo, tag, image.StateString(cs.State)})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].namespace != results[j].namespace {
			return results[i].namespace < results[j].namespace
		}
		if results[i].name != results[j].name {
			return results[i].name < results[j].name
		}
		return results[i].container < results[j].container
	})

	w := output.NewTable()
	defer w.Flush()

	if allNamespaces {
		fmt.Fprintln(w, "NAMESPACE\tPOD\tCONTAINER\tREPOSITORY\tVERSION\tSTATUS")

		for _, r := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				r.namespace, r.name, r.container, r.repository, r.version, r.status)
		}
	} else {
		fmt.Fprintln(w, "POD\tCONTAINER\tREPOSITORY\tVERSION\tSTATUS")
		for _, r := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				r.name, r.container, r.repository, r.version, r.status)
		}
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "No pods found.")
	}
	return nil
}
