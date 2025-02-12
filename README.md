[![datree-logo](https://raw.githubusercontent.com/datreeio/datree/main/images/datree_LOGO-180px.png)](#) 

![Travis (.com) branch](https://img.shields.io/travis/com/datreeio/datree/staging?label=build-staging)
![Travis (.com) branch](https://img.shields.io/travis/com/datreeio/datree/main?label=build-main)
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fdatreeio%2Fdatree&count_bg=%2379C83D&title_bg=%23555555&icon=github.svg&icon_color=%23E7E7E7&title=views+%28today+%2F+total%29&edge_flat=false)](https://hits.seeyoufarm.com)
![open issues](https://img.shields.io/github/issues-raw/datreeio/datree)

## What is Datree?
[Datree](https://datree.io/?utm_source=github&utm_medium=organic_oss) helps to prevent Kubernetes misconfigurations from ever making it to production.  

The CLI integration can be used locally or in the CI to provide a policy enforcement solution for Kubernetes. The policy runs automatic checks on every invocation (e.g. pull request, code change, etc.) for rule violations and misconfigurations. When rule violations are found, Datree will prevent them and will show the Engineers instructions on how and why to fix those issues.

## Quick start in two steps
#### 1. Install the latest release on your CLI  
**Linux & MacOS:** ``curl https://get.datree.io | /bin/bash``  
**Windows:** ``not supported yet :(``  

#### 2. Pass datree a Kuberntes manifest file to scan
``datree test <k8s-manifest-file>``  

...and voilà, you just made your first invocation! 🥳    
In your CLI, you will see something like that:  
[![datree-cli-output](https://raw.githubusercontent.com/datreeio/datree/main/images/CLI-output.png)](#) 

#### Ready to review our "Getting Started" guide?
All the information that is needed to explore a bunch of other cool feature, or how to set up your policy, can be found in [**our docs**](https://hub.datree.io/getting-started/?utm_source=github&utm_medium=organic_oss).

## Built-in rules
Right now, there are 30 battle-tested rules for you to choose from.  
The rules are covering different Kubernetes resources/use-cases:
* [Workload](https://hub.datree.io/workload/?utm_source=github&utm_medium=organic_oss)
* [CronJob](https://hub.datree.io/cronjob/?utm_source=github&utm_medium=organic_oss)
* [Containers](https://hub.datree.io/containers/?utm_source=github&utm_medium=organic_oss)
* [Networking](https://hub.datree.io/networking/?utm_source=github&utm_medium=organic_oss)
* [Deprecation](https://hub.datree.io/deprecation/?utm_source=github&utm_medium=organic_oss)
* [Other](https://hub.datree.io/other/?utm_source=github&utm_medium=organic_oss)

## Support

[Datree](https://datree.io/?utm_source=github&utm_medium=organic_oss) builds and maintains this project to make Kubernetes policies simple and accessible.  
Start with our [documentation](https://hub.datree.io/?utm_source=github&utm_medium=organic_oss) for quick tutorials and examples.  
If you need direct support you can contact us at support@datree.io.
