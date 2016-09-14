The service should act as a persistent http caching service for microsoft updates.
The files that should be save should have either a specific name, path or content type.
Windows will try to download these files mostly using a 206 range request with many ranges.

Do we want to support a refresh of requests?
Most clients would be windows services that do not force a refresh and will be OK with what recieved.
The other option is to patch the caching service to ignore the no-store or no-cache requests and if required responses.

Request key is: METHOD:url-without-query-terms

Issues with requests fetching is the memory consumption.
To overcome the memory consumption we need to write directly to disk and to not store requests in ram.

The idea was that we can create a disk object with an "in-progress" extention and then download only if both full and in-progress do not exists.
We can use a request queue "dispatcher" to verify that a specific object dosn't need to be fetched again.
We can use a go routing that will handle http requests and will decide on the next action for them.

