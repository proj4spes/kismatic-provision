package aws

type PlanAWS struct {
	Etcd                         []NodeDeets
	Master                       []NodeDeets
	Worker                       []NodeDeets
	MasterNodeFQDN               string
	MasterNodeShortName          string
	SSHUser                      string
	SSHKeyFile                   string
	HomeDirectory                string
	AllowPackageInstallation     bool
	AutoConfiguredDockerRegistry bool
	DockerRegistryIP             string
	DockerRegistryPort           int
	DockerRegistryCAPath         string
	AdminPassword                string
}

const planAWSOverlay = `cluster:
  name: kubernetes
  admin_password: {{.AdminPassword}}      # This password is used to login to the Kubernetes Dashboard and can also be used for administration without a security certificate
  allow_package_installation: true      # When false, installation will not occur if any node is missing the correct deb/rpm packages. When true, the installer will attempt to install missing packages for you.
  networking:
    type: overlay                        # overlay or routed. Routed pods can be addressed from outside the Kubernetes cluster; Overlay pods can only address each other.
    pod_cidr_block: 172.16.0.0/16        # Kubernetes will assign pods IPs in this range. Do not use a range that is already in use on your local network!
    service_cidr_block: 172.17.0.0/16    # Kubernetes will assign services IPs in this range. Do not use a range that is already in use by your local network or pod network!
    policy_enabled: false                # When true, enables network policy enforcement on the Kubernetes Pod network. This is an advanced feature.
    update_hosts_files: true            # When true, the installer will add entries for all nodes to other nodes' hosts files. Use when you don't have access to DNS.
  certificates:
    expiry: 17520h                       # Self-signed certificate expiration period in hours; default is 2 years.
  ssh:
    user: {{.SSHUser}}
    ssh_key: {{.SSHKeyFile}}             # Absolute path to the ssh public key we should use to manage nodes.
    ssh_port: 22
docker_registry:                         # Here you will provide the details of your Docker registry or setup an internal one to run in the cluster. This is optional and the cluster will always have access to the Docker Hub.
  setup_internal: true                  # When true, a Docker Registry will be installed on top of your cluster and used to host Docker images needed for its installation.
  address: ""                            # IP or hostname for your Docker registry. An internal registry will NOT be setup when this field is provided. Must be accessible from all the nodes in the cluster.
  port: 443                              # Port for your Docker registry.
  CA: ""                                 # Absolute path to the CA that was used when starting your Docker registry. The docker daemons on all nodes in the cluster will be configured with this CA.
etcd:
  expected_count: {{len .Etcd}}
  nodes:{{range .Etcd}}
  - host: {{.Hostname}}
    ip: {{.PublicIP}}
    internalip: {{.PrivateIP}}{{end}}
master:
  expected_count: {{len .Master}}
  nodes:{{range .Master}}
  - host: {{.Hostname}}
    ip: {{.PublicIP}}
    internalip: {{.PrivateIP}}{{end}}
  load_balanced_fqdn: {{.MasterNodeFQDN}}
  load_balanced_short_name: {{.MasterNodeShortName}}
worker:
  expected_count: {{len .Worker}}
  nodes:{{range .Worker}}
  - host: {{.Hostname}}
    ip: {{.PublicIP}}
    internalip: {{.PrivateIP}}{{end}}
`
