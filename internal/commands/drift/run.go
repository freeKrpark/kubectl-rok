package drift

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/freekrpark/kubectl-rok/internal/kube"
	"github.com/freekrpark/kubectl-rok/internal/output"
	"github.com/freekrpark/kubectl-rok/pkg/image"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func run(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) error {
	clientset, err := kube.NewClientset(configFlags)
	if err != nil {
		return err
	}

	namespace := kube.GetNamespace(configFlags, allNamespaces)

	deploy, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("deployment list failed: %w", err)
	}

	var results []driftInfo
	for _, d := range deploy.Items {
		if len(d.Spec.Template.Spec.Containers) == 0 {
			continue
		}

		mainContainer := d.Spec.Template.Spec.Containers[0]

		var replicas int32 = 1
		if d.Spec.Replicas != nil {
			replicas = *d.Spec.Replicas
		}

		results = append(results, driftInfo{
			namespace:     d.Namespace,
			deployment:    d.Name,
			container:     mainContainer.Name,
			desiredImage:  mainContainer.Image,
			replicas:      replicas,
			readyReplicas: d.Status.ReadyReplicas,
			actualImages:  map[string]int{},
			actualDigests: map[string]int{},
			podStates:     map[string]int{},
			selector:      labels.Set(d.Spec.Selector.MatchLabels).AsSelector().String(),
		})
	}

	for i := range results {
		pods, err := clientset.CoreV1().Pods(results[i].namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: results[i].selector})
		if err != nil {
			return fmt.Errorf("pod list for %s failed: %w", results[i].deployment, err)
		}

		for _, pod := range pods.Items {
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.Name != results[i].container {
					continue
				}

				results[i].actualImages[cs.Image]++
				results[i].podStates[image.StateString(cs.State)]++

				if digest := image.ParseDigest(cs.ImageID); digest != "" {
					results[i].actualDigests[digest]++
				}

			}
		}
		classifyDrift(&results[i])
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].namespace != results[j].namespace {
			return results[i].namespace < results[j].namespace
		}

		if results[i].deployment != results[j].deployment {
			return results[i].deployment < results[j].deployment
		}
		return results[i].container < results[j].container
	})

	w := output.NewTable()
	defer w.Flush()

	if allNamespaces {
		fmt.Fprintln(w, "NAMESPACE\tDEPLOYMENT\tDESIRED\tREPLICAS\tSTATUS\tDETAIL")
		for _, r := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d/%d\t%s %s\t%s\n",
				r.namespace, r.deployment, image.ShortImage(r.desiredImage),
				r.readyReplicas, r.replicas,
				driftIcon(r.driftType), r.driftType,
				r.driftDetail)
		}
	} else {
		fmt.Fprintln(w, "DEPLOYMENT\tDESIRED\tREPLICAS\tSTATUS\tDETAIL")
		for _, r := range results {
			fmt.Fprintf(w, "%s\t%s\t%d/%d\t%s %s\t%s\n",
				r.deployment, image.ShortImage(r.desiredImage),
				r.readyReplicas, r.replicas,
				driftIcon(r.driftType), r.driftType,
				r.driftDetail)
		}
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "No pods found")
	}

	return nil
}
