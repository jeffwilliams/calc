# TODO

Some things that would be nice to have:

  * Add built in support for standard unit conversions, for example kg to lbs and ounces
    t=kg_to_lbs(2.2); floor(t); lbs_to_oz( t-floor(t))

  * Adding an integer to an integer modifies the first integer instead of copying. For example:

    > i=5
    > i+1
    6
    > i
    6

    this is not really expected. Same problem occurs with lists:

    > l=[1]
    > l+[2]
    [3]
    > l
    [3]

  * Make a filter function (a la map, reduce)
  * Allow comments in calcrc
  * Cannot define a function that takes a list as parameter
    > def ipv4_to_hex(v) "Convert a list representing an IPv4 address to integer" unbytes(lrev(v))
    Error: last line:1:1 (0): rule "def statement": Parameter 1 is invalid: expected main.BigIntList but got *big.Float

  * Add License
  
