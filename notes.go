package tcptohttpgo

/*

we will use RFC 9110 and RFC 9112 to implement HTTP


in HTTP 1.1 we send headers (the same headers all the time or large cookies)


in HTTP 2 and 3 we use HPACK and QPACK (to shrink the headers)


HTTP message = everything is called an http message

METHOD /resource-path PROTOCOL-VERSION\r\n  -> GET /cats HTTP/1.1
field-name: value\r\n                       -> HOST: mewtown.catnet\r\n
field-name: value\r\n                       -> User-Agent: MysticOwl/3.14\r\n
field-name: value\r\n                       -> Accept: application/json\r\n
\r\n                                        -> headers end body folows
{
	"observer": "Whisker Watcher",
	"favourite_breed": "Maine Coon"
}



*/
