// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubernetes

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v1 "k8s.io/api/core/v1"
	api "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/types"
)

func expandEndpointSliceEndpoints(in *schema.Set) []api.Endpoint {
	if in == nil || in.Len() == 0 {
		return []api.Endpoint{}
	}
	endpoints := make([]api.Endpoint, in.Len())
	for i, endpoint := range in.List() {
		r := api.Endpoint{}
		endpointConfig := endpoint.(map[string]interface{})
		if v := endpointConfig["addresses"].([]interface{}); len(v) != 0 {
			r.Addresses = expandStringSlice(v)
		}
		if v, ok := endpointConfig["conditions"].(api.EndpointConditions); ok {
			r.Conditions = v
		}
		if v, ok := endpointConfig["hostname"].(string); ok && v != "" {
			r.Hostname = ptrToString(v)
		}
		if v, ok := endpointConfig["node_name"].(string); ok && v != "" {
			r.NodeName = ptrToString(v)
		}
		if v, ok := endpointConfig["target_ref"].(v1.ObjectReference); ok {
			r.TargetRef = &v
		}
		if v, ok := endpointConfig["zone"].(string); ok && v != "" {
			r.Zone = ptrToString(v)
		}

		endpoints[i] = r
	}
	return endpoints
}

func expandObjectReference(l []interface{}) *v1.ObjectReference {
	if len(l) == 0 || l[0] == nil {
		return &v1.ObjectReference{}
	}
	in := l[0].(map[string]interface{})
	obj := &v1.ObjectReference{}

	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}
	if v, ok := in["namespace"].(string); ok {
		obj.Namespace = v
	}
	if v, ok := in["resource_version"].(string); ok {
		obj.ResourceVersion = v
	}
	if v, ok := in["uid"]; ok {
		obj.UID = types.UID(v.(string))
	}
	if v, ok := in["field_path"].(string); ok {
		obj.FieldPath = v
	}

	return obj
}

func expandEndpointSlicePorts(in *schema.Set) []api.EndpointPort {
	if in == nil || in.Len() == 0 {
		return []api.EndpointPort{}
	}
	ports := make([]api.EndpointPort, in.Len())
	for i, port := range in.List() {
		r := api.EndpointPort{}
		portCfg := port.(map[string]interface{})
		if v, ok := portCfg["name"].(string); ok {
			r.Name = ptrToString(v)
		}
		if v, ok := portCfg["port"].(int32); ok {
			r.Port = &v
		}
		if v, ok := portCfg["protocol"].(v1.Protocol); ok {
			r.Protocol = &v
		}
		if v, ok := portCfg["app_protocol"].(string); ok {
			r.AppProtocol = ptrToString(v)
		}
		ports[i] = r
	}
	return ports
}

func flattenEndpointSliceEndpoints(in []api.Endpoint) *schema.Set {
	att := make([]interface{}, len(in), len(in))
	for i, e := range in {
		m := make(map[string]interface{})
		if e.Hostname != nil {
			m["hostname"] = e.Hostname
		}
		if e.NodeName != nil {
			m["node_name"] = e.NodeName
		}
		if e.Zone != nil {
			m["zone"] = e.Zone
		}
		if len(e.Addresses) != 0 {
			m["addresses"] = e.Addresses
		}
		if e.TargetRef != nil {
			m["target_ref"] = e.TargetRef
		}
		if &e.Conditions != nil {
			m["hostname"] = e.Hostname
		}
		att[i] = m
	}
	return schema.NewSet(hashEndpointSliceEndpoints(), att)
}

func flattenEndpointSlicePorts(in []api.EndpointPort) *schema.Set {
	att := make([]interface{}, len(in), len(in))
	for i, e := range in {
		m := make(map[string]interface{})
		if *e.Name != "" {
			m["name"] = e.Name
		}
		if e.Port != nil {
			m["port"] = int(*e.Port)
		}
		if e.Protocol != nil {
			m["protocol"] = string(*e.Protocol)
		}
		if e.AppProtocol != nil {
			m["app_protocol"] = string(*e.AppProtocol)
		}
		att[i] = m
	}
	return schema.NewSet(hashEndpointSlicePorts(), att)
}

func flattenObjectReference(in *v1.ObjectReference) []interface{} {
	att := make(map[string]interface{})
	if in.Name != "" {
		att["name"] = in.Name
	}
	if in.Name != "" {
		att["namespace"] = in.Name
	}
	if in.FieldPath != "" {
		att["field_path"] = in.Name
	}
	if in.ResourceVersion != "" {
		att["resource_version"] = in.Name
	}
	if in.UID != "" {
		att["uid"] = in.Name
	}

	return []interface{}{att}
}
