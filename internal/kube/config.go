package kube

import "k8s.io/cli-runtime/pkg/genericclioptions"

func GetNamespace(configFlags *genericclioptions.ConfigFlags, allNamespaces bool) string {
	if allNamespaces {
		return ""
	}

	if configFlags.Namespace != nil && *configFlags.Namespace != "" {
		return *configFlags.Namespace
	}

	rawConfig, err := configFlags.ToRawKubeConfigLoader().RawConfig()
	if err == nil {
		if ctx, ok := rawConfig.Contexts[rawConfig.CurrentContext]; ok && ctx.Namespace != "" {
			return ctx.Namespace
		}
	}

	return "default"
}
