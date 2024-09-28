// ## Imports and globals
package main

import "fmt"

// ### The "inner ring"

// A `Poem` contains some poetry and an abstract storage reference.
type Poem struct {
	content []byte
	storage PoemStorage
}

// `PoemStorage` is just an interface that defines the behavior of a poem storage.
// This is all that `Poem` knows (and needs to know) about storing and retrieving poems.
// Nothing from the "outer ring" appears here.
type PoemStorage interface {
	Type() string        // Return a string describing the storage type.
	Load(string) []byte  // Load a poem by name.
	Save(string, []byte) // Save a poem by name.
}

// `NewPoem` constructs a `Poem` object. We use this constructor to inject an object
// that satisfies the `PoemStorage` interface.
func NewPoem(ps PoemStorage) *Poem {
	return &Poem{
		content: []byte("I am a poem from a " + ps.Type() + "."),
		storage: ps,
	}
}

// `Save` simply calls `Save` on the interface type. The `Poem` object neither knows
// nor cares about which actual storage object receives this method call.
func (p *Poem) Save(name string) {
	p.storage.Save(name, p.content)
}

// `Load` also invokes the injected storage object without knowing it.
func (p *Poem) Load(name string) {
	p.content = p.storage.Load(name)
}

// `String` makes Poem a Stringer, allowing us to drop it anywhere a string would be
// expected.
func (p *Poem) String() string {
	return string(p.content)
}

// ### The "outer ring"

// #### The notebook

// A `Notebook` is the classic storage device of a poet.
type Notebook struct {
	poems map[string][]byte
}

func NewNotebook() *Notebook {
	return &Notebook{
		poems: map[string][]byte{},
	}
}

// After adding `Save` and `Load`, `Notebook` implicitly satisfies `PoemStorage`.
func (n *Notebook) Save(name string, contents []byte) {
	n.poems[name] = contents
}

func (n *Notebook) Load(name string) []byte {
	return n.poems[name]
}

// `Type` returns an informal description of the storage type.
func (n *Notebook) Type() string {
	return "Notebook"
}

// A `Napkin` is the emergency storage device of a poet.
// It can store only one poem.
type Napkin struct {
	poem []byte
}

func NewNapkin() *Napkin {
	return &Napkin{
		poem: []byte{},
	}
}

func (n *Napkin) Save(name string, contents []byte) {
	n.poem = contents
}

func (n *Napkin) Load(name string) []byte {
	return n.poem
}

func (n *Napkin) Type() string {
	return "Napkin"
}

// ### Wiring everything up

// Create and connect objects, then save and load a few poems from different storage objects.
func main() {
	notebook := NewNotebook()
	napkin := NewNapkin()

	// First, write a poem into a notebook.
	// `NewPoem()` injects the dependency.
	poem := NewPoem(notebook)
	poem.Save("My first poem")

	// Create a new poem object to prove that the notebook storage works.
	poem = NewPoem(notebook)
	poem.Load("My first poem")
	fmt.Println(poem)

	// Now we do the same with a napkin as storage.
	poem = NewPoem(napkin)
	// Note the poem still just uses `Save` and `Load`. "Notebook? Napkin? I don't care."
	poem.Save("My second poem")
	poem = NewPoem(napkin)
	poem.Load("My second poem")
	fmt.Println(poem)
}

/* As usual, you can `go get` the code from GitHub. Don't forget to use -d if you do not wish to have the exectuable in your $GOPATH/bin directory.

    go get -d github.com/appliedgo/di
	cd $GOPATH/src/github.com/appliedgo/di
	./di

## Conclusion

Outside the world of poetry, dependency injection is a useful tool for decoupling logical entities, especially in multi-layered architectures as we have seen above.

Besides its benefits for layered architectures, dependency injection can also help with testing. Instead of reading a poem from a real notebook, a test can read from a notebook mockup that either is easier to set up, or delivers consistent test data, or both.


## Further reading

I definitely recommend reading the aforementioned [article about the Clean Architecture](http://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html) by [Robert C. Martin, a.k.a. "Uncle Bob"](https://de.wikipedia.org/wiki/Robert_Cecil_Martin).

The excellent article [Applying The Clean Architecture to Go applications](http://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/) is a deep dive into implementing DI in Go that builds upon all four layers of the Clean Architecture. This is a great opportunity to see how entities, use cases, interfaces, and frameworks (speaking in Clean Architecture lingo) are utilized to build a (toy) shop system.

Dependency Injection can be seen as one specific form of *loose coupling*, a term referring to interconnecting components without making them too dependent on each other. Another option for loose coupling in Go (besides interfaces) is to use *higher-order functions*. I found a quick and easy intro to this topic in the blog article [Loose Coupling in Go lang](https://blog.8thlight.com/javier-saldana/2015/02/06/loose-coupling-in-go-lang.html).
*/

// what if our poet decides to write a poem on a napkin? Or on 4x6 index cards?
// The document layer would have to be modified and recompiled!
//  We have created an unwanted dependency on a particular storage type.

// How can we remove that dependency?
// Abstraction to the rescue
//
// As a first step, we can replace the storage service by an abstraction of that service.
// Using Go's interface type, this becomes really easy.
// The interface describes only a behavior,
// and our Poem object can call the interface functions without worring about the object that implements this interface.
// Now we can define the Poem struct without any dependency on the storage layer:
// Remember, PoemStorage is just an interface but we can assign any type to storage that satisfies this interface.

// Adding dependency injection

// Right now the Poem only talks to an empty abstraction.
//  we need a way to connect a real storage object to the Poem.
// In other words, we need to inject a dependency on a PoemStorage object into the Poem layer.

// The interface/constructor pattern is not the only approach to implementing dependency injection. Still, it is a quite appealing one because it is clear and concise and builds upon just a few basic language constructs.

// What is Dependency Injection?

// Dependency Injection is a design pattern, that helps you to decouple the external logic of your implementation.
// It’s common an implementation needs an external API, or a database, etc.
// It isn’t responsibility of the implementation to know these things, it should receive its dependencies and use them as it needs.

// An important point of injecting dependencies is to avoid injecting implementations (structs), you should inject abstractions (interfaces).
//  It allows you to switch easily the implementation of some dependency and, you could change the real implementation for a mock implementation.

// Kinds of Dependency Injection

// There are some kinds of dependency injection, and each one has its own use case. Here we'll cover three of them: Constructor, Property and, Method (or Setter).

// The most common kind is the Constructor Injection. It allows you to make your implementation immutable, nothing can change the dependencies (if your properties are private). Also, it requires all dependencies to be ready to create something. If they aren’t, it usually will generate an error.
