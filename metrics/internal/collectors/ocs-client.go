package collectors

import (
	ocsv1alpha1 "github.com/red-hat-storage/ocs-operator/api/v1alpha1"
	"github.com/red-hat-storage/ocs-operator/metrics/internal/options"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func GetOcsClient(opts *options.Options) (*rest.RESTClient, error) {
	ocsConfig, err := clientcmd.BuildConfigFromFlags("", opts.KubeconfigPath)
	if err != nil {
		return nil, err
	}
	scheme := runtime.NewScheme()
	utilruntime.Must(ocsv1alpha1.AddToScheme(scheme))
	codecs := serializer.NewCodecFactory(scheme)
	ocsConfig.GroupVersion = &ocsv1alpha1.GroupVersion
	ocsConfig.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: codecs}
	ocsConfig.APIPath = "/apis"
	ocsConfig.ContentType = runtime.ContentTypeJSON
	if ocsConfig.UserAgent == "" {
		ocsConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	ocsClient, err := rest.RESTClientFor(ocsConfig)
	if err != nil {
		return nil, err
	}
	return ocsClient, nil
}

type StorageConsumerLister interface {
	List(selector labels.Selector) (ret []*ocsv1alpha1.StorageConsumer, err error)
}

type storageConsumerLister struct {
	indexer cache.Indexer
}

func NewStorageConsumerLister(indexer cache.Indexer) StorageConsumerLister {
	return &storageConsumerLister{indexer: indexer}
}

func (s *storageConsumerLister) List(selector labels.Selector) (ret []*ocsv1alpha1.StorageConsumer, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*ocsv1alpha1.StorageConsumer))
	})
	return ret, err
}
