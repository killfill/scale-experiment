

## TIPS

DOMAIN=run.pivotal.io
curl -k -H "Accept: application/json" https://login.${DOMAIN}/login

curl -k -X POST https://login.${DOMAIN}/oauth/token -H "Accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "grant_type=password" -d "password=errr" -d "scope=" -d "username=userrrr"  -H "Authorization: Basic Y2Y6" -i|json

#^---- That could work with Content-Type: application/json too
#      And with a GET too, sending the params as query strings...

TOKEN=el_token

curl -k https://api.${DOMAIN}/v2/organizations -H "Accept: application/json" -H "Authorization: Bearer $TOKEN"


/v2/apps/48c90b4f-a80e-4961-b59f-eefa7027f44e/summary
/v2/apps/48c90b4f-a80e-4961-b59f-eefa7027f44e/instances #Muestra estado de las instancias
/v2/apps/48c90b4f-a80e-4961-b59f-eefa7027f44e/stats #Stats de las instancias mem, cpu! y el estado.. :P
/v2/apps/:guid/env
PUT /v2/apps/c1b4b87d-cb2b-4669-af77-2de1cea16118?async=true&inline-relations-depth=1  body:   {"instances":2}
