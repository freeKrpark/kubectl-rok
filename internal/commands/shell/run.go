package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/freekrpark/kubectl-rok/internal/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func run(configFlags *genericclioptions.ConfigFlags, podName, container string, debug bool, tool string) error {
	clientset, err := kube.NewClientset(configFlags)
	if err != nil {
		return err
	}

	namespace := kube.GetNamespace(configFlags, false)
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("pod not found: %w", err)
	}

	if container == "" {
		container = pickContainer(pod)
	}

	if debug {
		image := resolveToolImage(tool)
		return execDebug(namespace, podName, container, image)
	}

	for _, shell := range []string{"bash", "sh"} {
		fmt.Fprintf(os.Stderr, "Trying %s in %s/%s...\n", shell, podName, container)
		if err := execShell(namespace, podName, container, shell); err == nil {
			return nil
		}
	}

	return fmt.Errorf(
		"no shell (bash, sh) available in pod %s.\n"+
			"Hint: try --debug to launch an ephemeral container",
		podName,
	)
}

func pickContainer(pod *corev1.Pod) string {
	sidecars := map[string]bool{
		"istio-proxy":   true,
		"linkerd-proxy": true,
		"envoy":         true,
		"logship":       true,
	}

	for _, c := range pod.Spec.Containers {
		if !sidecars[c.Name] {
			return c.Name
		}
	}
	return pod.Spec.Containers[0].Name
}

func resolveToolImage(tool string) string {
	presets := map[string]string{
		"busybox":  "busybox:latest",
		"netshoot": "nicolaka/netshoot:latest",
		"ubuntu":   "ubuntu:latest",
		"alpine":   "alpine:latest",
	}

	if image, ok := presets[tool]; ok {
		return image
	}

	return tool
}

func execShell(namespace, podName, container, shell string) error {
	args := []string{"exec", "-it"}

	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if container != "" {
		args = append(args, "-c", container)
	}

	args = append(args, podName, "--", shell)

	cmd := exec.Command("kubectl", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func execDebug(namespace, podName, target, image string) error {
	args := []string{"debug", "-it"}

	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	args = append(args,
		"--image", image,
		"--target", target,
		"--profile", "general",
		podName)

	fmt.Fprintf(os.Stderr, "Launching ephemeral container (%s) targeting %s/%s...\n",
		image, podName, target)

	cmd := exec.Command("kubectl", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
