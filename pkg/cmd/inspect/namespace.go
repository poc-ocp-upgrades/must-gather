package inspect

import (
	"fmt"
	"log"
	"os"
	"path"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/errors"
)

func namespaceResourcesToCollect() []schema.GroupResource {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []schema.GroupResource{{Resource: "all"}, {Resource: "events"}, {Resource: "configmaps"}, {Resource: "secrets"}}
}
func (o *InspectOptions) gatherNamespaceData(baseDir, namespace string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	log.Printf("Gathering data for ns/%s...\n", namespace)
	destDir := path.Join(baseDir, namespaceResourcesDirname, namespace)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return err
	}
	ns, err := o.kubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}
	ns.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Namespace"))
	errs := []error{}
	filename := fmt.Sprintf("%s.yaml", namespace)
	if err := o.fileWriter.WriteFromResource(path.Join(destDir, "/"+filename), ns); err != nil {
		errs = append(errs, err)
	}
	log.Printf("    Collecting resources for namespace %q...\n", namespace)
	resourcesTypesToStore := map[schema.GroupVersionResource]bool{corev1.SchemeGroupVersion.WithResource("pods"): true}
	resourcesToStore := map[schema.GroupVersionResource]runtime.Object{}
	for gvr := range resourcesTypesToStore {
		list, err := o.dynamicClient.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		resourcesToStore[gvr] = list
	}
	log.Printf("    Gathering pod data for namespace %q...\n", namespace)
	for _, pod := range resourcesToStore[corev1.SchemeGroupVersion.WithResource("pods")].(*unstructured.UnstructuredList).Items {
		log.Printf("        Gathering data for pod %q\n", pod.GetName())
		structuredPod := &corev1.Pod{}
		runtime.DefaultUnstructuredConverter.FromUnstructured(pod.Object, structuredPod)
		if err := o.gatherPodData(path.Join(destDir, "/pods/"+pod.GetName()), namespace, structuredPod); err != nil {
			errs = append(errs, err)
			continue
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("one or more errors ocurred while gathering pod-specific data for namespace: %s\n\n    %v", namespace, errors.NewAggregate(errs))
	}
	return nil
}
