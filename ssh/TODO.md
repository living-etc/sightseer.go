# TODO
## Functionality to implement
- Add charmbracelet log library and make it easy to enable debug mode
- Rewrite the parsing logic to not use regex, instead having each query type
  implement its own parsing logic using Go functions instead of regex
## Research
- Research the scheme for systemctl output. The difference between "preset" and
  "vendor preset" tripped me up - ssh service only had preset, but kubelet had
  "vendor preset".
## Resources to implement
- [] cgroup
- [x] command
- [] cron
- [] default_gateway
- [] docker_container
- [] docker_image
- [x] file
- [] group
- [] interface
- [] ip6tables
- [] ipfilter
- [] ipnat
- [] iptables
- [] kernel_module
- [] linux_audit_system
- [x] linux kernel parameter
- [] lxc
- [] package
- [] port
- [] ppa
- [] process
- [] routing table
- [] selinux
- [] selinux module
- [] systemd
    - [] Handle the different load states
    - [x] service
    - [x] timer
    - [] socket
    - [] device
    - [] mount
    - [] automount
    - [] swap
    - [] target
    - [] path
    - [] snapshot
    - [] slice
    - [] scope
- [x] user
- [] yumrepo
- [] zfs
