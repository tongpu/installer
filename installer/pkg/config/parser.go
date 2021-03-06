package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/openshift/installer/pkg/rhcos"
)

// ParseConfig parses a yaml string and returns, if successful, a Cluster.
func ParseConfig(data []byte) (*Cluster, error) {
	cluster := defaultCluster

	if err := yaml.Unmarshal(data, &cluster); err != nil {
		return nil, err
	}

	// Deprecated: remove after openshift/release is ported to pullSecret
	if cluster.PullSecretPath != "" {
		if cluster.PullSecret != "" {
			return nil, errors.New("pullSecretPath is deprecated; just set pullSecret")
		}

		data, err := ioutil.ReadFile(cluster.PullSecretPath)
		if err != nil {
			return nil, err
		}
		cluster.PullSecret = string(data)
	}

	if cluster.EC2AMIOverride == "" {
		ami, err := rhcos.AMI(DefaultChannel, cluster.AWS.Region)
		if err != nil {
			return nil, fmt.Errorf("Failed to determine default AMI: %v", err)
		}
		cluster.EC2AMIOverride = ami
	}

	return &cluster, nil
}

// ParseConfigFile parses a yaml file and returns, if successful, a Cluster.
func ParseConfigFile(path string) (*Cluster, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(data)
}

// ParseInternal parses a yaml string and returns, if successful, an internal.
func ParseInternal(data []byte) (*Internal, error) {
	internal := &Internal{}

	if err := yaml.Unmarshal(data, internal); err != nil {
		return nil, err
	}

	return internal, nil
}

// ParseInternalFile parses a yaml file and returns, if successful, an internal.
func ParseInternalFile(path string) (*Internal, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseInternal(data)
}
