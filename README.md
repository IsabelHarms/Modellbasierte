# Modellbasierte

The parser is a recursive descent parser with 1 Subroutine (func in GO) for each syntactical element. The grammar allows for all decisions to be made upon the next token presented in the input stream. Upon Get() of a token not applicable to the current syntax element, the tokenizer accepts 1 Unget() in order to re-present that token to the next caller of Get().

Since the implemented simple  language IMP allows only integers and booleans, values are implemented by a structure with 1 int and 1 bool value plus a type identifyer to decide which one of those two is applicable. In addition to syntax checking, the parser does type checking, where multiple error messages are avoided by a type "Invalid" that is applicable to all operators. 
Error messages are handled and counted by the tokenizer, because it knows the position of the  last delivered token

Variables are entered into a table, mapping the name to their type and current value,  implemented by a GO "map" (using quick hashing). We use one map for each block nesting level, initialized upon the first declaration within that block, if any.
Upon end of a block, the scope of its variables ends by discarding the affiliated map.

Blocks are implemented as statement slices in GO, because they expand dynamically, comparable to "List" in other programming languages. 

In the tree constructed by the parser, we use a specific structure for each statement, but the same structure for all operators in expressions, in order to make the code shorter. 
The main program receives the root node from the parser and, provided no errors are found, hands it to the interpreter, which can then rely on compatible types in all statements.  During interpretation, variables are again entered into and removed from the variable table, where of course additionally assignments change their values.
