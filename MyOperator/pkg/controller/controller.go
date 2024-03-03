package controller

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	types "MyOperator/pkg/apis/foperator.test/v1alpha1"
	"MyOperator/pkg/generated/clientset/versioned/scheme"
	informer "MyOperator/pkg/generated/informers/externalversions/foperator.test/v1alpha1"
	"MyOperator/pkg/generated/listers/foperator.test/v1alpha1"
)

const (
	AddEvent    = "myOperatorAdd"
	UpdateEvent = "myOperatorUpdate"
	DeleteEvent = "myOperatorDel"
)

type Element struct {
	Type string      `json:"type"` // resource type, eg: pod
	Key  string      `json:"key"`  // MetaNamespaceKey, format: <namespace>/<name>
	Obj  interface{} `json:"obj"`  // object
}

type Controller struct {
	myOperator       v1alpha1.MyOperatorLister
	myOperatorSynced cache.InformerSynced
	queue            workqueue.RateLimitingInterface
}

func NewController(operators informer.MyOperatorInformer) *Controller {
	c := &Controller{
		myOperator:       operators.Lister(),
		myOperatorSynced: operators.Informer().HasSynced,
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "my-controller"),
	}

	runtime.Must(scheme.AddToScheme(k8sscheme.Scheme))
	operators.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				runtime.HandleError(err)
				return
			}
			e := Element{
				Type: AddEvent,
				Key:  key,
				Obj:  obj,
			}
			c.queue.Add(e)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldNetwork := old.(*types.MyOperator)
			newNetwork := new.(*types.MyOperator)
			if oldNetwork.ResourceVersion == newNetwork.ResourceVersion {
				return
			}
			key, err := cache.MetaNamespaceKeyFunc(new)
			e := Element{
				Type: UpdateEvent,
				Key:  key,
				Obj:  new,
			}
			if err != nil {
				runtime.HandleError(err)
				return
			}
			c.queue.Add(e)
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			e := Element{
				Type: DeleteEvent,
				Key:  key,
				Obj:  obj,
			}
			if err != nil {
				runtime.HandleError(err)
				return
			}
			c.queue.Add(e)
		},
	})

	return c
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	fmt.Printf("Starting controller\n")

	if !cache.WaitForCacheSync(stopCh, c.myOperatorSynced) {
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	fmt.Printf("Shutting down <NAME> controller\n")
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	defer c.queue.Done(key)

	err := func(obj interface{}) error {
		var element Element
		var ok bool
		if element, ok = obj.(Element); !ok {
			c.queue.Forget(obj)
			fmt.Printf("expected string in workqueue but got %#v \n", obj)
			return nil
		}
		if err := c.syncHandler(element); err != nil {
			c.queue.AddRateLimited(element)
			return fmt.Errorf("error syncing '%+v': %w, requeuing \n", element, err)
		}
		c.queue.Forget(obj)
		return nil
	}(key)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(element Element) error {
	fmt.Printf("syncHandler start process: %+v \n", element)

	switch element.Type {
	case AddEvent:
		fmt.Printf("add myOperator : %+v\n", element)
	case UpdateEvent:
		fmt.Printf("update myOperator : %+v\n", element)
	case DeleteEvent:
		fmt.Printf("del myOperator : %+v\n", element)
	default:
		fmt.Printf("element type: %s not supported \n", element.Type)
		return nil
	}

	return nil
}
