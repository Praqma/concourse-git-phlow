---
maintainer: groenborg, sharor
---
# Concourse Git Phlow - Tollgate Concourse resource
The Git Phlow is named after a working phlow introduced by Praqma, which is explained significantly better in [the repository where git phlow originates from.](https://github.com/praqma/git-phlow)

For a higher detail you can also visit our blogs at [www.praqma.com/stories](http://www.praqma.com/stories/a-pragmatic-workflow/) alongside other blogs about technology. 

The gist of it is to create a tollgate in front of your main integration branch, usually master, which runs some measure of quality check. 
This is done by pushing a 'ready' branch and letting this resource pick it up. 

For example:

Code => push ready branch => fetch dependencies/build/run fast tests (Concourse job) => on success, push to master.

This way you get a pristine master branch, which will always build and pass the unit test suite. It can also be used to protect other branches further down the line, by changing the ready and master branch pointers.

## Import the resource to your Concourse
Simply add the following under resource_types: 
```yaml
resource_types:
- name: git-phlow
 type: docker-image
 source:
   repository: groenborg/concourse-git-phlow
   tag: '1.0.10'
```

Note: versions below 0.2.31 are not finished and were used during development. 


## Using the resource (Source configuration)
Note: Whenever the tollgate does a rebase and fast-forward the check and in sha will differ. Do not be alarmed. 

Furthermore, for a `ready/` branch to be found by check it need a new unique sha. This means repushing a failed branch will not trigger check. 
- `prefixready`: *Required.* The branch prefix that Concourse uses to find new branches to integrate to tollgated branch. 
- `prefixwip`: *Required.* This the prefix for the branch that Concourse uses while the job runs on the resource. 
  * While a branch has the work in progress prefix, the job is either running or has failed. It is intended as a way to allow the developer to recover their failed branch post fast forward merge.
- `branch`: *Required.* The tollgated branch, usually master or equivalent.
- `url`: *Required.* Url pointing to the git repository.
- `username`: *Required.* Username for logging into the repository.
- `password`: *Required.* Password or token (in the case of Github) for logging into the repository.
- `datadog_api_key`: *Optional.* Datadog api key; required for using monitoring with datadog.
- `datadog_app_key`: *Optional.* Datadog app key; required for using monitoring with datadog.
- `metric_name`: *Optional.* Name of the metric in the datadog interface e.g. `tollgate_example_name`


An example for the resource can be found below: 

```yaml
resources:
- name: tollgate
  type: git-phlow
  source:
    prefixready: ready/   
    prefixwip: wip/
    branch: master
    url: https://github.com/praqma/concourse-phlow-test.git
    username: {{github-username}}
    password: {{github-password-or-token}}
```
As of right now, we do not yet support SSH as a means to authenticate against the git repository. [It is already an issue](https://github.com/Praqma/concourse-git-phlow/issues/11).

## Resource monitoring
Enabling monitoring for the tollgate will monitor the resource itself, not the job or builds. The metric output for datadog is: 
- Event 
  - `Error events` thrown when errors occur with stacktraces along side the error message displayed in concourse 
  - `Info event` sent when a ready branch is picked up and sent to [in](/in/in.go)
- Metric
  - Health check that runs every time check runs to see that the resource is still working properly. It will post `1` everytime a check happens. 
 
##### Example
An example configuration with a datadog setup. To read more about datadog for monitoring look [here](https://www.datadoghq.com/)

```yaml
resources:
- name: tollgate
  type: pip
  source:
    prefixready: ready/
    prefixwip: wip/
    master: master
    url: https://github.com/Praqma/mwriter-git-phlow.git
    username: {{github-username}}
    password: {{github-password}}
    datadog_api_key: {{dd_api_key}}
    datadog_app_key: {{dd_app_key}}
    metric_name: conocurse_git_phlow
```


## Concourse jobs with the tollgate
Since the intended workflow is to get a ready branch, and then perform a quality check followed by a push, most normal jobs only have one input and no outputs. 

It is required to have the tollgate resource as an input, as it stores a file in the .git folder which locates the work in progress branch. 

Should there be outside dependencies, copy them into the tollgate folder and remove them before using out. 

A normal example of the tollgate in action: 

```yaml
jobs:
- name: build-golang
  public: true
  plan:
  - get: test
    trigger: true
  - task: ls
    config:
      platform: linux
      image_resource:
        type: docker-image
        source: {repository: golang, tag: '1.8'}          
      inputs:
        - name: tollgate
      run:
        path: sh
        args:
        - -exc
        - |
          go build tollgate/main.go          
  - put: tollgate
    params:
      repository: tollgate
```
#### Contributing
We communicate through github issues, if there is a feature that you want which does not have an issue, simply create it and we will be in touch. 

To contribute to an existing issue, fork the repository and make a pull request when ready. :wink: 

## CodeScene Analysis report
[![](https://codescene.io/projects/1475/status.svg) Get more details at **codescene.io**.](https://codescene.io/projects/1475/jobs/latest-successful/results)
