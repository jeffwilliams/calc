# TODO

Some things that would be nice to have:

  * Add support for decoding IPv4 addresses from hex.
  * Add built in support for standard unit conversions, for example kg to lbs and ounces
    t=kg_to_lbs(2.2); floor(t); lbs_to_oz( t-floor(t))
  * Improve list support:
    - Create a function to make a list containing a repeated element
      - lrp(3,2) makes a list of size 2 containing 3,3
      - Can be used to simulate multiplying a list by a scalar:
        [1,2,3]*lrp(3,3)
    - add functions to extract list items
    - add function to get cardinality of a list

