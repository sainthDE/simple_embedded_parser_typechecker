# ReadMe

A simple parser typechecker embedded in Go to fulfill <https://sulzmann.github.io/ProgrammingParadigms/projectExp.html#(2)>

## Grammar

### Original grammer (ambiguous):

```
N ::= 0 | 1 | 2 
B ::= true | false
V ::= N | B
E ::= V | (E) | E + E | E * E | E && E | E || E
```

### Refactored grammar (unambiguous but with left recursion):

```
N  ::= 0 | 1 | 2 
B  ::= true | false
V  ::= N | B
F  ::= V | (EP)
EP ::= EP + EO | EO
EO ::= EO || EM | EM
EM ::= EM * EA | EA
EA ::= EA && F | F
```

### Final refactoring (unambiguous w/o left recursion):

```
N   ::= 0 | 1 | 2
B   ::= true | false
V   ::= N | B
F   ::= V | (EP)
EP  ::= EO EP2
EP2 ::= + EO EP2 |
EO  ::= EM EO2
EO2 ::= || EM EO2 |
EM  ::= EA EM2
EM2 ::= * EA EM2 |
EA  ::= F EA2
EA2 ::= && F EA2 |
```

## Some examples

The output of the examples consists of three lines:

1. The input value is behind Original
1. The parsed value is behind Parsed
1. The inferred type

### Positive examples

```
Original: 1
Parsed  : 1
 Int
```

```
Original: 1+0
Parsed  : (1+0)
 Int
```

```
Original: 1 * 2 
Parsed  : (1*2)
 Int
```

```
Original:  (1) 
Parsed  : 1
 Int
```

```
Original:  (1 * (2)) 
Parsed  : (1*2)
 Int
```

```
Original:  (1 + 2) * 0 
Parsed  : ((1+2)*0)
 Int
```

```
Original: true || false
Parsed  : (true||false)
 Bool
```

### Negative examples

```
Original: 1+
Parsed  : Syntax Error
 Illtyped
```

```
Original: + 1
Parsed  : Syntax Error
 Illtyped
```

```
Original: (((1))
Parsed  : Syntax Error
 Illtyped
```

```
Original: tru
Parsed  : Syntax Error
 Illtyped
```

```
Original: fal
Parsed  : Syntax Error
 Illtyped
```

```
Original: true | false
Parsed  : Syntax Error
 Illtyped
```

```
Original:  1 + true
Parsed  : (1+true)
 Illtyped
```

```
Original: 2 * false
Parsed  : (2*false)
 Illtyped
```

```
Original:  1 || true
Parsed  : (1||true)
 Illtyped
```

```
Original:  1 && false
Parsed  : (1&&false)
 Illtyped
```