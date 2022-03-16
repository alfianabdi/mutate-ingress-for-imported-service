package mutate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "k8s.io/api/admission/v1"
)

func TestMutatesValidRequest(t *testing.T) {
	rawJson := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1",
		"request" : {
			"uid": "29cf0025-ce3d-4d22-8468-8befe9027bd5",
			"kind": {v
				"group": "k8s.nginx.org",
				"kind": "VirtualServer",
				"version": "v1"
			},
			"resource": {
				"group": "k8s.nginx.org",
				"resource": "virtualserver",
				"version": "v1"
			},
			"requestKind": {
				"group": "k8s.nginx.org",
				"kind": "VirtualServer",
				"version": "v1"
			},
			"requestResource": {
				"group": "k8s.nginx.org",
				"resource": "virtualserver",
				"version": "v1"
			},
			"namespace": "p-one-dev",
			"operation": "CREATE",
			"userInfo": {
				"username": "kubernetes-admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"apiVersion": "k8s.nginx.org/v1",
				"kind": "VirtualServer",
				"metadata": {
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"k8s.nginx.org/v1\",\"kind\":\"VirtualServer\",\"metadata\":{\"annotations\":{},\"name\":\"virtualserver-service-a\",\"namespace\":\"p-one-dev\"},\"spec\":{\"host\":\"xtr-0000-p-one-service-a-dev.app-dev.deploud.com\",\"http-snippets\":\"#######################\\n# --- Http Common --- #\\n#######################\\n# --- Configure trusted IP addresses --- #\\nset_real_ip_from 130.176.0.0/17;\\nset_real_ip_from 130.176.128.0/18;\\nset_real_ip_from 130.176.192.0/19;\\nset_real_ip_from 130.176.224.0/20;\\n#real_ip_header X-Forwarded-For;\\n#real_ip_recursive on;\\n\\n#------- Secure URI Path --------\\nmap $request_uri $request_uri_path {\\n    ~(?\\u003cpath\\u003e[^?]+) $path;\\n}\\n\",\"routes\":[{\"action\":{\"pass\":\"certbot\"},\"path\":\"/.well-known/acme-challenge/\"},{\"action\":{\"pass\":\"service-a\"},\"path\":\"/\"}],\"server-snippets\":\"#################################\\n# --- Server Common Before  --- #\\n#################################\\n# --- Allow/disallow query string --- #\\n# --- Default disallow all -----------#\\nif ($args)\\n{\\n  set $ArgumentScheme \\\"${scheme}_args\\\";\\n}\\n####################################\\n# --- Server Website Specific  --- #\\n####################################\\n# --- Set CacheControl --- #\\n# --- Set Querystring Whitelist --- #\\n# --- Add Header --- #    \\n################################\\n# --- Server Common After  --- #\\n################################\\n# --- Set Cache-Control Headers --- #\\nproxy_hide_header Cache-Control;\\nadd_header \\\"Cache-Control\\\" $CacheControl always;\\n\\n# --- Argument Stripping --- #\\nif ($ArgumentScheme ~ \\\"http_(args|blocked)\\\") {\\n    return 301 $scheme://$host$request_uri_path;\\n}\\nif ($ArgumentScheme ~ \\\"https_(args|blocked)\\\") {\\n    return 301 $scheme://$host$request_uri_path;\\n}\\n\",\"upstreams\":[{\"name\":\"certbot\",\"port\":80,\"service\":\"certbot\"},{\"name\":\"service-a\",\"port\":80,\"service\":\"service-a\"}]}}\n"
					},
					"creationTimestamp": "2022-03-10T05:40:37Z",
					"generation": 23,
					"name": "virtualserver-service-a",
					"namespace": "p-one-dev",
					"resourceVersion": "17477616",
					"uid": "29cf0025-ce3d-4d22-8468-8befe9027bd5"
				},
				"spec": {
					"host": "xtr-0000-p-one-service-a-dev.app-dev.deploud.com",
					"http-snippets": "#######################\n# --- Http Common --- #\n#######################\n# --- Configure trusted IP addresses --- #\nset_real_ip_from 130.176.0.0/17;\nset_real_ip_from 130.176.128.0/18;\nset_real_ip_from 130.176.192.0/19;\nset_real_ip_from 130.176.224.0/20;\n#real_ip_header X-Forwarded-For;\n#real_ip_recursive on;\n\n#------- Secure URI Path --------\nmap $request_uri $request_uri_path {\n    ~(?\u003cpath\u003e[^?]+) $path;\n}\n",
					"routes": [
						{
							"action": {
								"pass": "certbot"
							},
							"path": "/.well-known/acme-challenge/"
						},
						{
							"action": {
								"pass": "service-a"
							},
							"path": "/"
						}
					],
					"server-snippets": "#################################\n# --- Server Common Before  --- #\n#################################\n# --- Allow/disallow query string --- #\n# --- Default disallow all -----------#\nif ($args)\n{\n  set $ArgumentScheme \"${scheme}_args\";\n}\n####################################\n# --- Server Website Specific  --- #\n####################################\n# --- Set CacheControl --- #\n# --- Set Querystring Whitelist --- #\n# --- Add Header --- #    \n################################\n# --- Server Common After  --- #\n################################\n# --- Set Cache-Control Headers --- #\nproxy_hide_header Cache-Control;\nadd_header \"Cache-Control\" $CacheControl always;\n\n# --- Argument Stripping --- #\nif ($ArgumentScheme ~ \"http_(args|blocked)\") {\n    return 301 $scheme://$host$request_uri_path;\n}\nif ($ArgumentScheme ~ \"https_(args|blocked)\") {\n    return 301 $scheme://$host$request_uri_path;\n}\n",
					"tls": {
						"secret": "tls-certificate"
					},
					"upstreams": [
						{
							"name": "service-a",
							"port": 80,
							"service": "service-a"
						}
					]
				},
				"status": {
					"externalEndpoints": [
						{
							"ip": "",
							"ports": "[80,443]"
						}
					],
					"message": "Configuration for p-one-dev/virtualserver-service-a was added or updated ",
					"reason": "AddedOrUpdated",
					"state": "Valid"
				}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1"
			}
		}
	}`
	response, err := MutateVirtualServer([]byte(rawJson))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}
	r := v1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"replace","path":"/spec/containers/0/image","value":"debian"}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutateme")
}
