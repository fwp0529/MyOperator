package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	myController "MyOperator/pkg/controller"
	"MyOperator/pkg/generated/clientset/versioned"
	_ "MyOperator/pkg/generated/clientset/versioned/scheme"
	"MyOperator/pkg/generated/informers/externalversions"
)

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("read kube config failed, err: %v\n", err)
		config, err = clientcmd.BuildConfigFromFlags("", "./kubeconfig")
	}
	clientSet, err := versioned.NewForConfig(config)
	if err != nil {
		fmt.Printf("new client set failed, err: %v\n", err)
	}
	clientSet.FoperatorV1alpha1().MyOperators("fns").List(context.TODO(), metav1.ListOptions{})
	myOperatorInformerFactory := externalversions.NewSharedInformerFactory(clientSet, time.Second)
	controller := myController.NewController(myOperatorInformerFactory.Foperator().V1alpha1().MyOperators())
	stopCh := make(chan struct{})

	go myOperatorInformerFactory.Start(stopCh)

	go controller.Run(2, stopCh)
	fmt.Printf("controller run!\n")

	listenSig(stopCh)
}

func listenSig(stopCh chan struct{}) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigCh
	stopCh <- struct{}{}
	close(stopCh)
	return
}
