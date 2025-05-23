/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha3

import (
	"github.com/blang/semver/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/version"
)

// +kubebuilder:object:root=true

// Metadata for a provider repository.
type Metadata struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// releaseSeries maps a provider release series (major/minor) with a Cluster API contract version.
	// +optional
	ReleaseSeries []ReleaseSeries `json:"releaseSeries"`
}

// ReleaseSeries maps a provider release series (major/minor) with a Cluster API contract version.
type ReleaseSeries struct {
	// major version of the release series
	// +optional
	Major int32 `json:"major,omitempty"`

	// minor version of the release series
	// +optional
	Minor int32 `json:"minor,omitempty"`

	// contract defines the Cluster API contract supported by this series.
	//
	// The value is an API Version, e.g. `v1alpha3`.
	// +optional
	Contract string `json:"contract,omitempty"`
}

func (rs ReleaseSeries) newer(release ReleaseSeries) bool {
	v := semver.Version{Major: uint64(rs.Major), Minor: uint64(rs.Minor)}
	ver := semver.Version{Major: uint64(release.Major), Minor: uint64(release.Minor)}
	return v.GTE(ver)
}

func init() {
	objectTypes = append(objectTypes, &Metadata{})
}

// GetReleaseSeriesForVersion returns the release series for a given version.
func (m *Metadata) GetReleaseSeriesForVersion(version *version.Version) *ReleaseSeries {
	for _, releaseSeries := range m.ReleaseSeries {
		if version.Major() == uint(releaseSeries.Major) && version.Minor() == uint(releaseSeries.Minor) {
			return &releaseSeries
		}
	}

	return nil
}

// GetReleaseSeriesForContract returns the release series for a given API Version, e.g. `v1alpha4`.
// If more than one release series use the same contract then the latest newer release series is
// returned.
func (m *Metadata) GetReleaseSeriesForContract(contract string) *ReleaseSeries {
	var rs ReleaseSeries
	var found bool
	for _, releaseSeries := range m.ReleaseSeries {
		if contract == releaseSeries.Contract {
			found = true
			if releaseSeries.newer(rs) {
				rs = releaseSeries
			}
		}
	}
	if !found {
		return nil
	}
	return &rs
}
