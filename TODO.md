# TODO

Some things that would be nice to have:

  * Add built in support for standard unit conversions, for example kg to lbs and ounces
    t=kg_to_lbs(2.2); floor(t); lbs_to_oz( t-floor(t))

  * Allow comments in calcrc
  * Add License
  * Add mod
  * add shifting

## VM/AST

How to allow a function to be assigned to a variable? And closures.
  - Currently a Func object is stored in the var.
  - Instead we would need to redefine Func to point to the global offset where the function code begins
  - How do we do this for a function defined immediately and assigned to a variable?
  - Put the function at the end of the global compiled code.
  - What if the function already exists?
  - Redefine it, and cut out the old code from the global compiled code. IF it's in the middle, we need to repack the code.
  - Then store the Func object in the variable
  - When a call is being compiled, we should be able to tell if it's a global function or a variable.
  - If a variable, generate the code so that it looks up the variable contents at runtime, and calls into the offset for the function


Must allow closures. If each line is compiled and evaluated, a variable may be set to a function 

Closures: when the function is defined it will need to have the scope (variables) stored in a data slot.
  Can make the variable a function struct that has pointer to code, and pointer to scope variables.
  Alternately just allocate new variables for the function's scope vars, and when compiling the code refer to them
  directly.

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
