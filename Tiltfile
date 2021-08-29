
allow_k8s_contexts(os.getenv("TILT_ALLOW_CONTEXT"))

k8s_yaml('setup.yaml')
k8s_yaml('verify.yaml')

k8s_resource(
  workload='network-health-admission-controller',
  port_forwards=[
     port_forward(local_port=8080, container_port=8080, name='metrics', link_path='/metrics'), 
     port_forward(local_port=9444, container_port=9443, name='readinessProbe', link_path='/healthz'), 
     port_forward(local_port=9443, container_port=9443, name='webhook', link_path='/webhook'), 
  ]
)

k8s_resource('network-health-webhook-available', resource_deps=['network-health-admission-controller'])
k8s_resource('network-health-test-pod', resource_deps=['network-health-webhook-available'])

target='prod'
live_update=[]
if os.environ.get('PROD', '') ==  '':
  target='build-env'
  live_update=[
    sync('main.go', '/app/main.go'),
    sync('mutatingwebhook.go', '/app/mutatingwebhook.go'),
    sync('go.mod', '/app/go.mod'),
    sync('go.sum', '/app/go.sum'),
    run('go install .'),
  ]

docker_build(
  ref='ghcr.io/damoon/network-health-admission-controller:latest',
  context='.',
  dockerfile='./Dockerfile',
  live_update=live_update,
  target=target,
  only=[ 'go.mod'
       , 'go.sum'
       , 'main.go'
       , 'mutatingwebhook.go'
       , 'entrypoint.sh'
  ],
  ignore=[ '.git'
         , '*/*_test.go'
         , '*.yaml'
  ],
)