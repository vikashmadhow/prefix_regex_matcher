# Prefix regular expression matcher and lexer
A regular expression implementation that can do partial prefix matching whereby it is 
supplied a string to match one character at a time and responds if the supplied prefix 
is a partial match, complete match or a failed match. This will be used to develop a 
fast lexer.

## Regular expression engine
The regular expression engine currently supports the following syntax:

| Expression | Meaning                                                                                                            |
|------------|--------------------------------------------------------------------------------------------------------------------|
| *x**       | Zero or more of *x*.                                                                                               |
| *x*+       | One or more of *x*.                                                                                                |
| *x*?       | Zero or one of *x*.                                                                                                |
| *x*{m,n}   | k of *x* where *m <= k <= n*. If *m* is not provided it is set to 0. If *n* is not provided it is set to infinity. |
| *x*{m}     | Same *x*{m,m}                                                                                                      |
| *x* \| *y* | *x* or *y*.                                                                                                        |
| *(x)*      | Precedence given to *x*.                                                                                           |
| *[a-zAB]*  | Character set: *a* to *z*, *A* and *B*.                                                                            |
| *[^a-zAB]* | Inverse of character set: any character other than *a* to *z*, *A* and *B*.                                        |

