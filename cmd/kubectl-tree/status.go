package main

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog"
)

type ReadyStatus string // True False Unknown or ""
type Reason string

func extractStatus(obj unstructured.Unstructured) (ReadyStatus, Reason, string) {
	jsonVal, _ := json.Marshal(obj.Object["status"])
	klog.V(6).Infof("status for object=%s/%s: %s", obj.GetKind(), obj.GetName(), string(jsonVal))
	statusF, ok := obj.Object["status"]
	if !ok {
		return "", "", ""
	}
	statusV, ok := statusF.(map[string]interface{})
	if !ok {
		return "", "", ""
	}
	conditionsF, ok := statusV["conditions"]
	if !ok {
		return "", "", ""
	}
	conditionsV, ok := conditionsF.([]interface{})
	if !ok {
		return "", "", ""
	}

	for _, cond := range conditionsV {
		condM, ok := cond.(map[string]interface{})
		if !ok {
			return "", "", ""
		}
		condType, ok := condM["type"].(string)
		if !ok {
			return "", "", ""
		}
		if condType == "Ready" || condType == "Succeeded" {
			condStatus, _ := condM["status"].(string)
			condReason, _ := condM["reason"].(string)
			message, _ := condM["message"].(string)

			// perform a manual rewrite for tekton pods
			if condStatus == "False" && condReason == "PodCompleted" {
				return ReadyStatus("True"), Reason(condReason), message
			}
			return ReadyStatus(condStatus), Reason(condReason), message
		}
	}
	return "", "", ""
}
