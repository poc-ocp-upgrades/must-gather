package util

import (
	"fmt"
	"k8s.io/api/core/v1"
)

func PodRunningReady(p *v1.Pod) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !hasReadyCondition(p) {
		return false, nil
	}
	if p.Status.Phase != v1.PodRunning {
		return false, fmt.Errorf("want pod '%s' on '%s' to be '%v' but was '%v'", p.ObjectMeta.Name, p.Spec.NodeName, v1.PodRunning, p.Status.Phase)
	}
	if !IsPodReady(p) {
		return false, fmt.Errorf("pod '%s' on '%s' didn't have condition {%v %v}; conditions: %v", p.ObjectMeta.Name, p.Spec.NodeName, v1.PodReady, v1.ConditionTrue, p.Status.Conditions)
	}
	return true, nil
}
func hasReadyCondition(pod *v1.Pod) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	conditionReady := true
	for _, cond := range pod.Status.Conditions {
		if cond.Type != v1.PodReady {
			continue
		}
		conditionReady = cond.Status == v1.ConditionTrue
		break
	}
	return conditionReady
}
func IsPodReady(pod *v1.Pod) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return IsPodReadyConditionTrue(pod.Status)
}
func IsPodReadyConditionTrue(status v1.PodStatus) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	condition := GetPodReadyCondition(status)
	return condition != nil && condition.Status == v1.ConditionTrue
}
func GetPodReadyCondition(status v1.PodStatus) *v1.PodCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, condition := GetPodCondition(&status, v1.PodReady)
	return condition
}
func GetPodCondition(status *v1.PodStatus, conditionType v1.PodConditionType) (int, *v1.PodCondition) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if status == nil {
		return -1, nil
	}
	return GetPodConditionFromList(status.Conditions, conditionType)
}
func GetPodConditionFromList(conditions []v1.PodCondition, conditionType v1.PodConditionType) (int, *v1.PodCondition) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if conditions == nil {
		return -1, nil
	}
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return i, &conditions[i]
		}
	}
	return -1, nil
}
