# Behavior-driven development (BDD) in Go

## Introduction

![Gophers](https://raw.githubusercontent.com/hedhyw/gherkingen/main/assets/gophers-bdd-by-hedhyw.svg)

The article will introduce a generation of BDD tests using [gherkingen](https://github.com/hedhyw/gherkingen#gherkingen) for Golang.

`gherkingen` is a BDD boilerplate generator and framework. It accepts a *.feature [Cucumber/Gherkin](https://cucumber.io/docs/gherkin/reference/) file and generates a test boilerplate. All that remains is to change the tests a little. The generator provides a standart Golang way of testing stuffed with a BDD approach.

The generator is very customizable, it is possible to customize an output for any golang testing framework or even for another language, because it uses [Golang templates](https://github.com/hedhyw/gherkingen#creating-templates).

## Why BDD?

The main goal of using BDD testing is for improving collaboration between all employees (developers, product managers, business analysts and testers) by a simple converstaion of scenarios. Instead of using developer's way of describing tests, it gives a possibility to describe tests using human-readable language. BDD provides high quality tests that can be then implemented using test driven development (TDD). Also, it brings a good documentation of a feature and a good knowledge of a buisness behavior.

Let's compare, a developer's way:
```c++
assert sum(2, 3) == 5
```
and a BDD's way:
```feature
Feature: The Universe Summator
    Scenario: Adding two numbers
        Given John has the number 2 on the left hand
        And John has the number 3 on the right hand
        When he puts his hands into the summator
        Then the summator outputs 5
```


## Task

Pictures can have size up to 3840×2160! And we don't want a user to wait until it will be loaded. The user can see loading even if he uses the fastest connection. So let's create an HTTP-server that resizes pictures. API will be simple, it will be just one endpoint that accepts desirable size and a link to the picture. The server will fetch an image, resize it, and return to the user.

### Scenarios

Let's define user scenarios that we discussed with a team. I prepared those scenarios with a help of my friends, so this is much closer to a real example.

```feature
// internal/features/image_resizing.feature
Feature: Picture resizing
    Scenario: Duyên wants to receive a resized picture
        """ Offtopic:
                Duyên is a vietnamese girl firstname. It means destiry.
                It is interesting that "D" reads as english "Z",
                and "Đ" reads as english "D".
        """
        Given Duyên selects the size '<width>x<height>'
        And a link to existen picture of a bigger size of type 'image/jpeg'
        When Duyên calls an endpoint
        Then she receives an image of a content type 'image/jpeg'
        And the size of the image is '<width>x<height>'
        Examples:
        | <width> | <height> |
        | 256     | 128      |
        | 128     | 128      |
        | 1       | 1        |

    Scenario: Duyên provides an invalid size
        Given Duyên selects the size <size>
        And a link to an existen picture of a bigger size
        When Duyên calls an endpoint
        Then she receives an error
        Examples:
        | <size> |
        | 0x0    |
        | 0x1    |
        | ax10   |
        | 10,10  |
        | -1x10  |
    
    Scenario: Duyên provides an invalid link
        Given Duyên selects the size 256x256
        And a link to an unexistent picture
        When Duyên calls an endpoint
        Then she receives an error
```

### Prepare service boilerplate

The server will be very simple in this step, it will have only one endpoint that always responds with InternalServerError (500) status.

```go
// internal/server/server.go
func New() *Server {
    mux := http.NewServeMux()
    s := &Server{
        handler: mux,
    }

    mux.HandleFunc("/api/image.jpg", s.handleJPEGImage)

    return s
}

func (s *Server) handleJPEGImage(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
	_ = q.Get("size")
	_ = q.Get("url")

    http.Error(w, "unimplemented", http.StatusInternalServerError)
}
```

### Install generator

Let's install gherkingen first. For simple Go installation, run:
```
go install github.com/hedhyw/gherkingen/cmd/gherkingen@latest
```

Other ways are described in [the documentation](https://github.com/hedhyw/gherkingen#install).

### Generating tests

It is a time for generating BDD-tests

I prefere the following structure:
- internal/features/feature_name.feature
- internal/features/feature_name_test.go

There is a folder "features" that contains a set of feature-files and their test implementations.

Let's generate a test
```sh
cd internal/features
gherkingen -package features_test image_resizing.feature > image_resizing_test.go
```

Alternative way of a generation with docker:
```sh
docker run --rm -it --read-only --network none \
	--volume $PWD:/host/ \
	hedhyw/gherkingen:latest \
	-- /host/internal/features/image_resizing.feature \
    -package features_test \
    > internal/features/image_resizing_test.go
```

The generator always prints to the stdout, because it safer to avoid overriding of existing file. So in real generations, prefere to copy a code from the stdout by hand.

The package name is specified for convinience, by default it is "example_test".

It will generate a test:

```go
// internal/features/image_resizing_test.go
func TestPictureResizing(t *testing.T) {
	f := bdd.NewFeature(t, "Picture resizing")

	f.Scenario("Duyên wants to receive a resized picture", func(t *testing.T, f *bdd.Feature) {
		type testCase struct {
			Size string `field:"<size>"`
		}

		testCases := map[string]testCase{
			"256x128": {"256x128"},
			"128x128": {"128x128"},
			"1x1":     {"1x1"},
		}

		f.TestCases(testCases, func(t *testing.T, f *bdd.Feature, tc testCase) {
			f.Given("Duyên selects the size <size>", func() {

			})
			f.And("a link to existen picture of a bigger size", func() {

			})
			f.When("Duyên calls an endpoint", func() {

			})
			f.Then("she receives and an image of a content type 'image/jpeg'", func() {

			})
			f.And("the size of the image is <size>", func() {

			})
		})
	})
    /*
        Other part of code is hidden for convinience.
    */
}
```

This test should pass by default.

### Preparation

The next out step is to create an implementation. We need to call a method of running server. For this we are going to create some helpers:

```go
// internal/features/features.go
type testHelper struct {
    addr string
}

func newTestHelper(tb testing.TB) *testHelper {
    return &testHelper{addr: runTestServer(tb)}
}

func (th *testHelper) CallAPIGet(tb testing.TB, path string, query url.Values) (resp *http.Response) {
    tb.Helper()
    /* Implementation is hidden. */
    return resp
}
```

And start implementing our test.

### Test implementation

I advise implement small steps, no more than 5 lines. In this case, the tests will still remain readable. For this, separate all large sections into separate helper functions.

1. First we will define a context that will be shared between tests.

```go
// internal/features/image_resizing_test.go
f.TestCases(testCases, func(t *testing.T, f *bdd.Feature, tc testCase) {
    th := newTestHelper(t) // It was created before.
    q := make(url.Values) // We will put here our input parameters.
    var resp *http.Response // Result.

    /* Part of code is hidden. */
})
```
2. Next, let's define preconditions:
```go
f.Given("Duyên selects the size <size>", func() {
    q.Set("size", tc.Size) // We use current test case.
})
f.And("a link to an existen picture of a bigger size", func() {
    // newPhotoServer is a helper server that serves a big JPEG image.
    photoURL := newPhotoServer(t)
    q.Set("url", photoURL)
})
```
3. And execution part.
```go
f.When("Duyên calls an endpoint", func() {
    resp = th.CallAPIGet(t, pathImageResize, q)
})
f.Then("she receives an error", func() {
    require.Equal(t, http.StatusBadRequest, resp.StatusCode)
})
```
4. That is all. Now we just run the test `go test ./...`. It should fail because our server is not implemented.

### Final part

The rest of our work is the implementation of our feature. I will not be described here, because the goal of the article was about generation.

## Other tools

|  Tool                                                  | Description                       | Generator? | Fraemwork? |
|--------------------------------------------------------|-----------------------------------|------------|------------|
| [godog](https://github.com/cucumber/godog)             | official BDD generator for golang | ✅         | ✅         |
| [ginkgo](https://github.com/onsi/ginkgo)               | BDD fraemwork, not a generator    | ❌         | ✅         |
| [goconvey](https://github.com/smartystreets/goconvey/) | BDD-style framework with web UI and live reload |  ❌         | ✅         |

What a difference between godog and gherkingen? gherkingen uses go-templates, so it is fully customizable, and if it is desired, it can be easily extended to any language. It also has a default simple BDD-style implementation for a testing that is compitable with standart go tests.

### Conclusion

Thank you for your attention! I hope you will rate and use BDD testing approach in your Golang projects. The full example can be found in the repository [https://github.com/hedhyw/bdd-resizer-example](https://github.com/hedhyw/bdd-resizer-example).