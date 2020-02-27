package migrate

import (
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	ProwClusterURL = "https://api.ci.openshift.org"
)

var (
	migratedRepos = sets.NewString(
		"openshift/jenkins-openshift-login-plugin/*",
		//find ./ci-operator/config -type d -depth 2 | head -n 230 | while read i; do echo "\"${i#./ci-operator/config/}/.*\","; done
		"ostreedev/ostree/.*",
		"openshift-priv/csi-external-attacher/.*",
		"openshift-priv/cluster-api-provider-azure/.*",
		"openshift-priv/cluster-update-keys/.*",
		"openshift-priv/vertical-pod-autoscaler-operator/.*",
		"openshift-priv/multus-cni/.*",
		"openshift-priv/oauth-server/.*",
		"openshift-priv/template-service-broker-operator/.*",
		"openshift-priv/ci-experiment-origin/.*",
		"openshift-priv/kubernetes-kube-storage-version-migrator/.*",
		"openshift-priv/openshift-state-metrics/.*",
		"openshift-priv/cluster-api-provider-baremetal/.*",
		"openshift-priv/kube-state-metrics/.*",
		"openshift-priv/dedicated-admin-operator/.*",
		"openshift-priv/loki/.*",
		"openshift-priv/cluster-capacity/.*",
		"openshift-priv/cluster-version-operator/.*",
		"openshift-priv/operator-lifecycle-manager/.*",
		"openshift-priv/presto/.*",
		"openshift-priv/cluster-dns-operator/.*",
		"openshift-priv/crd-schema-gen/.*",
		"openshift-priv/operator-registry/.*",
		"openshift-priv/oauth-proxy/.*",
		"openshift-priv/cluster-nfd-operator/.*",
		"openshift-priv/pagerduty-operator/.*",
		"openshift-priv/descheduler/.*",
		"openshift-priv/client-go/.*",
		"openshift-priv/leader-elector/.*",
		"openshift-priv/openshift-tuned/.*",
		"openshift-priv/cluster-autoscaler-operator/.*",
		"openshift-priv/service-ca-operator/.*",
		"openshift-priv/jenkins-client-plugin/.*",
		"openshift-priv/ocs-operator/.*",
		"openshift-priv/rbac-permissions-operator/.*",
		"openshift-priv/cluster-config-operator/.*",
		"openshift-priv/kubecsr/.*",
		"openshift-priv/kuryr-kubernetes/.*",
		"openshift-priv/cluster-bootstrap/.*",
		"openshift-priv/whereabouts-cni/.*",
		"openshift-priv/service-catalog/.*",
		"openshift-priv/cluster-api/.*",
		"openshift-priv/cluster-api-provider-gcp/.*",
		"openshift-priv/compliance-operator/.*",
		"openshift-priv/sdn/.*",
		"openshift-priv/hyperconverged-cluster-operator/.*",
		"openshift-priv/kubefed/.*",
		"openshift-priv/machine-api-operator/.*",
		"openshift-priv/images/.*",
		"openshift-priv/node-feature-discovery/.*",
		"openshift-priv/route-override-cni/.*",
		"openshift-priv/cluster-svcat-controller-manager-operator/.*",
		"openshift-priv/elasticsearch-operator/.*",
		"openshift-priv/ghostunnel/.*",
		"openshift-priv/containernetworking-plugins/.*",
		"openshift-priv/configure-alertmanager-operator/.*",
		"openshift-priv/cluster-api-provider-libvirt/.*",
		"openshift-priv/apiserver-library-go/.*",
		"openshift-priv/must-gather/.*",
		"openshift-priv/operator-marketplace/.*",
		"openshift-priv/csi-node-driver-registrar/.*",
		"openshift-priv/csi-driver-nfs/.*",
		"openshift-priv/prometheus-operator/.*",
		"openshift-priv/oc/.*",
		"openshift-priv/csi-external-resizer/.*",
		"openshift-priv/template-service-broker/.*",
		"openshift-priv/ovn-kubernetes/.*",
		"openshift-priv/cluster-api-actuator-pkg/.*",
		"openshift-priv/cluster-kube-controller-manager-operator/.*",
		"openshift-priv/multus-admission-controller/.*",
		"openshift-priv/library-go/.*",
		"openshift-priv/baremetal-operator/.*",
		"openshift-priv/cluster-api-provider-openstack/.*",
		"openshift-priv/cluster-openshift-controller-manager-operator/.*",
		"openshift-priv/ironic-static-ip-manager/.*",
		"openshift-priv/cluster-kube-descheduler-operator/.*",
		"openshift-priv/node_exporter/.*",
		"openshift-priv/sriov-network-device-plugin/.*",
		"openshift-priv/cluster-image-registry-operator/.*",
		"openshift-priv/cluster-ingress-operator/.*",
		"openshift-priv/sriov-dp-admission-controller/.*",
		"openshift-priv/elasticsearch-proxy/.*",
		"openshift-priv/openshift-apiserver/.*",
		"openshift-priv/grafana/.*",
		"openshift-priv/operator-metering/.*",
		"openshift-priv/kubefed-operator/.*",
		"openshift-priv/cluster-machine-approver/.*",
		"openshift-priv/kube-rbac-proxy/.*",
		"openshift-priv/aws-account-operator/.*",
		"openshift-priv/ironic-rhcos-downloader/.*",
		"openshift-priv/cluster-etcd-operator/.*",
		"openshift-priv/build-machinery-go/.*",
		"openshift-priv/runtime-utils/.*",
		"openshift-priv/cluster-policy-controller/.*",
		"openshift-priv/cluster-kube-scheduler-operator/.*",
		"openshift-priv/cluster-node-tuning-operator/.*",
		"openshift-priv/origin-aggregated-logging/.*",
		"openshift-priv/cluster-csi-snapshot-controller-operator/.*",
		"openshift-priv/federation-v2-operator/.*",
		"openshift-priv/ovirt-csi-driver/.*",
		"openshift-priv/certman-operator/.*",
		"openshift-priv/oauth-apiserver/.*",
		"openshift-priv/deadmanssnitch-operator/.*",
		"openshift-priv/cluster-api-provider-aws/.*",
		"openshift-priv/metal3-smart-exporter/.*",
		"openshift-priv/csi-external-snapshotter/.*",
		"openshift-priv/cluster-resource-override-admission-operator/.*",
		"openshift-priv/windows-machine-config-bootstrapper/.*",
		"openshift-priv/sriov-network-operator/.*",
		"openshift-priv/machine-config-operator/.*",
		"openshift-priv/console-operator/.*",
		"openshift-priv/csi-operator/.*",
		"openshift-priv/k8s-prometheus-adapter/.*",
		"openshift-priv/cluster-resource-override-admission/.*",
		"openshift-priv/ironic-hardware-inventory-recorder-image/.*",
		"openshift-priv/configmap-reload/.*",
		"openshift-priv/ansible-service-broker/.*",
		"openshift-priv/insights-operator/.*",
		"openshift-priv/hadoop/.*",
		"openshift-priv/csi-cluster-driver-registrar/.*",
		"openshift-priv/sig-storage-local-static-provisioner/.*",
		"openshift-priv/jenkins-sync-plugin/.*",
		"openshift-priv/image-registry/.*",
		"openshift-priv/cloud-credential-operator/.*",
		"openshift-priv/etcd/.*",
		"openshift-priv/thanos/.*",
		"openshift-priv/linuxptp-daemon/.*",
		"openshift-priv/ptp-operator/.*",
		"openshift-priv/cluster-kube-storage-version-migrator-operator/.*",
		"openshift-priv/node-problem-detector/.*",
		"openshift-priv/openshift-ansible/.*",
		"openshift-priv/hive/.*",
		"openshift-priv/api/.*",
		"openshift-priv/cluster-svcat-apiserver-operator/.*",
		"openshift-priv/node-problem-detector-operator/.*",
		"openshift-priv/ironic-inspector-image/.*",
		"openshift-priv/ironic-image/.*",
		"openshift-priv/cloud-provider-openstack/.*",
		"openshift-priv/jenkins-openshift-login-plugin/.*",
		"openshift-priv/ironic-ipa-downloader/.*",
		"openshift-priv/prometheus/.*",
		"openshift-priv/origin/.*",
		"openshift-priv/cluster-openshift-apiserver-operator/.*",
		"openshift-priv/cluster-samples-operator/.*",
		"openshift-priv/jenkins/.*",
		"openshift-priv/prom-label-proxy/.*",
		"openshift-priv/coredns/.*",
		"openshift-priv/csi-driver-registrar/.*",
		"openshift-priv/ocp-release-operator-sdk/.*",
		"openshift-priv/installer/.*",
		"openshift-priv/cluster-api-provider-ovirt/.*",
		"openshift-priv/cluster-kube-apiserver-operator/.*",
		"openshift-priv/local-storage-operator/.*",
		"openshift-priv/builder/.*",
		"openshift-priv/baremetal-runtimecfg/.*",
		"openshift-priv/csi-livenessprobe/.*",
		"openshift-priv/openshift-controller-manager/.*",
		"openshift-priv/cluster-network-operator/.*",
		"openshift-priv/kubernetes-autoscaler/.*",
		"openshift-priv/cluster-logging-operator/.*",
		"openshift-priv/cluster-storage-operator/.*",
		"openshift-priv/mdns-publisher/.*",
		"openshift-priv/external-storage/.*",
		"openshift-priv/csi-external-provisioner/.*",
		"openshift-priv/helm/.*",
		"openshift-priv/telemeter/.*",
		"openshift-priv/cluster-authentication-operator/.*",
		"openshift-priv/managed-cluster-config/.*",
		"openshift-priv/prometheus-alertmanager/.*",
		"openshift-priv/openshift-tests/.*",
		"openshift-priv/file-integrity-operator/.*",
		"openshift-priv/sriov-cni/.*",
		"openshift-priv/console/.*",
		"openshift-priv/cluster-monitoring-operator/.*",
		"openshift-priv/cluster-api-provider-kubemark/.*",
		"openshift-priv/router/.*",
		"heketi/heketi/.*",
		"cri-o/cri-o/.*",
		"monstorak/monstorak-operator/.*",
		"openshift-cnv/cnv-ci/.*",
		"openshift-knative/serving-operator/.*",
		"openshift-knative/serverless-operator/.*",
		"openshift-knative/kourier/.*",
		"openshift-kni/cnf-features-deploy/.*",
		"openshift-kni/baremetal-deploy/.*",
		"openshift-kni/performance-addon-operators/.*",
		"integr8ly/cloud-resource-operator/.*",
		"integr8ly/integreatly-operator/.*",
		"integr8ly/installation/.*",
		"integr8ly/heimdall/.*",
		//"integr8ly/ansible-tower-configuration/.*",
		"fabric8-services/toolchain-operator/.*",
		"openshift/csi-external-attacher/.*",
		"openshift/cluster-api-provider-azure/.*",
		"openshift/knative-serving/.*",
		"openshift/cluster-update-keys/.*",
		"openshift/vertical-pod-autoscaler-operator/.*",
		"openshift/multus-cni/.*",
		"openshift/oauth-server/.*",
		"openshift/template-service-broker-operator/.*",
		"openshift/kubernetes-kube-storage-version-migrator/.*",
		"openshift/openshift-state-metrics/.*",
		"openshift/cluster-api-provider-baremetal/.*",
		"openshift/kube-state-metrics/.*",
		"openshift/dedicated-admin-operator/.*",
		"openshift/loki/.*",
		"openshift/cluster-capacity/.*",
		"openshift/cluster-version-operator/.*",
		"openshift/odo/.*",
		"openshift/cluster-dns-operator/.*",
		"openshift/crd-schema-gen/.*",
		"openshift/oauth-proxy/.*",
		"openshift/cluster-nfd-operator/.*",
		"openshift/pagerduty-operator/.*",
		"openshift/descheduler/.*",
		"openshift/client-go/.*",
		"openshift/leader-elector/.*",
		"openshift/openshift-tuned/.*",
		"openshift/managed-cluster-validating-webhooks/.*",
		"openshift/cluster-autoscaler-operator/.*",
		"openshift/service-ca-operator/.*",
		"openshift/jenkins-client-plugin/.*",
		"openshift/ocs-operator/.*",
		"openshift/rbac-permissions-operator/.*",
		"openshift/cluster-config-operator/.*",
		"openshift/kubecsr/.*",
		"openshift/kuryr-kubernetes/.*",
		"openshift/cluster-bootstrap/.*",
		"openshift/whereabouts-cni/.*",
		"openshift/origin-web-console-server/.*",
		"openshift/service-catalog/.*",
	)
	migratedRegexes []*regexp.Regexp
)

func init() {
	for _, migratedRepo := range migratedRepos.List() {
		migratedRegexes = append(migratedRegexes, regexp.MustCompile(migratedRepo))
	}
}

func Migrated(org, repo, branch string) bool {
	for _, regex := range migratedRegexes {
		if regex.MatchString(fmt.Sprintf("%s/%s/%s", org, repo, branch)) {
			return true
		}
	}
	return false
}
