FROM ansibleplaybookbundle/apb-base

LABEL "com.redhat.apb.version"="0.1.0"
LABEL "com.redhat.apb.spec"=\
"bmFtZTogbWNwLWFwYgppbWFnZTogZmVlZGhlbnJ5L21jcC1hcGIKZGVzY3JpcHRpb246IFRoaXMg\
aXMgYSBzYW1wbGUgYXBwbGljYXRpb24gZ2VuZXJhdGVkIGJ5IGFwYiBpbml0CmJpbmRhYmxlOiBG\
YWxzZQphc3luYzogb3B0aW9uYWwKdGFnczogCiAgLSBtb2JpbGUtc2VydmljZQptZXRhZGF0YToK\
ICBkaXNwbGF5TmFtZTogbWNwCnBsYW5zOgogIC0gbmFtZTogZGVmYXVsdAogICAgZGVzY3JpcHRp\
b246IFRoaXMgZGVmYXVsdCBwbGFuIGRlcGxveXMgbWNwLWFwYgogICAgZnJlZTogVHJ1ZQogICAg\
bWV0YWRhdGE6IHt9CiAgICBwYXJhbWV0ZXJzOiBbXQ=="

COPY playbooks /opt/apb/actions
COPY roles /opt/ansible/roles
RUN chmod -R g=u /opt/{ansible,apb}
USER apb
