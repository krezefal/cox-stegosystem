# Cox Algorithm (Frequency steganography)

## Description

This algorithm is robust to many signal processing operations. Watermark message bits
are embedded into AC DCT image coefficients with the highest energy per 1 data unit 
(8x8 pixels). Detection of an embedded digital watermark in it is performed using the
original image.

There are 3 realisations of the Cox algorithm watermark embedding:

<p align="center">
    <img src="examples/readme_img/cox_formulas.png"  width="120"/>
</p>

where `α` is the weight coefficient (normally [0; 1]) and `s` is the bit to embed
{-1; 1}.

Current implementation is based on the **first formula**. This variant can be used
in case when the energy of the watermark is comparable to the energy of the modified 
coefficient. Otherwise, either the watermark will be non-robust, or the distortion is
too large. Therefore, it is possible to embed information with an insignificant range
of variation in values of the coefficients' energy.

## Implementation

The program is a CLI application that performs watermark embedding and extracting
to/from image-container. This program is a Go module.

To build the binary, execute

`go build -o bin/cox cmd/main.go`

Or just run it w/o building the binary:

`go run cmd/main.go`

Folder with examples contains a couple of files on which program can be tested "out
of box"

### Using

To run the program in embedding mode use next example:

`<cox> -src <path/to/src.bmp> -m <binary message> -tg <path/to/tg.bmp>`

To run the program in extracting mode:

`<cox> -src <path/to/src.bmp> -ext -tg <path/to/tg.bmp>`

where\
`<cox>` is the name of the binary for particular platform\
`-src` is the path to source file\
`-ext` means *extract*. It is the optional flag, specified if extraction procedure is
required\
`-m` is the optional flag after which follows the binary sequence (watermark), 
specified if embedding procedure is required\
`-tg` is the path to target file (which will be created in case of the embedding
procedure)

**For example:**\
`./cox --src examples/peppers.bmp --m 10101010 --tg examples/result.bmp`\
or\
`./cox --src examples/peppers.bmp --ext --tg examples/result.bmp`

Info about other optional flags is provided with `-h` option

### Principle of operation

*ADD INFO*

## Theory. Example of how it works

*ADD INFO*

## Research

*add my observations*

## Other

*add info about the force flag*