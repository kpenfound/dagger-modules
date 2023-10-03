# Encircle

Dagger module to run circleci workflows

Example:

`echo "{git(url:'github.com/kpenfound/encircle') {branch(name:'main') {tree {encircleWorkflow(workflow:'test')}}}}" | dagger query -m github.com/kpenfound/dagger-modules/encircle --progress=plain`
