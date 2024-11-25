# Prefix regular expression matcher and lexer
A regular expression implementation that can do partial prefix matching whereby it is 
supplied a string to match one character at a time and responds if the supplied prefix 
is a partial match, complete match or a failed match. 

This is then used to build a fast lexer.

## Regular expression engine
The regular expression engine currently supports the following patterns:

| Expression  | Meaning                                                                                                                                  |
|-------------|------------------------------------------------------------------------------------------------------------------------------------------|
| `x*`        | Zero or more of `x`.                                                                                                                     |
| `x+`        | One or more of `x`.                                                                                                                      |
| `x?`        | Zero or one of `x`.                                                                                                                      |
| `x{m,n}`    | `k` of `x` where `m <= k <= n`. If `m` is not provided it is set to 0. If `n` is not provided it is set to infinity.                     |
| `x{m}`      | Same as `x`{m,m}                                                                                                                         |
| `x \| y`    | `x` or `y`.                                                                                                                              |
| `(x)`       | `x` as a numbered capturing group, starting from 1. Group 0 is reserved for the whole expression. Precedence is also overridden by `()`. |

### Character and character classes
| Expression | Meaning                                                                            |
|------------|------------------------------------------------------------------------------------|
| `[a-zAB]`  | Character set: `a` to `z`, `A` and `B`.                                            |
| `[^a-zAB]` | Inverse of character set: any character other than `a` to `z`, `A` and `B`.        |
| `.`        | Matches any character.                                                             |
| `\d`       | Digits characters `[0-9]`.                                                         |
| `\D`       | Not digits characters `[^0-9]`.                                                    |
| `\s`       | Whitespace characters `[ \t\n\f\r]`.                                               |
| `\S`       | Not whitespace characters `[^ \t\n\f\r]`.                                          |
| `\w`       | Word characters `[0-9a-zA-Z_]`.                                                    |
| `\W`       | Not word characters `[^0-9a-zA-Z_]`.                                               |

