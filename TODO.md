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

  * Make a reduce function for lists that takes a function and a list as parameters
  * Make a map function
  
