/*
   Copyright 2023 Sangfor Technologies. All rights reserved.
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "MyOperator/pkg/generated/clientset/versioned/typed/foperator.test/v1alpha1"

	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeFoperatorV1alpha1 struct {
	*testing.Fake
}

func (c *FakeFoperatorV1alpha1) MyOperators(namespace string) v1alpha1.MyOperatorInterface {
	return &FakeMyOperators{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeFoperatorV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}