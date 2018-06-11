# Calc

Calc is a simple interactive calculator for your terminal.

# How's it work?

Input an expression and hit enter to evaluate:

    > 20+5+4
    29

Basic operators are supported:

    > (3+4)^(9-10/5)
    823543

And literal numbers may be specified in decimal, hex, or binary. The display base may be set to one of those as well:

    > 10+0xa+0b1010
    30
    > set obase hex
    > 10+0xa+0b1010
    0x1e
    > set obase bin
    > 10+0xa+0b1010
    0b11110
    > set obase dec
    > 10+0xa+0b1010
    30

And operands may be decimals:

    > 5.1+6.7
    11.800000

If one operand is a decimal the type of the expression is automatically "upcast" to decimals:

    > 1+2+3+5.1
    11.100000

Some operators useful to programmers are supported:

    > (0b100&0b110)|1
    5

Standard operations are performed as arbitrary length numbers:

    > (2^800)^8
    3908159226643238733174614283614836731126768107046341227250667472076853556843013838172049588691174330707818243450118448302812729951247321576513166624369426519566175521146060765081679976758041720679300544142578899999682218145590731711215852375861916259684264866953944533853780206232109689860956552348006732061695789938930646213944790078445226547509325302609329069445117085739611168420511100727807029960219755304683343928522835681236580544749345772997687789272005780505845692810279598474814928801575139205787258048513369093710961312514682075899525393917486191722044010992536554254433346757797401453317441157668323475117378300300675380311411783720892291860208840935427421876878495172143644788379083826461955232814452670024106634029782444314857725946984340662589552137681187058855370840678104158373102913142935322248166923135606488577076176786571072087787275721151671449698445547824376581179228842823243076579850082699440910879690172701654935338585419237394528910219466400398711390146945174827484382033620349117027739850870285498619302796390659011030983899571129475531951900699941995728715779973392015185546940934704246331341905671381265137758178677705185913030002603348040480378452842334938835344896384986472605870643265219489006098190067654108321170943392355400613564249445645639510323755781809289920764767965983067048457794256504910798982324928635401543742404992470685295256712710005713066462569470457811357440928814052826871480405082643768534413838217558052879567404685419337070919454165543913926254797535062175403291840282437565784510890527884282609693788352728845075427184666736270904718599967377231129340602704541090595660344392945021559935525438319887090353617135488670209943492849139965846896740313626495887105269676175557097165018916853148260794391708438199220888781289029685829505315773902138853990717604142885719601697094720405328129745119056694810317474513400330333517233611193133379954007384384395034058799871873253376

Values may be stored in variables and referenced later:

    > usd_per_cad = 0.77
    > 60*usd_per_cad
    46.200000

Where'd my money go? The special variable `last` refers to the value of the last evaluation:

    > 67*8
    536
    > last
    536
    > 200+last
    736

Functions may also be defined:

    > def lbs_n_oz_to_kg(lbs,oz) "convert pounds and ounces to kg" (lbs+oz/oz_per_pound)/lbs_per_kg
    > lbs_per_kg=2.20462
    > oz_per_pound=16
    > lbs_n_oz_to_kg(160,6)
    72.574866

Help lists the defined functions:

    > help
    &(p1, p2): return p1 & p2 (bitwise and)
    *(p1, p2): return p1 * p2
    +(p1, p2): return p1 + p2
    -(p1, p2): return p1 - p2
    /(p1, p2): return p1 / p2
    ^(p1, p2): return p1 ^ p2
    abs(p1): absolute value. This function only has the precision of a float64.
    acos(p1): arccosine. This function only has the precision of a float64.
    acosh(p1): inverse hyperbolic cosine. This function only has the precision of a float64.
    asin(p1): arcsine. This function only has the precision of a float64.
    asinh(p1): inverse hyperbolic sine. This function only has the precision of a float64.
    atan(p1): arctangent. This function only has the precision of a float64.
    atanh(p1): inverse hyperbolic tangent. This function only has the precision of a float64.
    binom(p1, p2): binmomial coeffient of (p1, p2)
    bit(p1, p2): return the value of bit p2 in p1, counting from 0
    bytes(p1): return a list of each byte composing an integer
    cbrt(p1): cube root. This function only has the precision of a float64.
    ceil(p1): ceiling. This function only has the precision of a float64.
    choose(p1, p2): p1 choose p2. Same as binom
    cos(p1): cosine. This function only has the precision of a float64.
    cosh(p1): hyperbolic cosine. This function only has the precision of a float64.
    erf(p1): error function. This function only has the precision of a float64.
    erfc(p1): error function compliment. This function only has the precision of a float64.
    exp(p1): calculates e^p1, the base-e exponential of p1. This function only has the precision of a float64.
    exp10(p1): calculates 10^p1, the base-10 exponential of p1. This function only has the precision of a float64.
    exp2(p1): calculates 2^p1, the base-2 exponential of p1. This function only has the precision of a float64.
    floor(p1): floor. This function only has the precision of a float64.
    gamma(p1): gamma function. This function only has the precision of a float64.
    hex_to_ipv4(p1): Convert a hex value to an IPv4 address
    hypot(p1, p2): calculates sqrt(p1*p1 + p2*p2). This function only has the precision of a float64.
    j0(p1): order zero bessel function of the first kind. This function only has the precision of a float64.
    j1(p1): order one bessel function of the first kind. This function only has the precision of a float64.
    lbs_n_oz_to_kg(p1, p2): convert pounds and ounces to kg
    li(p1, p2): return element at index p2 in list p1
    llen(p1): return length of a list
    log(p1): natural logarithm. This function only has the precision of a float64.
    log10(p1): base-10 logarithm. This function only has the precision of a float64.
    log2(p1): base-2 logarithm. This function only has the precision of a float64.
    lrev(p1): return a copy of list p1 with elements in reverse order
    lrp(p1, p2): return a list consisting of p1 repeated p2 times
    map(p1, p2): return a new list which is the result of applying the function p2 to each element in p1
    neg(p1): return -p1 
    now(): return the number of milliseconds since epoch
    reduce(p1, p2, p3): apply a dyadic function p2 to each element in the list p1 and an accumulator (having initial value p3), returning the final value of the accumulator
    roll(p1, p2): roll p1 dice each having p2 sides and sum the outcomes
    sin(p1): sine. This function only has the precision of a float64.
    sinh(p1): hyperbolic sine. This function only has the precision of a float64.
    sqrt(p1): square root. This function only has the precision of a float64.
    tan(p1): tangent. This function only has the precision of a float64.
    tanh(p1): hyperbolic tangent. This function only has the precision of a float64.
    unbytes(p1): treat the list as a list of bytes and convert it to an integer
    y0(p1): order zero bessel function of the second kind. This function only has the precision of a float64.
    y1(p1): order one bessel function of the second kind. This function only has the precision of a float64.
    |(p1, p2): return p1 | p2 (bitwise or)
    ~(p1): return p1 | p2 (bitwise not)


Note that `lbs_n_oz_to_kg` is in there. There are a number of predefined functions. The ones that apply to decimal numbers only support double precision: 

    > exp(2^8)
    1511427665004103527714100498092829891603482697174374415092350456743517150826614334359230562343706299625849749504.000000
    > exp(2^800)
    0.000000

Basic list/vector support is included as well:

    > [2,3,4,5]+[1,2,3,4]
    [3, 5, 7, 9]
    > l=[1,2,3,4]
    > l
    [1, 2, 3, 4]
    > llen(l)
    4
    > li(l,1)
    2
    > lrev(l)
    [4, 3, 2, 1]
    > lrp(4,10)
    [4, 4, 4, 4, 4, 4, 4, 4, 4, 4]

Conversion to and from bytes is possible, as demonstrated in the following example. It converts a list representing an IPv4 address to a hex value as it would appear as a 32-bit word in memory on a big-endian machine:

    > set obase hex
    > unbytes([127,0,0,1])
    0x7f000001

And on a little endian:

    > set obase hex
    > unbytes(lrev([127,0,0,1]))
    0x100007f

In reverse:

    > bytes(0x7f000001)
    [127, 0, 0, 1]
    > lrev(bytes(0x100007f))
    [127, 0, 0, 1]

One might define a convenience function for the little endian conversion above, if one often works with encoded IP addresses in gdb:

    > def hex_to_ipv4(v) lrev(bytes(v))
    > hex_to_ipv4(0x100007f)
    [127, 0, 0, 1]

Commonly used user-defined functions (such as `hex_to_ipv4`) and variables may be defined in `~/.calcrc`, which is loaded on startup. 

Functions, while not fully first-class, can be assigned to variables and passed to functions. This is useful when applying a function
to a list of values using `map`:

    > map([25.0,9.0,81.0], sqrt)
    [5, 3, 9]

The basic arithmetic operators are internally defined as functions, and may be called as functions:

    > 1+2
    3
    > +(1,2)
    3

This comes in handy when reducing a list:

    > reduce([1,2,3], +, 0)
    6

Note that the unary negation function is `neg` not `-`. The `-` is used for subtraction.

Calc supports readline-like line editing: UP moves to the previous expression, arrow keys, home, end, CTRL-A, CTRL-E, CTRL-U, CTRL-R, and CTRL-K all behave as expected. On a blank line the TAB key auto-completes against defined functions, variables, and keywords.

Use CTRL-C to exit.

# Install

Grab the latest archived binary from the [releases](https://github.com/jeffwilliams/calc/releases) page and unpack it.

# Building

Calc requires the Go programming language version 1.8 to build (Needs big.Float to implement fmt.Scanner). First, install Go, then:

## On Linux

    make 

## On other systems

    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
    go get github.com/mna/pigeon github.com/cheekybits/genny
    go generate
    dep ensure
    go install
    go test

# Acknowledgements

A special thanks goes out to Termux, the Android terminal emulator. A large part of calc was written in Termux using Vim on a Blackberry Priv.
