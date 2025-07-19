$ORIGIN internal.curnowtopia.com.
$TTL 60; 1 minute
@ IN SOA ns1.internal.curnowtopia.com. curnowtopia.gmail.com. (
    162       ; serial
    7200      ; refresh (2 hours)
    600       ; retry (10 minutes)
    3600000   ; expire (~42 days)
    60        ; minimum (2 hours)
)
                               NS    ns1.internal.curnowtopia.com.
                               NS    ns2.internal.curnowtopia.com.
ansible                        A     10.2.2.100
atomicpi                       A     10.2.2.25
atomicpi-wifi                  A     10.20.0.75
backup                         CNAME devpi
core-switch                    A     10.0.0.2
debian-template                A     10.2.2.224
devpi                          A     10.2.2.26
devpi-wifi                     A     10.20.0.209
dmzproxy                       A     10.3.0.10
firewall                       CNAME gateway
gateway                        A     10.0.0.1
ha                             CNAME ingress-nginx
hasvr                          A     10.20.0.10
ingress-nginx                  A     10.98.0.2
k8sapi                         A     10.2.2.70
k8sapi1                        A     10.2.2.71
k8sapi2                        A     10.2.2.72
k8sapi3                        A     10.2.2.73
k8sworker1                     A     10.2.2.81
k8sworker2                     A     10.2.2.82
k8sworker3                     A     10.2.2.83
keycloak                       CNAME ingress-nginx
kronk                          A     10.2.2.10
kronk-wifi                     A     10.20.0.70
kvm                            A     10.20.0.210
minio                          CNAME ingress-nginx
minio-console                  CNAME ingress-nginx
mymac                          A     10.20.0.20
ns1                            A     10.0.0.3
ns2                            A     10.0.0.4
octopi                         A     10.20.0.21
octoprint                      CNAME octopi
pool                           A     10.20.0.63
pvectrl                        A     10.2.0.10
pvectrl-c                      A     10.2.1.10
pve1                           A     10.2.0.20
pve1-c                         A     10.2.1.20
pve2                           A     10.2.0.30
pve2-c                         A     10.2.2.30
pve3                           A     10.2.0.40
pve3-c                         A     10.2.1.40
pve4                           A     10.2.0.50
pve4-c                         A     10.2.1.50
pve5                           A     10.2.0.60
pve5-c                         A     10.2.1.60
pve5-mgmt                      A     10.2.0.61
pve6                           A     10.2.0.70
pve6-c                         A     10.2.1.70
talosmgr                       A     10.2.2.101
yzma                           A     10.2.2.20
yzma-wifi                      A     10.20.0.94
