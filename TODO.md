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

Lambdas:
  - Need to make builtins work as lambdas. Try it now, and it will print an error woth the place needed to be changed

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
