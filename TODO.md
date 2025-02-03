# TODO

- Add an etcdctl command executor to the kubernetes client and refactor the test
  for the etcd cluster members
- Document things, but wait a bit until the design has settled down a bit -
  things are chruning a lot right now

- Add recovers to prevent any panics from escaping the bounds of the module, and
  convert them to an error instead.
- Rewrite error handling to be more idiomatic.
