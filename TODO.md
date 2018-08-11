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

- make sure different programs leave the stack empty. 
  - builtin stored as lambda in variable, and called. 


- congirm 

Lambdas:
  - Need to make builtins work as lambdas. Try it now, and it will print an error woth the place needed to be changed
  - If an identifier resolves to a builtin, make a new function for it in the code. Name it like @builtin-X where X 
    is the builtin name. Before making it check if it exists. Then push the address of this function. The function just
    calls the builtin.

  - With lambdas, since they get stored in a variable, if the variable is set to a new value the lambda's environment is 
    leaked (since a new one is needed each time). So need basic mark-and-sweep garbage collector for labda envs: mark
    all lambda envs reachable from vars, and any not reachable sweep.

Closures:
  - make a variable in the compiler for the current closure environment. When a real function defn is being compiled
    initialize it.
  - In identifier handling, if the identifier resolves to a parameter of the first non-lambda function, then 
    - pick an environment index for it, starting from 0
    - add this identifier name to the closure env. with it's index
    - generate the code that loads the index from the environment of the lambda
      pushparm 0  // first parameter is the closure environment
      push <index>
      iadd // add the ints and push the result: the final address of the variable in the environ
      loads // load the value at the address on the top of the stack and put it on the stack
  - In the function where the closure is defined, read the closure environment builder var to see what
    variables are there with which indexes and names, and make a variable slot for each in order. 
    - How will the function dynamically allocate those values? Mostly this acts like a stack really. Maybe
      we can store it on the stack, and access from there. 
    - Alternately have an instruction that adds a new slot to the end of the data segment and returns the index.
      - alloc/free

  code in function where closure is defined:

  code in closure function
  pushparm 0  // first parameter is the closure environment
  


Closures: In the outer function that is defining the inner function, create an environment (use some new memory slots specifically
for this inner function) and when generating the instructions for the inner function have it use the environment slots for the 
variables. Basically pass a pointer to the environment to the inner function as a new parameter, and it uses it to access the
variables.

- For compile errors: store line number and col in the AST nodes. Can be gotten in the PEG actions using c.pos.line, c.pos.col.
 


--

dataSegment []interface{}

(exists only in AST:)
var:
  name
  type
  where (data segment or stack relative to bp)
  index (into data OR relative to BP when function parm or local)
  ---->  points to data that is one of:
    big.Int
    big.Float
    Func
    
- In the VM store symbol tables (name->index) for vars (points into Data) and funcs (points into instructions). This is needed
  for the repl so that when we link a new line and execute, if it modifies an existing variable it is able to do so. Same for 
  function redefinition.

Compile:
  input: code
  output: 
    - executable code
    - defined function code (shared)
      - with symbol table
    - symbol table for detected variables


Link together:
  - executable code
  - defined function code

At Run:
  - allocate a data segment big enough for the symbol table for variables
