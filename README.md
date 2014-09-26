# Auto Scaler

This is a experiment to build a [Cloud Foundry](http://cloudfoundry.org/) [Service](http://docs.cloudfoundry.org/services/overview.html) that monitor apps, and scale them up or down depending on their CPU usage.

Its a naive aproach, what it does is just creates more instances when CPU > 90%, and deletes some when the CPU is < 40%, its just a prototye!

It uses [GoBro](https://github.com/killfill/go-bro) as a lib, to expose itself as a broker in a Cloud Foundry instalation.

# How to use

Clone the repo, push it as an app into Cloud Foundry and:

## Enable the service
```BASH
$ cf create-service-broker scale us3r passw0rd http://auto-scaler.domain.com
Creating service broker scale as admin...
OK
```

```BASH
$ cf service-access
getting service access as admin...
broker: scale
   service       plan   access   orgs
   auto-scaler   cpu    all
```

```BASH
$ cf enable-service-access auto-scaler
Enabling access to all plans of service auto-scaler for all orgs as admin...
OK
```

```BASH
$ cf m
Getting services from marketplace in org virtu / space dev as admin...
OK

service           plans                                         description
auto-scaler       cpu                                           Auto Scaler Experiment
```

## Use the service

User can now create a service and bind apps to it to have them auto-scale:

```BASH
$ cf cs auto-scaler cpu scaler
Creating service scaler in org virtu / space dev as admin...
OK
```

```BASH
$ cf bs myapp scaler
Binding service scaler to app myapp in org virtu / space dev as admin...
OK
TIP: Use 'cf restage' to ensure your env variable changes take effect
```

The app will be monitored.

Unbind the service to undo.

```BASH
$ cf us myapp scaler
Unbinding app myapp from service scaler in org virtu / space dev as admin...
OK
```

