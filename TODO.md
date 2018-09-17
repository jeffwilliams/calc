# TODO

Some things that would be nice to have:

  * Add built in support for standard unit conversions, for example kg to lbs and ounces
    t=kg_to_lbs(2.2); floor(t); lbs_to_oz( t-floor(t))

  * Allow comments in calcrc
  * Add License
  * Add mod
  * add shifting

## VM/AST

http://hokstad.com/how-to-implement-closures


Closures:
 
  - Closure tables are allocated dynamically. When stored to a variable, and then the variable is pointed to something else,
    the allocation will leak. Must implement garbage collection (mark and sweep)

