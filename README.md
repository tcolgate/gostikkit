# Client for Stikked paste bin

This is a client library, and a command line client, for
the Stikked paste bin (https://github.com/claudehohl/Stikked)

It's main distinguishing features are:

* Support for client side encryption for pasting, and reading
  (compatible with the javascript implementation in the web
  interface)

## Installation

```
go get github.com/tcolgate/gostikkit/cmd/gostikkit
```

# Usage
```
Usage of gostikkit:
  -author string
    	Author of the post
  -encrypt
    	Encrypt the post
  -expire string
    	Expiration time, in minutes (or "never", or "burn")
  -file string
    	Post contents of this file
  -key string
    	API key, if needed
  -lang string
    	The language to render the post as
  -title string
    	Title of the post
  -url string
    	Post contents of this file
```
You can also use ```STIKKED_URL``` and ```STIKKED_KEY``` to pass in the base
url and the api key for your stikked deployment.

