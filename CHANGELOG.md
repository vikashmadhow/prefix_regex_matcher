# Changelog

## [0.2.1] - 2024-12-30
- Filtering, mapping and flat-mapping can now work on both the pull and push versions of iter.Seq
  and iter.Seq2.
- Lexer `Lex` methods for returning a push version of iter.Seq2, which is simple to iterate over.
  However, these version does not allow for token pushback.

## [0.2.0] - 2024-12-30
- A set of utility functions for filtering, mapping and flat-mapping over iter.Seq and iter.Seq2.
- Lexer now reads its input from an `io.Reader` which is more memory efficient.
- Flatmap functions can be attached to a lexer to modulate its output, arbitrarily changing the
  token stream. This can be used to ignore certain tokens, modify tokens or insert new tokens
  at arbitrary points in the token stream.
- Simple context-free-grammar definition and predictive parser (untested).

## [0.1.3] - 2024-12-02
- Lexer will not provide the token(s) which failed and the expected next character(s) for each.
- Error information is now provided if the last token in the input stream has an error,

## [0.1.2] - 2024-11-30
- Improved lexer matching, error and position tracking.

## [0.1.1] - 2024-11-27
- Line and column where each token matched reported in `Token` by `Lexer`. 

## [0.1.0] - 2024-11-25
- Numbered capturing groups with group 0 matching whole string and each opening
  parenthesis `(` starting a new capturing group. Capturing groups can be nested.

## [0.0.0-20241124] - 2024-11-24
- Character classes shortcuts (`\d`, `\w`, etc.)
- Matching any character with `.`.
- More regular expression tests.

## [0.0.0-20241123] - 2024-11-23
- Lexer and regex tests.
- `match` method in `CompiledRegex` match regular expression to string exactly.
- Range repetition in regular expression (`re{m,n}`).

## [0.0.0-20241122] - 2024-11-22
- Working DFA implementation for regular expression matching.
- Base incremental lexer working and tested.
