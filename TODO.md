# TODO

- Add an etcdctl command executor to the kubernetes client and refactor the test
  for the etcd cluster members
- Document things, but wait a bit until the design has settled down a bit -
  things are chruning a lot right now
- Research the scheme for systemctl output. The difference between "preset" and
  "vendor preset" tripped me up - ssh service only had preset, but kubelet had
  "vendor preset".
