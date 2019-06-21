# tekton-listener
The project is an addon for [tektoncd-pipeline](https://github.com/tektoncd/pipeline)

The [tektoncd-pipeline](https://github.com/tektoncd/pipeline) is a k8s native CI/CD tool, it can create pipeline for tasks processing,
but there is no trigger to kick off a pipeline unless use `kubectl create -f ` to create a `CRD` named `pipelinerun`.

It's not good, so this project is try to build a `triger`.

The basic idea is:
- Leverage `Knative/event-contrib`(https://github.com/knative/eventing-contrib) to be a event adepter to introduce different event source, such as `github`, `kafka` .etc.
- `tekton-listener` will create a `knative service` when a `EventBinding` is applied. The `KSV` will standby for eventing from `knative/event-contrib`
and handle the event when received.
- When received event, it will get template from `ListenerTemplate` and bind the event, then create `PipelineResource` and `PipelineRun` to trigger the `Pipeline`
- Event should be `CloudEvent`(https://github.com/cloudevents/spec)


For now (2019.6.21) the project is still not complete...
