package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func main() {
	// fetch the existing &appsv1.Deployment via API
	expected := appsv1.Deployment{}

	// fill in the fields to generate your expected state

	scheme.Scheme.Default(&expected)
	// now you should have your empty values filled in

	// Serializer = Decoder + Encoder.
	serializer := jsonserializer.NewSerializerWithOptions(
		jsonserializer.DefaultMetaFactory, // jsonserializer.MetaFactory
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectCreater
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectTyper
		jsonserializer.SerializerOptions{
			Yaml:   true,
			Pretty: false,
			Strict: false,
		},
	)

	// Typed -> YAML
	// Runtime.Encode() is just a helper function to invoke Encoder.Encode()
	yaml, err := runtime.Encode(serializer, &expected)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Serialized:\n%s", string(yaml))
}
