static_resources:
  listeners:
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 14443
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
            {{- range $name, $ports := . }}
            {{- range $ports }}
            - name: "{{ $name }}:{{ . }}"
              domains:
              - "{{ $name }}:{{ . }}"
              {{- if eq . 80 }}
              - "{{ $name }}"
              {{- end }}
              routes:
              - match:
                  prefix: /
                route:
                  cluster: "{{ $name }}:{{ . }}"
            {{- end }}
            {{- end }}
          http_filters:
          - name: envoy.router
            typed_config: {}
  clusters:
  {{- range $name, $ports := . }}
  {{- range $ports }}
  - name: "{{ $name }}:{{ . }}"
    connect_timeout: 0.25s
    type: strict_dns
    lb_policy: round_robin
    load_assignment:
      cluster_name: "{{ $name }}:{{ . }}"
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: {{ $name }}
                port_value: {{ . }}
  {{- end }}
  {{- end }}
admin:
  access_log_path: "/dev/stdout"
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 8001
