api
---

h2forward api definitions:

-	get all forwards

	method: get
	uri:	/virtualhost.json?host=


-	add a new forward

	method: put
	uri:    /virtualhost.json

	data:
			{
				"host":"www.example.com",
				"url":"http://127.0.0.1:8080"
			}
				
-	del an old forward

	method: delete
	uri: 	/virtualhost.json?host=


