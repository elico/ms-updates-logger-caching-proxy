## Basics

The service should act as a semi-persistent http caching service for microsoft updates.
The files that should be saved should have either a specific name, path or content type.
Windows will try to download these files mostly in as 206 range requests.

Do we want to support a refresh of requests?
Most clients would be windows services that do not force a refresh and will be OK with what is in the cache.
The other option is to patch the caching service to ignore the no-store or no-cache requests and if required responses.

Request key is: METHOD:url-without-query-terms
or
Request key is: hash(sha256 default)(METHOD-string:url_string-without-query-terms)

Issues with requests fetching is the RAM consumption.
To overcome the RAM consumption we need to write directly to disk and to not store requests in RAM.

The idea was that we can create a disk object with an "in-progress" extention and then download it only if both full and in-progress objects do not exists.
We can use a request queue "dispatcher" to verify that a specific object dosn't need to be fetched again.
We can use a GO routin that will handle http requests and will decide on the next action for them.

## Reality
It's not a perfect service but it worth publishing.

Thanks,
Eliezer Croitoru

