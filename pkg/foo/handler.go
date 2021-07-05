package foo

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"google.golang.org/grpc"
	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	devicePluginVersion = "v1beta1"
	devicePluginsDir = "/var/lib/kubelet/device-plugins"
	socketName = "sample-device-plugin.sock"
	resourceName = "hardware-vendor.example/foo"
)

func Serve() error {
	socket := filepath.Join(devicePluginsDir, socketName)
	_ = os.Remove(socket)
	listener, err := net.Listen("unix", socket)
	if err != nil {
		return fmt.Errorf("failed to listen %s: %s", socket, err)
	}
	defer listener.Close()

	server := grpc.NewServer()
	v1beta1.RegisterDevicePluginServer(server, &FooDevicePluginServer{})
	return server.Serve(listener)
}

func Register() error {
	kubeletSocket := filepath.Join(devicePluginsDir, "kubelet.sock")
	conn, err := grpc.Dial("unix://"+kubeletSocket, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to dial %s: %s", kubeletSocket, err)
	}
	defer conn.Close()

	client := v1beta1.NewRegistrationClient(conn)
	req := &v1beta1.RegisterRequest{
		Version:      devicePluginVersion,
		Endpoint:     socketName,
		ResourceName: resourceName,
	}
	_, err = client.Register(context.Background(), req)
	return err
}

func WatchKubeletRestart(stop chan struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %s", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Remove == 0 {
					continue
				}
				if fileName, err := filepath.Rel(devicePluginsDir, event.Name); err == nil && fileName == socketName {
					log.Printf("socket has been removed, exiting")
					os.Exit(0)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("watcher error: %s", err)
			}
		}
	}()

	if err := watcher.Add(devicePluginsDir); err != nil {
		return fmt.Errorf("failed to start watching %s: %s", devicePluginsDir, err)
	}

	<-stop
	return nil
}
