package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func create(obj runtime.Object, gvk schema.GroupVersionKind) error {
	// find GVR corresponding to given GVK, can use: meta.UnsafeGuessKindToResource(schemaGvk)
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("error mapping gvk=%v to gvr, error=%v", gvk, err)
	}

	// convert runtime.Object to unstructured object
	unstructuredObj, err := getUnstructured(obj)
	if err != nil {
		return err
	}
	namespace, name := unstructuredObj.GetNamespace(), unstructuredObj.GetName()
	log.Debugf("processing object %s/%s, gvk=%v, gvr= %v", namespace, name, gvk, mapping.Resource)

	// check if it already exists
	_, err = dynClient.Resource(mapping.Resource).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("unable to get k8s object %s/%s, error %v", namespace, name, err)
	}

	// create given k8s object if it doesn't exist
	if errors.IsNotFound(err) {
		_, err = dynClient.Resource(mapping.Resource).Namespace(namespace).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to create k8s object %s/%s, error %v", namespace, name, err)
		}
		log.Infof("successfully created k8s object %s/%s", namespace, name)
		return nil
	}

	return fmt.Errorf("k8s object %s/%s, gvk=%v, already exists", namespace, name, gvk)
}

func update(obj runtime.Object, gvk schema.GroupVersionKind) error {
	// find GVR corresponding to given GVK, can use: meta.UnsafeGuessKindToResource(schemaGvk)
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("error mapping gvk=%v to gvr, error=%v", gvk, err)
	}

	// convert runtime.Object to unstructured object
	unstructuredObj, err := getUnstructured(obj)
	if err != nil {
		return err
	}
	namespace, name := unstructuredObj.GetNamespace(), unstructuredObj.GetName()
	log.Debugf("processing object %s/%s, gvk=%v, gvr= %v", namespace, name, gvk, mapping.Resource)

	// fetch current state of an object
	currentObj, err := dynClient.Resource(mapping.Resource).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("unable to get k8s object %s/%s, error %v", namespace, name, err)
	}

	// create given k8s object if it doesn't exist
	if errors.IsNotFound(err) {
		return create(obj, gvk)
	}

	// update existing k8s object for given config
	unstructuredObj.SetResourceVersion(currentObj.GetResourceVersion()) // latest resource version required for updating k8s object
	_, err = dynClient.Resource(mapping.Resource).Namespace(namespace).Update(context.TODO(), unstructuredObj, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("unable to update k8s object %s/%s, error %v", namespace, name, err)
	}

	log.Infof("successfully updated k8s object %s/%s", namespace, name)
	return nil
}

func delete(obj runtime.Object, gvk schema.GroupVersionKind) error {
	// find GVR corresponding to given GVK, can use: meta.UnsafeGuessKindToResource(schemaGvk)
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("error mapping gvk=%v to gvr, error=%v", gvk, err)
	}

	// convert runtime.Object to unstructured object
	unstructuredObj, err := getUnstructured(obj)
	if err != nil {
		return err
	}
	namespace, name := unstructuredObj.GetNamespace(), unstructuredObj.GetName()
	log.Debugf("processing object %s/%s, gvk=%v, gvr= %v", namespace, name, gvk, mapping.Resource)

	// check if it exists
	_, err = dynClient.Resource(mapping.Resource).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("unable to get k8s object %s/%s, error %v", namespace, name, err)
	}

	// create given k8s object if it doesn't exist
	if errors.IsNotFound(err) {
		log.Debugf("object %s/%s, gvk=%v, does not exist, might be deleted", namespace, name, gvk)
		return nil
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	if err = dynClient.Resource(mapping.Resource).Namespace(namespace).Delete(context.TODO(), name, deleteOptions); err != nil {
		return fmt.Errorf("unable to delete k8s object %s/%s, error %v", namespace, name, err)
	}

	return nil
}

func getUnstructured(obj runtime.Object) (*unstructured.Unstructured, error) {
	toUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, fmt.Errorf("unable to convert k8s object to unstructured: %v", err)
	}

	return &unstructured.Unstructured{Object: toUnstructured}, nil
}

func main() {
	initDynamicClient()

}
