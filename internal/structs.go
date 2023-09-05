package internal

type Deployment struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Labels    struct {
			App string `yaml:"app"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Strategy struct {
			Type          string `yaml:"type"`
			RollingUpdate struct {
				MaxSurge       int `yaml:"maxSurge"`
				MaxUnavailable int `yaml:"maxUnavailable"`
			} `yaml:"rollingUpdate"`
		} `yaml:"strategy"`
		RevisionHistoryLimit int `yaml:"revisionHistoryLimit"`
		Selector             struct {
			MatchLabels struct {
				App string `yaml:"app"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Template struct {
			Metadata struct {
				Name   string `yaml:"name"`
				Labels struct {
					App string `yaml:"app"`
				} `yaml:"labels"`
				Annotations struct {
					PrometheusIoPath   string `yaml:"prometheus.io/path"`
					PrometheusIoPort   string `yaml:"prometheus.io/port"`
					PrometheusIoScrape string `yaml:"prometheus.io/scrape"`
				} `yaml:"annotations"`
			} `yaml:"metadata"`
			Spec struct {
				Containers []struct {
					Name            string   `yaml:"name"`
					Image           string   `yaml:"image"`
					ImagePullPolicy string   `yaml:"imagePullPolicy"`
					Command         []string `yaml:"command"`
					Args            []string `yaml:"args"`
					Ports           []struct {
						ContainerPort int    `yaml:"containerPort"`
						Protocol      string `yaml:"protocol"`
					} `yaml:"ports"`
					Resources struct {
						Limits struct {
							CPU    string `yaml:"cpu"`
							Memory string `yaml:"memory"`
						} `yaml:"limits"`
						Requests struct {
							CPU    string `yaml:"cpu"`
							Memory string `yaml:"memory"`
						} `yaml:"requests"`
					} `yaml:"resources"`
					VolumeMounts []struct {
						MountPath string `yaml:"mountPath"`
						Name      string `yaml:"name"`
					} `yaml:"volumeMounts"`
					Env []struct {
						Name      string `yaml:"name"`
						ValueFrom struct {
							SecretKeyRef struct {
								Key  string `yaml:"key"`
								Name string `yaml:"name"`
							} `yaml:"secretKeyRef"`
						} `yaml:"valueFrom"`
					} `yaml:"env"`
					LivenessProbe struct {
						HTTPGet struct {
							Path   string `yaml:"path"`
							Port   int    `yaml:"port"`
							Scheme string `yaml:"scheme"`
						} `yaml:"httpGet"`
						InitialDelaySeconds int `yaml:"initialDelaySeconds"`
						PeriodSeconds       int `yaml:"periodSeconds"`
					} `yaml:"livenessProbe"`
					ReadinessProbe struct {
						HTTPGet struct {
							Path   string `yaml:"path"`
							Port   int    `yaml:"port"`
							Scheme string `yaml:"scheme"`
						} `yaml:"httpGet"`
						InitialDelaySeconds int `yaml:"initialDelaySeconds"`
						PeriodSeconds       int `yaml:"periodSeconds"`
					} `yaml:"readinessProbe"`
				} `yaml:"containers"`
				TerminationGracePeriodSeconds int `yaml:"terminationGracePeriodSeconds"`
				ImagePullSecrets              []struct {
					Name string `yaml:"name"`
				} `yaml:"imagePullSecrets"`
				Volumes []struct {
					Name   string `yaml:"name"`
					Secret struct {
						SecretName string `yaml:"secretName"`
					} `yaml:"secret,omitempty"`
					ConfigMap struct {
						Name string `yaml:"name"`
					} `yaml:"configMap,omitempty"`
				} `yaml:"volumes"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

type Config struct {
	Firestates map[string]struct {
		Resources struct {
			Limits struct {
				CPU    string `yaml:"cpu"`
				Memory string `yaml:"memory"`
			} `yaml:"limits"`
			Requests struct {
				CPU    string `yaml:"cpu"`
				Memory string `yaml:"memory"`
			} `yaml:"requests"`
		} `yaml:"resources"`
	} `yaml:"deployment"`
}
