# servicesresolver

Since Zuul proxy is not supporting websockets yet, and inside our application we need to use websockets, we need a way to talk in some cases directly to our micro services from the ui sample application (create websocket connection in dashboard mask). We need therefore a service, which will be capable of grabbing metadata about registered micro services from eureka's discovery service, so that those data could be used directly inside client.

---

Part of xdevices application - simple micro services-based system for monitoring. <br/>
**Developed for udemy course "Beginners guide to microservices with Go, Spring and RaspPi".**
