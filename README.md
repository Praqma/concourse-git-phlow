---
maintainer: groenborg, sharor
---
# concourse git phlow


#### using the resource
```yaml
resources:
- name: test
 type: git-phlow
 source:
   prefixready: ready   
   prefixwip: wip
   branch: master
   url: https://github.com/praqma/phlow-test.git
   username: {{github-username}}
   password: {{github-password}}

```

#### params
```yaml
resource_types:
- name: git-phlow
 type: docker-image
 source:
   repository: groenborg/concourse-git-phlow

resources:
- name: test
 type: git-phlow
 source:
   url: https://github.com/praqma/phlow-test.git
   username: {{github-username}}
   password: {{github-password}}

jobs:
- name: pip
 public: true
 plan:
 - get: test
   trigger: true
 - task: ls
   config:
     platform: linux
     image_resource:
       type: docker-image
       source:
         repository: groenborg/concourse-git-phlow
     inputs:
       - name: test
     outputs:
       - name: phlow
     run:
       path: sh
       args:
       - -exc
       - |
         ls
         cd test
         ls
         touch blahblahblah
         git rev-parse HEAD
         ls .git
         cat .git/git-phlow-ready-branch
 - put: test
   params:
     repository: test
```
