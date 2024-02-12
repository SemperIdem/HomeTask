### Number Processing Service
Handle requests with query parameters that contain URL's, then process and request these URLs, if their JSON response contains "numbers", retrieve it and process
![Test drawio](https://github.com/SemperIdem/homeTask/assets/24257533/29855d4c-b62b-4f10-86bf-0a16591f99f7)



### Running
```
make run 
make test
```

or
```
docker build -t test_image 
docker run -d --name test_container test_image
```
If you run by Docker make sure you configured the network setting correctly to request external APi correctly.


### Example of Usage
```
curl --request GET 'http://127.0.0.1:8080/numbers?u=http://localhost:8090/fibo&u=http://localhost:8090/primes&u=http://localhost:8090/odd'
```

Where ```http://127.0.0.1:8080/``` - endpoint of the service

```'?u=http://localhost:8090/fibo&u=http://localhost:8090/primes&u=http://localhost:8090/odd'``` - external API urls


Output
```
{"numbers":[0,1,2,3,5,7,8,11,13,17,19,21,23,29,31,34,37,39,41,43,45,47,49,51,53,55,57,59,61,63,65,67,69,71,73,75,77,79,81,83,85,87,89,97,144,233,377,610,987,1597,2584,4181]}
```

More examples you can run from requests.http 


### Requirements 

• Write an HTTP service that exposes an endpoint "/numbers". This endpoint receives a list of URLs through a "GET" query parameter. In the example below the parameter is called "u", but you can design it as you see fit. - **Done**</br>

• Write unit tests for your code - **Done**, covered by unit tests critical functionality: handler and controller</br>

• The endpoint needs to return the result as quickly as possible, but always within 500 milliseconds - **Done**, using context.WithTimeOut</br>

• All URLs that were successfully retrieved within the given Fmeframe must influence the result of the endpoint. - **Done**</br>

• If one URL takes longer to respond, keep loading it in the background and cache the response for future use. - **Done**, implemented MemoryCache into HTTPClient which cached with set TTL succesful requests</br>

• It is valid to return an empty list as a result only if all URLs returned errors or took too long to respond and no previous response is stored in the cache - **Done**</br>

