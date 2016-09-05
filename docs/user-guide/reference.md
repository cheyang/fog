# configuration file reference

The way of fog is creating a number of machines with IAAS API, then handing on ansible(in docker or not) to deploy the docker cloud(infrastructure) you need.

The configuration file is `yaml` file defining `vmspecs`, `Run` or `DockerRun`, `Global`

## vmspec configuration reference

This section contains a list of all configuration options supported by a vmspec definition.

### name

The name of the VM spec, after instantiation, the vm instance's name will be like `name-0`, `name-1`. This name must be unique

```yaml
Name: kube-master
```

### instances

The number of instances which is created from the VM spec.

```yaml
Instances: 2
```

### Properties

The properties are the input of [docker machine](https://docs.docker.com/machine/drivers/)

Take openstack as example, it should be like

```yaml
Properties: 
  aliyunecs-access-key-id: abc
  aliyunecs-access-key-secret: ecd
  aliyunecs-image-id: centos7u2_64_40G_cloudinit_20160520.raw
  aliyunecs-instance-type: ecs.n1.small
  aliyunecs-internet-max-bandwidth: 100
  aliyunecs-private-address-only: false
  aliyunecs-region: cn-hongkong
  aliyunecs-security-group: k8s
  aliyunecs-system-disk-category: cloud_efficiency
  aliyunecs-io-optimized: optimized
```

### Roles




