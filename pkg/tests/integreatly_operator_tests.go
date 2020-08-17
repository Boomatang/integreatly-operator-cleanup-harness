package tests

import (
	"github.com/integr8ly/integreatly-operator-test-harness/pkg/metadata"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	k8sv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
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
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		//apiextensions, err := clientset.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		// get aws values
		secret, err := clientset.CoreV1().Secrets("kube-system").Get("aws-creds", metav1.GetOptions{})
		AWS_ACCESS_KEY_ID := string(secret.Data["aws_access_key_id"])
		AWS_SECRET_ACCESS_KEY := string(secret.Data["aws_secret_access_key"])
		logrus.Infof("AWS ACCESS %v", AWS_ACCESS_KEY_ID)
		logrus.Infof("AWS SECRET %v", AWS_SECRET_ACCESS_KEY)

		// get cluster id
		cluster_id := "a4ebf449-c3a7-4fb4-bdd1-7066efd0815d"

		// configure container args
		container_args := []string{"cleanup", cluster_id}

		if Contains(args, "dry-run") {
			logrus.Info("running cluster-service as dry-run")
			container_args = append(container_args, "--dry-run=true")
		}

		// create cluster-service
		var replicas = int32(1)
		d := &k8sv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster-service",
			},
			Spec: k8sv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "cluster-service",
					},
				},

				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "cluster-service",
						},
					},
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
					},
				},
			},
		}

		_, err = clientset.AppsV1().Deployments("kube-system").Create(d)

		if err != nil {
			metadata.Instance.FoundCRD = false
		} else {
			metadata.Instance.FoundCRD = true
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
