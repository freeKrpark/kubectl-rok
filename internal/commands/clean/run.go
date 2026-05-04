package clean

import (
	"context"
	"fmt"
	"os"

	"github.com/freekrpark/kubectl-rok/internal/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func run(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) error {
	clientset, err := kube.NewClientset(configFlags)

	if err != nil {
		return err
	}
	namespace := kube.GetNamespace(configFlags, allNamespaces)

	rs, err := clientset.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("replicaset list failed: %w", err)
	}

	for _, r := range rs.Items {
		if *r.Spec.Replicas == 0 {
			err := clientset.AppsV1().ReplicaSets(r.Namespace).Delete(context.TODO(), r.Name, metav1.DeleteOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "error deleting rs/%s: %v\n", r.Name, err)
				continue
			}
			fmt.Fprintf(os.Stdout, "deleted rs/%s in %s\n", r.Name, r.Namespace)
		}
	}
	return nil
}
