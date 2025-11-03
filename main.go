package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"theia-workspace-garbage-collector/api/types/v1beta5"
	clientV1Beta5 "theia-workspace-garbage-collector/clientset/v1beta5"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const (
	defaultTTL           = int64(14 * 24 * time.Hour) // 14 days in hours
	defaultCheckInterval = int64(30 * time.Minute)    // 30 minutes in minutes
	defaultNamespace     = "theia-prod"               // Default namespace
)

/*
Adapted from https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html
 @ https://github.com/martin-helmich/kubernetes-crd-example/tree/master
*/

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	var config *rest.Config
	var err error

	namespace := getEnv("K8S_NAMESPACE", defaultNamespace)
	checkInterval, err := time.ParseDuration(getEnv("CHECK_INTERVAL", strconv.FormatInt(defaultCheckInterval, 10)))
	if err != nil {
		panic(err.Error())
	}

	config, err = rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	v1beta5.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1Beta5.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Println("Starting garbage collector...")

	// Print configuration
	fmt.Printf("- Namespace: %s\n", namespace)
	fmt.Printf("- Check interval: %s\n", checkInterval)
	fmt.Printf("- Workspace TTL: %s\n", getEnv("WORKSPACE_TTL", strconv.FormatInt(defaultTTL, 10)))

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	// Run garbage collection loop
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down garbage collector...")
			return
		case <-ticker.C:
			fmt.Println("Running garbage collection...")
			err := garbageCollectWorkspaces(clientSet, namespace)
			if err != nil {
				fmt.Printf("Error during garbage collection: %v\n", err)
			} else {
				fmt.Println("Garbage collection completed successfully")
			}
			fmt.Printf("Will check again at %s\n\n", time.Now().Add(checkInterval).Format(time.DateTime))
		}
	}
}

func garbageCollectWorkspaces(clientSet *clientV1Beta5.V1Beta5Client, namespace string) error {
	maxDuration, err := time.ParseDuration(getEnv("WORKSPACE_TTL", strconv.FormatInt(int64(14*24*time.Hour), 10)))
	if err != nil {
		return err
	}

	workspaces, err := clientSet.Workspaces(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	now := time.Now()
	for _, workspace := range workspaces.Items {
		creationTime := workspace.CreationTimestamp.Time
		age := now.Sub(creationTime)

		if age > maxDuration {
			fmt.Printf("Deleting workspace created on %s %s\n", creationTime, workspace.Name)
			err := clientSet.Workspaces(namespace).Delete(workspace.Name, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
