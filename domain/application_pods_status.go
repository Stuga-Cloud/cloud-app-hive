package domain

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
)

type PodList struct {
	ListMetaData any   `json:"listMetaData"`
	Items        []Pod `json:"items"`
}

type PodMetaData struct {
	Name string `json:"name"`
}

type Pod struct {
	MetaData        PodMetaData `json:"metadata"`
	Status          PodStatus   `json:"status"`
	HumanizedStatus string      `json:"humanizedStatus"`
}

type PodStatus struct {
	Phase             string            `json:"phase"`
	Conditions        []PodCondition    `json:"conditions"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses"`
}

type PodCondition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	LastProbeTime      time.Time `json:"lastProbeTime"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Reason             string    `json:"reason"`
	Message            string    `json:"message"`
}

type ContainerStatus struct {
	Name                 string         `json:"name"`
	State                ContainerState `json:"state"`
	LastTerminationState ContainerState `json:"lastTerminationState"`
	Ready                bool           `json:"ready"`
	RestartCount         int32          `json:"restartCount"`
	Image                string         `json:"image"`
	ImageID              string         `json:"imageID"`
	ContainerID          string         `json:"containerID"`
	Started              bool           `json:"started"`
}

type ContainerState struct {
	Waiting    *ContainerStateWaiting    `json:"waiting,omitempty"`
	Running    *ContainerStateRunning    `json:"running,omitempty"`
	Terminated *ContainerStateTerminated `json:"terminated,omitempty"`
}

type ContainerStateWaiting struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type ContainerStateRunning struct {
	StartedAt time.Time `json:"startedAt"`
}

type ContainerStateTerminated struct {
	ExitCode    int32     `json:"exitCode"`
	Signal      int32     `json:"signal"`
	Reason      string    `json:"reason"`
	Message     string    `json:"message"`
	StartedAt   time.Time `json:"startedAt"`
	FinishedAt  time.Time `json:"finishedAt"`
	ContainerID string    `json:"containerID"`
}

func ConvertPods(pods *v1.PodList) PodList {
	podList := PodList{
		ListMetaData: convertAny(pods.ListMeta),
		Items:        make([]Pod, len(pods.Items)),
	}

	for i, pod := range pods.Items {
		podList.Items[i] = convertPod(pod)
	}

	return podList
}

func convertPod(pod v1.Pod) Pod {
	return Pod{
		MetaData: PodMetaData{
			Name: pod.Name,
		},
		Status: convertPodStatus(pod.Status),
	}
}

func convertPodStatus(status v1.PodStatus) PodStatus {
	podStatus := PodStatus{
		Phase:             string(status.Phase),
		Conditions:        make([]PodCondition, len(status.Conditions)),
		ContainerStatuses: make([]ContainerStatus, len(status.ContainerStatuses)),
	}

	for i, condition := range status.Conditions {
		podStatus.Conditions[i] = convertPodCondition(condition)
	}

	for i, containerStatus := range status.ContainerStatuses {
		podStatus.ContainerStatuses[i] = convertContainerStatus(containerStatus)
	}

	return podStatus
}

func convertPodCondition(condition v1.PodCondition) PodCondition {
	return PodCondition{
		Type:               string(condition.Type),
		Status:             string(condition.Status),
		LastProbeTime:      condition.LastProbeTime.Time,
		LastTransitionTime: condition.LastTransitionTime.Time,
		Reason:             condition.Reason,
		Message:            condition.Message,
	}
}

func convertContainerStatus(status v1.ContainerStatus) ContainerStatus {
	containerStatus := ContainerStatus{
		Name:                 status.Name,
		State:                convertContainerState(status.State),
		LastTerminationState: convertContainerState(status.LastTerminationState),
		Ready:                status.Ready,
		RestartCount:         status.RestartCount,
		Image:                status.Image,
		ImageID:              status.ImageID,
		ContainerID:          status.ContainerID,
		Started:              *status.Started,
	}

	return containerStatus
}

func convertContainerState(state v1.ContainerState) ContainerState {
	containerState := ContainerState{}

	if state.Waiting != nil {
		containerState.Waiting = &ContainerStateWaiting{
			Reason:  state.Waiting.Reason,
			Message: state.Waiting.Message,
		}
	}

	if state.Running != nil {
		containerState.Running = &ContainerStateRunning{
			StartedAt: state.Running.StartedAt.Time,
		}
	}

	if state.Terminated != nil {
		containerState.Terminated = &ContainerStateTerminated{
			ExitCode:    state.Terminated.ExitCode,
			Signal:      state.Terminated.Signal,
			Reason:      state.Terminated.Reason,
			Message:     state.Terminated.Message,
			StartedAt:   state.Terminated.StartedAt.Time,
			FinishedAt:  state.Terminated.FinishedAt.Time,
			ContainerID: state.Terminated.ContainerID,
		}
	}

	return containerState
}

func convertAny(in interface{}) any {
	// Implement your own conversion logic here if needed
	// You can use a library like json.Marshal and json.Unmarshal to perform the conversion
	return in
}

func ComputeHumanizedPodStatus(pods *[]Pod) []Pod {
	for _, pod := range *pods {
		pod.HumanizedStatus = fmt.Sprintf("Pod %s is in an unknown state", pod.MetaData.Name)
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Waiting != nil {
				pod.HumanizedStatus = fmt.Sprintf("Pod %s is waiting", pod.MetaData.Name)
				if containerStatus.State.Waiting.Reason == "CrashLoopBackOff" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s is crashing and restarting repeatedly", pod.MetaData.Name)
				}
				if containerStatus.State.Waiting.Reason == "ImagePullBackOff" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s is unable to pull the image", pod.MetaData.Name)
				}
				if containerStatus.State.Waiting.Reason == "ErrImagePull" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s is unable to pull the image", pod.MetaData.Name)
				}
				if containerStatus.State.Waiting.Reason == "CreateContainerConfigError" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s is unable to create the container", pod.MetaData.Name)
				}
				if containerStatus.State.Waiting.Reason == "InvalidImageName" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s has an invalid image name", pod.MetaData.Name)
				}
				if containerStatus.State.Waiting.Reason == "InvalidImageName" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s has an invalid image name", pod.MetaData.Name)
				}
			}
			if containerStatus.State.Terminated != nil {
				pod.HumanizedStatus = fmt.Sprintf("Pod %s has terminated", pod.MetaData.Name)
				if containerStatus.State.Terminated.ExitCode != 0 {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s might have failed", pod.MetaData.Name)
				}
				if containerStatus.State.Terminated.Reason == "ContainerCreating" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s is being created", pod.MetaData.Name)
				}
				if containerStatus.State.Terminated.Reason == "Completed" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s has completed", pod.MetaData.Name)
				}
				if containerStatus.State.Terminated.Reason == "OOMKilled" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s has been killed because it ran out of memory", pod.MetaData.Name)
				}
				if containerStatus.State.Terminated.Reason == "Error" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s has failed", pod.MetaData.Name)
				}
				if containerStatus.State.Terminated.Reason == "DeadlineExceeded" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s has failed", pod.MetaData.Name)
				}
				if containerStatus.State.Terminated.Reason == "NodeLost" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s has failed", pod.MetaData.Name)
				}
				if containerStatus.State.Terminated.Reason == "Evicted" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s has been evicted", pod.MetaData.Name)
				}
				if containerStatus.State.Terminated.Reason == "Killing" {
					pod.HumanizedStatus = fmt.Sprintf("Pod %s is being killed", pod.MetaData.Name)
				}
			}
			if containerStatus.State.Running != nil {
				pod.HumanizedStatus = fmt.Sprintf("Pod %s is running", pod.MetaData.Name)
			}
		}
	}
	return *pods
}
