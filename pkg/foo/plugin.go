package foo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type FooDevicePluginServer struct {
	v1beta1.UnimplementedDevicePluginServer
}

var _ v1beta1.DevicePluginServer = &FooDevicePluginServer{}

func (s *FooDevicePluginServer) GetDevicePluginOptions(ctx context.Context, req *v1beta1.Empty) (*v1beta1.DevicePluginOptions, error) {
	log.Printf("GetDevicePluginOptions, %#v", req)
	return &v1beta1.DevicePluginOptions{}, nil
}

func (s *FooDevicePluginServer) ListAndWatch(req *v1beta1.Empty, stream v1beta1.DevicePlugin_ListAndWatchServer) error {
	log.Printf("ListAndWatch")
	for {
		resp := &v1beta1.ListAndWatchResponse{
			Devices: []*v1beta1.Device{{
				ID:     "b",
				Health: v1beta1.Healthy,
			}, {
				ID:     "a",
				Health: v1beta1.Unhealthy,
			}, {
				ID:     "r",
				Health: v1beta1.Healthy,
			}},
		}
		if err := stream.Send(resp); err != nil {
			return fmt.Errorf("failed to send response: %s", err)
		}
		time.Sleep(time.Second)
	}
}

func (s *FooDevicePluginServer) Allocate(ctx context.Context, req *v1beta1.AllocateRequest) (*v1beta1.AllocateResponse, error) {
	log.Printf("Allocate: %+v", req)
	var containerResps []*v1beta1.ContainerAllocateResponse
	for _, containerReq := range req.ContainerRequests {
		containerResp := &v1beta1.ContainerAllocateResponse{
			Envs: map[string]string{
				"foos": strings.Join(containerReq.DevicesIDs, ","),
			},
		}
		containerResps = append(containerResps, containerResp)
	}
	return &v1beta1.AllocateResponse{
		ContainerResponses: containerResps,
	}, nil
}

