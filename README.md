In this challenge you will have to use what we learned with Multithreading and APIs to get the fastest result between two different APIs.

The two requests will be made simultaneously to the following APIs:

https://cdn.apicep.com/file/apicep/" + cep + ".json

http://viacep.com.br/ws/" + cep + "/json/

The requirements for this challenge are:

- Accept the API that delivers the fastest response and discard the slowest response.

- The result of the request should be displayed on the command line, as well as which API sent it.

- Limit response time to 1 second. Otherwise, the timeout error should be displayed.