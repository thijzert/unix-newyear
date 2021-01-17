# UNIX New Year

This is a website that lets you celebrate UNIX New Year.

## Background

As you may know, UNIX date/time is measured as the number of seconds since the 1st of January 1970, 00:00:00 UTC. If you express that number in hexadecimal notation, once every half-a-year-and-a-bit™ (16,777,216 seconds) the number will end in six zeros. These occasions I like to celebrate.

For example, UNIX time 0x56000000 was the 21st of September 2015, 13:02:56 UTC.

Here's an example of a UNIX millennium occurring:
![Happy 0x60, everyone.](https://github.com/thijzert/unix-newyear/blob/master/.readme/moneyshot.gif?raw=true)

## Building

### Development Build

Accesses assets from disk directly:

```bash
go build -tags=dev
```

### Production Build

All assets are statically embedded in the binary, so it can run standalone in any folder:

```bash
go generate
go build
```

## License

This repository is available under the terms of the [BSD 3-clause license](https://opensource.org/licenses/BSD-3-Clause).
It is based on the [gopherpen](https://github.com/gopherjs/gopherpen) project template, which is licensed under the [MIT License](https://opensource.org/licenses/MIT). It also includes jQuery, which is also licensed under MIT.

### Credits
The firework animations are based upon Roy Whittle's fireworks at javascript-fx.com.
