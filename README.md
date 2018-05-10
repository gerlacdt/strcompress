## Description

Decompresses a given compressed string.
You can find a detailed problem description here:
https://techdevguide.withgoogle.com/paths/advanced/compress-decompression/#!

``` bash
# Examples:
# 3[a]          ->  aaa
# 3[abc]4[ab]c  ->  abcabcabcababababc
# 2[3[a]b]      ->  aaabaaab
```

The implementation uses a recursive-descent parser with the following
grammar:

```bash
# <Exp> ::= <Number> '[' <Exp> ']' <Letter> | <Letter>
# <Number> ::= [0..9]+
# <Letter> ::= [a-z]*
```
