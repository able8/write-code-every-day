// ## Imports and globals
//
// The package name must be `main` to make the code run in the Go playground. In real life, the package would be a library package named, for example, `auth`.
package main

// None of these packages (except `time`) are actually required for the token refreshing code. They are used by the code that simulates the token refresh API, the test code, and for printing out what's going on.
import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	rnd "math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	// We want to refresh the token *before* it expires. The `lifeSpanSafetyMargin` duration shall provide enough time for this.
	// For the simulation, it is set to an unrealistically small value, to make the test run fast.
	lifeSpanSafetyMargin = 10 * time.Millisecond
	// If the token cannot be refreshed, we want to retry after a short delay.
	retryDelay = 11 * time.Millisecond
)

// The authorization API returns either a token or an error. We collect either of these in a `tokenResponse` and pass the result on to the client.
type tokenResponse struct {
	Token string
	Err   error
}

// `Token` represents an access token. It refreshes itself in the background by calling the API's authorization endpoint before the current token expires.
type Token struct {
	// The `accessToken` channel is used to send the new access token to the client.
	accessToken chan tokenResponse
	// The `authorize` field allows setting a custom authorization function that implements the call to the actual authorization endpoint.
	authorize func() (string, time.Duration, error)
}

// Method `refreshToken` fetches a new access token from the authorization API if there is none yet or if the current one expires. It sends the results (a token or an error) to the `accessToken` channel.
func (a *Token) refreshToken(ctx context.Context) {
	var token string
	var expiration time.Duration
	var err error

	// Set the initial token, before any client can request it.
	// `authorize()` is defined below. Its purpose is to fetch a new token from an authorization endpoint, including the token's lifespan.
	token, expiration, err = a.authorize()

	// Set a new timer to fire when 90% of the expiration duration has passed. We want a new token *before* the current one expires.
	expired := time.After(expiration - lifeSpanSafetyMargin)

	for {
		select {
		// When a client requests a token, this `case` condition writes one to the `accessToken` channel. It does nothing else, hence the body of the case is empty.
		case a.accessToken <- tokenResponse{Token: token, Err: err}:

		// The expiration timer has fired and wrote the current time to `expired`.
		case <-expired:
			// Refresh the token.
			log.Println("Token expired")
			token, expiration, err = a.authorize()
			if err != nil {
				log.Println("Error refreshing token:", err)
				// If the token cannot be fetched, retry frequently instead of waiting for the token's normal timeout (which could be minutes away).
				expiration = retryDelay
			} else {
				log.Println("Token refreshed")
			}
			// Set a new timer to fire when 90% of the expiration duration has passed, or, if the token could not be refreshed, to fire when the retry delay has passed.
			expired = time.After(expiration - lifeSpanSafetyMargin)

		// The context has been canceled. Stop the goroutine.
		case <-ctx.Done():
			return
		}
	}
}

// The Token constructor receives the authorization function to call. It takes care of spawning the goroutine that refreshes the token in the background.
func NewToken(ctx context.Context, auth func() (string, time.Duration, error)) *Token {
	a := &Token{
		accessToken: make(chan tokenResponse),
		authorize:   auth,
	}
	go a.refreshToken(ctx) // This call sets a.token and a.apiErr.
	return a
}

// Method `Get()` returns the current token or an error.
func (a *Token) Get() (string, error) {
	t := <-a.accessToken
	return t.Token, t.Err
}

/*

### Simulating an authorization endpoint

Next, let me implement a flaky authorization function that we can pass to the `Token` constructor.

The function `authFunc()` simulates fetching a new access token that expires after `lifespan` milliseconds.

But the simulated authorization endpoint is not very stable. With a probability of `apiFailureRate`, the call to the endpoint fails, and the failure then persists for `apiErrorDuration`.

When this happens, the function simulates a half-heartedly backoff strategy ("sleep, then try just once again") that fails half of the time. (Exponential backoff with jitter, anyone? Take it as a homework assignment.)

Feel free to skip reading through that part of the code, it is not relevant for the implementation. For any production purposes, you would insert a real API call here.

*/

// These constants help simulate a not very reliable authorization endpoint.
// To finish the tests quickly, the durations are set to absurdly short values.
// (Typically, access tokens of web APIs have lifespans that are counted in
// minutes, not milliseconds.)

const (
	tokenLifeSpan       = 100 * time.Millisecond
	averageCallDuration = 8 * time.Millisecond
	apiFailureRate      = 0.2
	apiErrorDuration    = tokenLifeSpan * 150 / 100
)

// `tempError` is set to `true` during a simulated transient API/network outage. To avoid races, it is an atomic value.
var tempError atomic.Bool

// `authFunc()` simulates fetching a new access token that expires after `lifespan` milliseconds.
func authFunc() (token string, lifespan time.Duration, err error) {
	b := make([]byte, 8)

	_, err = rand.Read(b)
	if err != nil {
		return "", tokenLifeSpan, err
	}

	// Simulate the delay of fetching a new token
	time.Sleep(averageCallDuration)

	// Simulate an API call error with a probability of `apiFailureRate`.
	// The error lasts for `apiErrorDuration`, then disappears.
	// The code pretends to do a backoff strategy that fails half of the time.
	if !tempError.Load() && rnd.Float64() < apiFailureRate {
		log.Println("API error")
		log.Println("Backing off...")
		time.Sleep(tokenLifeSpan)

		if rnd.Float64() < 0.5 {
			// Backoff strategy was not successful
			log.Println("API is still not back, giving up")
			tempError.Store(true)

			// The API/network outage resolves itself after `apiErrorDuration`.
			go func() {
				<-time.After(apiErrorDuration)
				log.Println("API error disappeared")
				tempError.Store(false)
			}()
		} else {
			log.Println("API error disappeared during backoff")
		}
	}

	if tempError.Load() {
		return "", tokenLifeSpan, fmt.Errorf("temporary API error")
	}

	return fmt.Sprintf("%x", b), tokenLifeSpan, nil

}

/*

## Testing

Instead of building a complete web app with an HTTP server, handlers, and more, let's just write a simple test that simulates a few clients requesting new tokens in regular intervals.

*/
//
func TestTokenGet(t *testing.T) {
	log.SetFlags(0) // no extra log info
	log.Println("Starting test")
	ctx, cancel := context.WithCancel(context.Background())
	a := NewToken(ctx, authFunc)
	defer cancel()

	// A test client requests an API token regularly, so that it can call the API.
	client := func(n int, token *Token, done <-chan struct{}, wg *sync.WaitGroup) {
		for {
			select {
			// The clients can be stopped by closing the done channel.
			case <-done:
				wg.Done()
				log.Printf("Client %d done\n", n)
				return
			// The clients shall request a token multiple times during the token's lifespan.
			default:
				t, err := token.Get()
				log.Printf("Client %d token: %s, err: %v\n", n, t, err)
				time.Sleep(tokenLifeSpan / 5)
			}

			time.Sleep(time.Second)
		}
	}

	// Start the clients in separate goroutines.
	var wg sync.WaitGroup
	done := make(chan struct{})
	log.Println("Starting clients")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go client(i, a, done, &wg)
	}

	// After 5 seconds, stop the clients.
	<-time.After(5 * time.Second)
	log.Println("Stopping clients")
	close(done)
	log.Println("Waiting for clients to clean up")
	wg.Wait()
}

/*
## Exploring an alternative using mutexes

How did this article arrive at the presented solution using channels and `select`?

I started from modeling futures in Go, especially, futures that deliver the computed result continuously. The beauty of that approach is that you do not need to share memory that you'd need to protect with mutexes at every access path.

I developed this "continuous future" further into a future that can update itself in the background.

I added a timer to signal when it's time to update the token. I added a `select` statement that listens to the timer and to consumers. This `select` statement takes the role of a `mutex` construct for guarding access to the token data.

How would the same code look like when using mutexes?

### Don't communicate by sharing memory...

...except if the scenario is simple enough. So let's try a version that uses shared memory protected by mutexes.
*/
// The modified token struct uses an RW mutex to serialize write access while allowing readers to read the token concurrently.
type MToken struct {
	token     string
	tokenErr  error
	mu        sync.RWMutex
	authorize func() (string, time.Duration, error)
	ctx       context.Context
}

// Create a new token object.
func NewMToken(ctx context.Context, auth func() (string, time.Duration, error)) *MToken {
	m := &MToken{
		authorize: authFunc,
		ctx:       ctx,
	}
	go m.refreshToken(ctx) // This call sets a.token and a.apiErr.
	return m
}

// Get provides protected read access to `m.token`. Locking the mutex for reading (through `RLock()`) allows mutliple readers to read the protected value at the same time. If a writer requests a write lock through `Lock()`, then `RLock()` does not let any new readers aquire a read lock until the write lock is released.
func (m *MToken) Get() (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.token, m.tokenErr
}

// In `refreshToken`, the `tokenResponse` channel is replaced by a direct, mutex-protected write to `m.token`.
func (m *MToken) refreshToken(ctx context.Context) {
	var expiration time.Duration

	// Set the initial token, before any client can request it.
	m.mu.Lock()
	m.token, expiration, m.tokenErr = m.authorize()
	m.mu.Unlock()

	// Set a new timer to fire when 90% of the expiration duration has passed. We want a new token *before* the current one expires.
	expired := time.After(expiration - lifeSpanSafetyMargin)

	for {
		select {
		// The expiration timer has fired and wrote the current time to `expired`.
		case <-expired:
			// Refresh the token.
			log.Println("Token expired")

			m.mu.Lock()
			// The call to authorize() is inside the lock, so that clients cannot request the current token while it is in the process of getting invalidated. They have to wait for the new one.
			m.token, expiration, m.tokenErr = m.authorize()
			// To avoid unnecessary delay, error logging does not need to be inside the lock. Creating an unshared copy of the error value allows logging it later without a read lock.
			err := m.tokenErr
			m.mu.Unlock()

			if err != nil {
				log.Println("Error refreshing token:", err)
				// If the token cannot be fetched, retry frequently instead of waiting for the token's normal timeout (which could be minutes away).
				expiration = retryDelay
			} else {
				log.Println("Token refreshed")
			}

			// Set a new timer to fire when 90% of the expiration duration has passed.
			expired = time.After(expiration - lifeSpanSafetyMargin)
		case <-ctx.Done():
			return
		}
	}
}

/*

## Test the token version

*/
//
func TestMTokenGet(t *testing.T) {
	log.SetFlags(0) // no extra log info
	log.Println("Starting mutex test")
	ctx, cancel := context.WithCancel(context.Background())
	m := NewMToken(ctx, authFunc)
	defer cancel()

	// A test client requests an API token regularly, so that it can call the API.
	client := func(n int, token *MToken, done <-chan struct{}, wg *sync.WaitGroup) {
		for {
			select {
			// The clients can be stopped by closing the done channel.
			case <-done:
				wg.Done()
				log.Printf("Mutex client %d done\n", n)
				return
			// The clients shall request a token multiple times during the token's lifespan.
			default:
				t, err := token.Get()
				log.Printf("Mutex client %d token: %s, err: %v\n", n, t, err)
				time.Sleep(tokenLifeSpan / 5)
			}
		}
	}

	// Start the clients in separate goroutines.
	var wg sync.WaitGroup
	done := make(chan struct{})
	log.Println("Starting mutex clients")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go client(i, m, done, &wg)
	}

	// After 5 seconds, stop the clients.
	<-time.After(5 * time.Second)
	log.Println("Stopping mutex clients")
	close(done)
	log.Println("Waiting for mutex clients to clean up")
	wg.Wait()
}

/*
## Run the code

The code runs fine in the Go playground and therefore is just [one click away][playground].

Alternatively, clone the [code repository][repository] and run `go test`.

```
git clone https://github.com/appliedgo/refresh
cd refresh
go test -v .
```

## Thoughts on using the code in real life

The code shown here is meant for demo purposes. To use it in production, consider a couple of improvements.

### Add a good backoff strategy

The above code's backoff "strategy" is to sleep once and try again.

A tested-and-proven backup strategy for production code is "exponential backoff with jitter". Exponential means that the time between retries becomes exponentially longer (for example, the first delay is 1 second, the second is 2 seconds the third is 4, the fourth is 8, and so on). Jitter means adding a random amount of time to the delay, to avoid that clients that happen to go into backoff at the same time all retry the call at the same times.


### Which approach is faster, channels or mutexes?

You probably don't need to care. The performance difference between sharing data via guarded access versus sharing data via channels is certainly negligible, given that the shared data—an access token—is used for calling an API endpoint (or even multiple endpoints), which is *orders of magnitudes* slower than reading a string shared via mutex-locking or a channel.

But for the sake of it, I have run benchmarks on Token.Get() and MToken.Get(). (The benchmark code is available [in the repository][repository] only.) The results speak for themselves.

```sh
> go test -run NONE -bench=.
goos: darwin
goarch: arm64
pkg: github.com/appliedgo/refresh
BenchmarkToken_Get-12            4263207               296.5 ns/op
BenchmarkMToken_Get-12           4758291               317.2 ns/op
PASS
ok      github.com/appliedgo/refresh    4.552s
> go test -run NONE -bench=.
goos: darwin
goarch: arm64
pkg: github.com/appliedgo/refresh
BenchmarkToken_Get-12            3404556               336.8 ns/op
BenchmarkMToken_Get-12           3348738               315.0 ns/op
PASS
ok      github.com/appliedgo/refresh    3.054s
> go test -run NONE -bench=.
goos: darwin
goarch: arm64
pkg: github.com/appliedgo/refresh
BenchmarkToken_Get-12            4190270               302.7 ns/op
BenchmarkMToken_Get-12           4691931               295.7 ns/op
PASS
ok      github.com/appliedgo/refresh    4.108s
```

None of the two variants is considerably faster than the other. And by "considerably", I mean by at least an order of magnitude faster.



### Why do you use a timer-based token refresher?

A typical API client would call into the API until it gets a token expiration error. Then the client would refresh the token and continue where it left off.

In contrast to this, our web app receives session-less requests from an unknown number of site visitors. It could be one user or a thousand visitors requesting a transfer search or booking at the same time. The web app does not create a separate API connection for each of them. It maintains a single one.

In this scenario, there is no point in letting thousands of request run into the token timeout. They all make use of the same API connection, so the better approach is to refresh the token once for all visitors, before the timeout.

But keep in mind that this approach has also a downside. `refreshToken()` is a goroutine that enters an infinite loop. It may therefore outlive the context that it was created in, unless it gets stopped from the outside.

This is a potential goroutine leak. If an app needs a lot of separate API connections, it would create many `Token` variables, and each of them would spawn a neverending goroutine that continues to consume resources.

Remember:

> Always think about your goroutines' lifetime.

**As a rule of thumb, keep goroutines short-lived wherever possible. Do not let goroutines enter infinite loops without ensuring they can—and will—be stopped when their job is done.**

The current solution uses a cancelable context, but the app must ensure to eventually cancel the context. That's an additional burden for the `Token` consumer, unless the consumer already uses a cancelable context that the token refresher can be hooked into.


## Conclusion

We went from a simple channel-based future to a self-refreshing value. We compared a channel-based approach and an approach using mutexes, and we had a critical look at the token refresher that uses a timer instead of letting API calls run into the token timeout. I hope you had some fun time while reading this article and trying out the code. thank you

## Links

- [Futures in Go • Applied Go][futures]
- [The Singleflight package][singleflight]
- [Select statements (Go language specification)][select]
- [This code on GitHub][repository]
- [Go playground][playground]

[futures]: https://appliedgo.net/futures
[singleflight]: https://pkg.go.dev/golang.org/x/sync/singleflight
[select]: https://go.dev/ref/spec#Select_statements
[repository]: https://github.com/AppliedGo/refresh
[playground]: https://go.dev/play/p/jvUlqDwWyxq

**Happy coding!**

___

Changelog:

2023-10-25: Add context and an atomic value to address a potential goroutine leak and a race condition.

2023-10-26: Add a mutex version and benchmarks.

2023-10-28: Replace footnotes with Commonmark-compliant link references.

*/
