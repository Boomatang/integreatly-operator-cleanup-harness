package tests

import (
	"github.com/integr8ly/integreatly-operator-test-harness/pkg/metadata"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	//buildv1 "github.com/openshift/client-go/config/clientset/versioned/clientset"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"time"
)

var _ = ginkgo.Describe("Integreatly Operator Cleanup", func() {
	defer ginkgo.GinkgoRecover()
	args := os.Args[1:]
	if !Contains(args, "cleanup") {
		logrus.Info("not doing clean up phase")
		os.Exit(0)
	}

	ginkgo.It("cleanup AWS resources using cluster-service", func() {

		config, err := rest.InClusterConfig()

		if err != nil {
			panic(err)
		}

		// Creates the clientset
		//buildClient, err := buildv1.NewForConfig(config)
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}
		Expect(err).NotTo(HaveOccurred())

		// get aws values
		secret, err := clientset.CoreV1().Secrets("kube-system").Get("aws-creds", metav1.GetOptions{})
		AWS_ACCESS_KEY_ID := string(secret.Data["aws_access_key_id"])
		AWS_SECRET_ACCESS_KEY := string(secret.Data["aws_secret_access_key"])
		logrus.Infof("AWS ACCESS %v", AWS_ACCESS_KEY_ID)
		logrus.Infof("AWS SECRET %v", AWS_SECRET_ACCESS_KEY)

		// get cluster id
		// TODO get id, it can be found in ClusterVersion version
		cluster_id := "a4ebf449-c3a7-4fb4-bdd1-7066efd0815d"
		//id, err := buildClient.RESTClient().Get().Resource("ClusterVersion").Name("version").Do().Get()
		id, err := clientset.CoordinationV1().RESTClient().Get().Resource("ClusterVersion").Name("version").Do().Get()

		logrus.Info("info found == ", id)
		logrus.Info("error found == ", err)

		// configure cluster-service pod args
		container_args := []string{"cleanup", cluster_id, "--watch"}

		if Contains(args, "dry-run") {
			logrus.Info("running cluster-service as dry-run")
			container_args = append(container_args, "--dry-run=true")
		}

		// create cluster-service pod
		pod := &v1.Pod{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{Name: "cluster-service"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "cluster-service",
						Image: "quay.io/integreatly/cluster-service:v0.4.0",
						Args:  container_args,
						Env: []v1.EnvVar{
							{
								Name:  "AWS_ACCESS_KEY_ID",
								Value: AWS_ACCESS_KEY_ID,
							}, {
								Name:  "AWS_SECRET_ACCESS_KEY",
								Value: AWS_SECRET_ACCESS_KEY,
							},
						},
					},
				},
				RestartPolicy: "Never",
			},
		}

		_, err = clientset.CoreV1().Pods("kube-system").Create(pod)

		// watch cluster-service pod for completion
		timeout := 35 * time.Minute
		delay := 30 * time.Second

		err = wait.Poll(timeout, delay, func() (done bool, err error) {
			pod, err = clientset.CoreV1().Pods("kube-system").Get("cluster-service", metav1.GetOptions{})
			if err != nil {
				return false, nil
			}

			if pod.Status.Phase == "Succeeded" {
				logrus.Info("pod status is completed")
				return true, nil
			}
			return false, nil
		})

		// add reported value
		if err != nil {
			metadata.Instance.CleanupCompleted = false
		} else {
			metadata.Instance.CleanupCompleted = true
		}

		Expect(err).NotTo(HaveOccurred())
	})
})

func Contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
