{{ $metadata := .ObjectMeta }}
{{ $spec := .Spec }}
static_resources:
  listeners:
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 80
    filter_chains:
    - filters:
      - name: envoy.http_connection_manager
        typed_config:
          "@type": type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager
          codec_type: auto
          stat_prefix: ingress_http
          access_log:
          - name: envoy.file_access_log
            config:
              path: "/dev/stdout"
          route_config:
            name: local_route
            virtual_hosts:
            {{- range $index, $port := $spec.RemotePorts }}
            - name: "{{ $metadata.Name }}.{{ $metadata.Namespace }}:{{ $port }}"
              domains:
              - "{{ $metadata.Name }}:{{ $port }}"
              - "{{ $metadata.Name }}.{{ $metadata.Namespace }}:{{ $port }}"
              {{- if eq $port 80 }}
              - "{{ $metadata.Name }}"
              - "{{ $metadata.Name }}.{{ $metadata.Namespace }}"
              {{- end }}
              routes:
              - match:
                  prefix: /
                route:
                  host_rewrite: {{ $spec.RemoteHost }}:{{ $port }}
                  cluster: remote
            {{- end }}
          http_filters:
          - name: envoy.router
            typed_config: {}
  clusters:
  - name: remote
    connect_timeout: 0.25s
    type: strict_dns
    lb_policy: round_robin
    http2_protocol_options: {}
    load_assignment:
      cluster_name: remote
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: {{ .Spec.Tesseract }}
                port_value: 14443
    tls_context:
      common_tls_context:
        tls_certificates:
          certificate_chain: { "filename": "/secret/client.crt" }
          private_key: { "filename": "/secret/client.key" }
        validation_context:
          trusted_ca:
            filename: /secret/ca.crt

admin:
  access_log_path: "/dev/stdout"
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 8001
