# syngo <!-- omit in toc -->

**syngo** provides special-purpose synchronisation types (for Go)

![Language](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![GitHub release](https://img.shields.io/github/v/release/synesissoftware/syngo.svg)](https://github.com/synesissoftware/syngo/releases/latest)
[![Last Commit](https://img.shields.io/github/last-commit/synesissoftware/syngo)](https://github.com/synesissoftware/syngo/commits/master)
[![Go](https://github.com/synesissoftware/syngo/actions/workflows/go.yml/badge.svg)](https://github.com/synesissoftware/syngo/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/synesissoftware/syngo.svg)](https://pkg.go.dev/github.com/synesissoftware/syngo)


## Table of Contents <!-- omit in toc -->

- [Introduction](#introduction)
- [Installation](#installation)
- [Components](#components)
  - [Counters](#counters)
  - [Latches](#latches)
  - [Version](#version)
- [Examples](#examples)
- [Project Information](#project-information)
  - [Where to get help](#where-to-get-help)
  - [Contribution guidelines](#contribution-guidelines)
  - [Dependencies](#dependencies)
    - [Development/Testing Dependencies](#developmenttesting-dependencies)
  - [Related projects](#related-projects)
  - [License](#license)


## Introduction

**syngo** provides special-purpose synchronisation types for **Go** — one-way latches and counters that may be operated safely from multiple concurrent goroutines.

Unlike the broad primitives in the standard library `sync` package (mutexes, condition variables, wait groups, and so on), the types in **syngo** are unidirectional — once a latch has flipped, it stays flipped — and are intended for simple, low-overhead signalling and counting patterns.


## Installation

```Go
import syngo "github.com/synesissoftware/syngo"
```

```Go
import syngo_sync "github.com/synesissoftware/syngo/sync"
```


## Components


### Counters

```Go
// in "github.com/synesissoftware/syngo/sync"

// A unidirectional counter that counts down from an initial value and may
// be operated safely by multiple concurrent goroutines.
type DownCounter struct { /* ... */ }

func NewDownCounter(initialValue int64) DownCounter
func (l *DownCounter) Step() (newCount int64)
func (l *DownCounter) Load() (count int64)
```

```Go
// in "github.com/synesissoftware/syngo/sync"

// A unidirectional counter that counts up from an initial value and may be
// operated safely by multiple concurrent goroutines.
type UpCounter struct { /* ... */ }

func NewUpCounter(initialValue int64) UpCounter
func (l *UpCounter) Step() (count int64)
func (l *UpCounter) Load() (count int64)
```


### Latches

```Go
// in "github.com/synesissoftware/syngo/sync"

// A one-way switch that may be operated safely by multiple concurrent
// goroutines.
type BoolLatch struct { /* ... */ }

func NewBoolLatch() BoolLatch
func (l *BoolLatch) Set() (flipped bool)
func (l *BoolLatch) Load() bool
```

```Go
// in "github.com/synesissoftware/syngo/sync"

// A unidirectional latch that counts down from an initial value to a lower
// threshold that may be operated safely by multiple concurrent goroutines.
type DownLatch struct { /* ... */ }

func NewDownLatch(initialValue, threshold int64) DownLatch
func (l *DownLatch) Step() (flipped, isLatched bool, newCount int64)
func (l *DownLatch) Load() (isLatched bool, count int64)
```

```Go
// in "github.com/synesissoftware/syngo/sync"

// A unidirectional latch that counts up from an initial value to a higher
// threshold that may be operated safely by multiple concurrent goroutines.
type UpLatch struct { /* ... */ }

func NewUpLatch(initialValue, threshold int64) UpLatch
func (l *UpLatch) Step() (flipped, isLatched bool, newCount int64)
func (l *UpLatch) Load() (isLatched bool, count int64)
```


### Version

```Go
// in "github.com/synesissoftware/syngo"

const (
  VersionMajor uint16 = /* ... */
  VersionMinor uint16 = /* ... */
  VersionPatch uint16 = /* ... */
  VersionAB    uint16 = /* ... */
)

func Version() uint64
func VersionString() string
```


## Examples

Examples are provided in the `examples` directory, along with a markdown description for each. A detailed list TOC of them is provided in [EXAMPLES.md](./EXAMPLES.md).


## Project Information


### Where to get help

[GitHub Page](https://github.com/synesissoftware/syngo "GitHub Page")


### Contribution guidelines

Defect reports, feature requests, and pull requests are welcome on https://github.com/synesissoftware/syngo.


### Dependencies

* [**ver2go**](https://github.com/synesissoftware/ver2go/);


#### Development/Testing Dependencies

* [**require**](https://github.com/stretchr/testify/);
* [**STEGoL**](https://github.com/synesissoftware/STEGoL/);


### Related projects

**syngo** is a **Go**-only library; there are no sibling implementations in other languages at this time.

Other Synesis **Go** libraries include:

* [**ANGoLS**](https://github.com/synesissoftware/ANGoLS/);
* [**CLiC4.Go**](https://github.com/synesissoftware/CLiC4.Go/);
* [**Diagnosticism.Go**](https://github.com/synesissoftware/Diagnosticism.Go/);
* [**STEGoL**](https://github.com/synesissoftware/STEGoL/);
* [**ver2go**](https://github.com/synesissoftware/ver2go/);


### License

**syngo** is released under the 3-clause BSD license. See [LICENSE](./LICENSE) for details.


<!-- ########################### end of file ########################### -->
