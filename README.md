# DoS Attack Demo and Prevention in Go

Demonstrates and provides solution to following DoS (Denial of Service) attacks 
in Go:

- Slowloris
- Large file

The solutions are trivially simple to implement in Go often consisting of a 
simple configuration directive or making use of a standard library function. 

## Installation and Running

Nothing to install.

Run the server:
```
cd cmd
go run ./server.go
```

Run the client
```
cd cmd
go run ./client.go
```

To demo the attacks and run the solutions, see below regarding which lines to
toggle.

## Slowloris

Send requests to the server extremely slow. The notable thing about this DoS
is it takes very little resource in terms of memory or CPU. The goal is to 
send a large number of requests that all send extremely slow requests - thus
making the server use all its connection while waiting on the requests to 
complete. 

### server.go:

To see the attack, keep commented: `//ReadTimeout: 1 * time.Second,`
This will wait on the slow request to complete.

```go
srv := &http.Server{
    Addr:        ":3000",
    Handler:     mux,
    //ReadTimeout: 1 * time.Second,
    //WriteTimeout: 10 * time.Second,
    //IdleTimeout:  1 * time.Minute,
}
```

To see the solution, uncomment: `//ReadTimeout: 1 * time.Second,`
This will limit the read to 1 second.


### client.go

Uncomment the following line: `//postSlow("I've seen...` and run `go run client.go`.
This will make a slow request.

```go
// Slowloris
//postSlow("I've seen things you people wouldn't believe... Attack ships on fire off the shoulder of Orion... I watched C-beams glitter in the dark near the Tannhauser Gate. All those moments will be lost in time, like tears in rain...")
```

You can modify the speed here in the method "func (s *slowloris) Read(p []byte) (int, error)".
Change `time.Sleep(100 * time.Millisecond)`:

```go
func (s *slowloris) Read(p []byte) (int, error) {
    i := 0
    for i < len(p) {
        time.Sleep(100 * time.Millisecond)
```

## Large File

Send very large requests to server. The goal is to overwhelm the server as it
tries to process the large requests often gigabytes in size. In contrast to 
Slowloris, this type of attack requires more memory and CPU to send large files
from the client.

### server.go:

To see the attack, uncomment: `//n, err := io.Copy(io.Discard, r.Body)`
This will process the full request.

To see the solution, comment above and uncomment: 
`//n, err := io.Copy(io.Discard, io.LimitReader(r.Body, 100_000))`
This will limit the reader to process 100K bytes.


```go
// Without limiting
//n, err := io.Copy(io.Discard, r.Body)

// With limiting (solution)
//n, err := io.Copy(io.Discard, io.LimitReader(r.Body, 100_000))
```

### client.go

Uncomment the following line: `//postLarge(1_000_000_000)` and run `go run client.go`
to make a 1GB large request.

```go
// Large file
//postLarge(1_000_000_000)
```

Based on the article and video here:

- https://oak.dev/2024/07/21/dos-attack-demo-and-prevention-in-go/
- https://youtu.be/iZaAOmJGW7s