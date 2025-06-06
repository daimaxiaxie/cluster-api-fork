//go:build !race

/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package upstreamv1beta2

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"

	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

// Test is disabled when the race detector is enabled (via "//go:build !race" above) because otherwise the fuzz tests would just time out.

func TestFuzzyConversion(t *testing.T) {
	t.Run("for ClusterConfiguration", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &bootstrapv1.ClusterConfiguration{},
		Spoke: &ClusterConfiguration{},
		// NOTE: Kubeadm types does not have ObjectMeta, so we are required to skip data annotation cleanup in the spoke-hub-spoke round trip test.
		SkipSpokeAnnotationCleanup: true,
		FuzzerFuncs:                []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))
	t.Run("for ClusterStatus", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &bootstrapv1.ClusterStatus{},
		Spoke: &ClusterStatus{},
		// NOTE: Kubeadm types does not have ObjectMeta, so we are required to skip data annotation cleanup in the spoke-hub-spoke round trip test.
		SkipSpokeAnnotationCleanup: true,
		FuzzerFuncs:                []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))
	t.Run("for InitConfiguration", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &bootstrapv1.InitConfiguration{},
		Spoke: &InitConfiguration{},
		// NOTE: Kubeadm types does not have ObjectMeta, so we are required to skip data annotation cleanup in the spoke-hub-spoke round trip test.
		SkipSpokeAnnotationCleanup: true,
		FuzzerFuncs:                []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))
	t.Run("for JoinConfiguration", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Hub:   &bootstrapv1.JoinConfiguration{},
		Spoke: &JoinConfiguration{},
		// NOTE: Kubeadm types does not have ObjectMeta, so we are required to skip data annotation cleanup in the spoke-hub-spoke round trip test.
		SkipSpokeAnnotationCleanup: true,
		FuzzerFuncs:                []fuzzer.FuzzerFuncs{fuzzFuncs},
	}))
}

func fuzzFuncs(_ runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		clusterConfigurationFuzzer,
		dnsFuzzer,
		initConfigurationFuzzer,
		joinControlPlanesFuzzer,
		bootstrapv1ControlPlaneComponentFuzzer,
		bootstrapv1LocalEtcdFuzzer,
		bootstrapv1InitConfigurationFuzzer,
		bootstrapv1JoinConfigurationFuzzer,
		bootstrapv1NodeRegistrationOptionsFuzzer,
	}
}

// Custom fuzzers for kubeadm v1beta2 types.
// NOTES:
// - When fields do does not exist in cabpk v1beta1 types, pinning it to avoid kubeadm v1beta2 --> cabpk v1beta1 --> kubeadm v1beta2 round trip errors.

func clusterConfigurationFuzzer(obj *ClusterConfiguration, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.UseHyperKubeImage = false
}

func dnsFuzzer(obj *DNS, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.Type = ""
}

func initConfigurationFuzzer(obj *InitConfiguration, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.CertificateKey = ""
}

func joinControlPlanesFuzzer(obj *JoinControlPlane, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.CertificateKey = ""
}

// Custom fuzzers for CABPK v1beta1 types.
// NOTES:
// - When fields do not exist in kubeadm v1beta2 types, pinning it to avoid cabpk v1beta1 --> kubeadm v1beta2 --> cabpk v1beta1 round trip errors.

func bootstrapv1ControlPlaneComponentFuzzer(obj *bootstrapv1.ControlPlaneComponent, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.ExtraEnvs = nil
}

func bootstrapv1LocalEtcdFuzzer(obj *bootstrapv1.LocalEtcd, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.ExtraEnvs = nil
}

func bootstrapv1InitConfigurationFuzzer(obj *bootstrapv1.InitConfiguration, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.Patches = nil
	obj.SkipPhases = nil
}

func bootstrapv1JoinConfigurationFuzzer(obj *bootstrapv1.JoinConfiguration, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.Patches = nil
	obj.SkipPhases = nil

	if obj.Discovery.File != nil {
		obj.Discovery.File.KubeConfig = nil
	}
}

func bootstrapv1NodeRegistrationOptionsFuzzer(obj *bootstrapv1.NodeRegistrationOptions, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	obj.ImagePullPolicy = ""
	obj.ImagePullSerial = nil
}
