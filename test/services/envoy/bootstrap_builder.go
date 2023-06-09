package envoy

import "text/template"

const boostrapText = `
node:
 cluster: ingress
 id: {{.ID}}
 metadata:
  role: {{.Role}}

static_resources:
  clusters:
  - name: xds_cluster
    connect_timeout: 5.000s
    load_assignment:
      cluster_name: xds_cluster
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: {{.GlooAddr}}
                    port_value: {{.Port}}
    http2_protocol_options: {}
    type: STATIC

layered_runtime:
  layers:
  - name: static_layer
    static_layer:
      upstream:
        healthy_panic_threshold:
          value: 0
  - name: admin_layer
    admin_layer: {}

dynamic_resources:
  ads_config:
    transport_api_version: {{ .ApiVersion }}
    api_type: GRPC
    grpc_services:
    - envoy_grpc: {cluster_name: xds_cluster}
  cds_config:
    resource_api_version: {{ .ApiVersion }}
    ads: {}
  lds_config:
    resource_api_version: {{ .ApiVersion }}
    ads: {}

admin:
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: {{.AdminPort}}

`

var bootstrapTemplate = template.Must(template.New("bootstrap").Parse(boostrapText))
