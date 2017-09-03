package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

type Config struct {
	Taints []Taints `yaml:"taints"`
}

func (cfg *Config) Load(file string) {
	// Load the config file from disk
	cfgB, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Sprint("unable to read config: ", err))
	}
	err = yaml.Unmarshal(cfgB, &cfg)
	if err != nil {
		panic(fmt.Sprint("unable to parse config: ", err))
	}
	cfg.Verify()
}

func (cfg *Config) Verify() {
	for _, t := range cfg.Taints {
		if len(t.Labels) == 0 {
			panic("labels may not be empty!")
		}
		if len(t.Taints) == 0 {
			panic("list of taints to apply may not be empty")
		}
	}
}

// Add taints to given node object. Return the object, and a bool value
// indicating if modifications were made.
func (cfg *Config) AddTaints(node *v1.Node) (*v1.Node, bool) {
	var modified = false
CfgLoop:
	for _, t := range cfg.Taints {
		for k, v := range t.Labels {
			nodeLabelValue, ok := node.Labels[k]
			if !ok || nodeLabelValue != v {
				// no match
				continue CfgLoop
			}
		}
		//fmt.Printf("node %v matched all labels: %v\n", node.Name, t.Labels)
		currentTaints := len(node.Spec.Taints)
		newTaints := []api.Taint{}
	TaintLoop:
		for _, nt := range t.Taints {
			typeCasted := v1.Taint{
				Key:    nt.Key,
				Value:  nt.Value,
				Effect: nt.Effect,
			}
			for _, existingTaint := range node.Spec.Taints {
				// reset time field for comparison
				existingTaint.TimeAdded = metav1.Time{}
				if existingTaint == typeCasted {
					continue TaintLoop
				}
			}
			node.Spec.Taints = append(node.Spec.Taints, typeCasted)
			newTaints = append(newTaints, api.Taint{
				Key:    nt.Key,
				Value:  nt.Value,
				Effect: api.TaintEffect(nt.Effect),
			})
		}
		if currentTaints == len(node.Spec.Taints) {
			// If there are no new taints to apply, we can check the next config
			continue CfgLoop
		}
		modified = true
	}
	return node, modified
}

type Taints struct {
	Labels map[string]string `yaml:"labels"`
	Taints []Taint           `yaml:"taints"`
}

// api.Taint sans the Time field
type Taint struct {
	Key    string         `yaml:"key"`
	Value  string         `yaml:"value"`
	Effect v1.TaintEffect `yaml:"effect"`
}
