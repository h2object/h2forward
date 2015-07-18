api
---

h2forward api definitions:

-	get all forwards

	method: get
	uri:	/forward/list

	data:
			host:url

-	add a new forward

	method: post
	uri:    /forward/new

	data:
			host:url
				
-	del an old forward

	method: delete
	uri: 	/forward/:host

second, should implement the rate limit control for the client accesses.

