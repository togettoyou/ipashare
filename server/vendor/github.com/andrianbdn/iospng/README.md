# iOS PNG File Normalizer 

This package reverts optimizations that are done by Xcode for PNG files when packaging iOS apps:

- Removes CgBI chunks
- Fixes compressed IDAT chunks  
- removes alpha pre-multiply

The package does similar things like ipin.py or <code>xcrun -sdk iphoneos pngcrush -revert-iphone-optimizations</code>

## Installation


The import path for the package is *github.com/andrianbdn/iospng*.

To install it, run:

    go get github.com/andrianbdn/iospng


## Usage

#### func PngRevertOptimizationWithSize

```go
func PngRevertOptimizationWithSize(reader io.Reader, writer io.Writer) (int, int, error)
```

This function actually does everything: reads PNG from reader and in case it is iOS-optimized, reverts optimization. 
Function does not change data if PNG does not have CgBI chunk.

First two return parameters (ints) are PNG width and height.

#### func PngRevertOptimization

```go
func PngRevertOptimization(reader io.Reader, writer io.Writer) error
```

Old-style function that does not return size also available.

 
## See also 

- [CgBI file format](http://iphonedevwiki.net/index.php/CgBI_file_format)
