package main

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	Taints []Taint `yaml:"taints"`
}

type Taints struct {
	Selector v1.LabelSelector `yaml:"selector"`
	Taints   []Taint          `yaml:"taints"`
}

type Taint struct {
	Key    string `yaml:"key"`
	Value  string `yaml:"value"`
	Effect string `yaml:"effect"`
}
