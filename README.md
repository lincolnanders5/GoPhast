# EasyServer
A simple, customizable backend engine written in [Go][golang] by 
[Lincoln Anders][LACOM]. EasyServer uses the [Gin Framework][gin] and a 
backend pattern developed by Lincoln Anders.

## Route Patterns
EasyServer uses common backend route patterns to remove the backend difficulty 
common in rolling custom engines.

### Static Assets
Assets such as images and PDFs and front-end materials such as JavaScript and 
CSS files are all pulled from the [public directory][config].



## Subsites
Subsites may be used when running multiple websites that share 
[static assets][static-assets]. Multiple EasyServer instances may then access 
the shared static assets while presenting different content through dedicated
HTML pages in their subsite directory.

Subsites may be configured with the [`subsite` configuration field][config], 
must be located in the [public directory][config], and must be a directory in 
the format `_[subsite]`.

*Note*: All static assets may be accessed from all subsites, i.e. the images 
used for `subsite1` may be accessed by `subsite2`, given the proper URL is
accessed.

## API Proxying

## Configuration
The following fields are valid in the `config.yml` file:
| Field Name		| Default			| Purpose							|
|-------------------|-------------------|-----------------------------------|
| `name`			| (none)			| If used, displays the configured name for clarity.					|
| `port`			| `4000`			| The port to run EasyServer from 										|
| `release`			| `false`			| Whether to suppress stdio logging or not 								|
| `public`			| `./public`		| The [static asset directory][static-assets] to pull assets from 		|
| `subsite`			| `(none)`			| If used, specifies which [subsite][subsites] directory to pull from 	|
| `api_route`		| `(none)`			| If used, specifies which routes to [proxy to the provided API][api] 	|



[api]: #api-proxying
[config]: #configuration
[static-assets]: #static-assets
[subsites]: #subsites

[gin]: https://github.com/gin-gonic/gin
[golang]: https://golang.org
[LACOM]: https://lincolnanders.com