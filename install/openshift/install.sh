namespace=$1
oc create -f dc.json
oc create -f sa.json
oc create -f svc.json
oc create -f route.json
oc policy add-role-to-user edit system:serviceaccount:${namespace}:mcp-standalone -n ${namespace}
