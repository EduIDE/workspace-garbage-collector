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
	defaultTTL           = 14 * 24 * time.Hour
	defaultCheckInterval = 30 * time.Minute
	defaultNamespace     = "theia-prod"
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

// getDuration reads a duration from the environment. Values without a unit
// ("1209600") count as seconds, matching the README.
func getDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		seconds, convErr := strconv.ParseInt(value, 10, 64)
		if convErr != nil {
			return 0, fmt.Errorf("invalid %s %q: use a duration like 30m or a number of seconds", key, value)
		}
		d = time.Duration(seconds) * time.Second
	}
	if d <= 0 {
		return 0, fmt.Errorf("invalid %s %q: must be positive", key, value)
	}
	return d, nil
}

func main() {
	var config *rest.Config
	var err error

	namespace := getEnv("K8S_NAMESPACE", defaultNamespace)
	checkInterval, err := getDuration("CHECK_INTERVAL", defaultCheckInterval)
	if err != nil {
		panic(err.Error())
	}

	maxDuration, err := getDuration("WORKSPACE_TTL", defaultTTL)
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
	fmt.Printf("- Workspace TTL: %s\n", maxDuration)

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
			err := garbageCollectWorkspaces(clientSet, namespace, maxDuration)
			if err != nil {
				fmt.Printf("Error during garbage collection: %v\n", err)
			} else {
				fmt.Println("Garbage collection completed successfully")
			}
			fmt.Printf("Will check again at %s\n\n", time.Now().Add(checkInterval).Format(time.DateTime))
		}
	}
}

func garbageCollectWorkspaces(clientSet *clientV1Beta5.V1Beta5Client, namespace string, maxDuration time.Duration) error {
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
