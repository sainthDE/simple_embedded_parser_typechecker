Original grammer (ambiguous):

```
N ::= 0 | 1 | 2 
B ::= true | false
V ::= N | B
E ::= V | (E) | E + E | E * E | E && E | E || E
```

Refactored grammar (unambiguous but with left recursion):

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

Final refactoring (unambiguous w/o left recursion):

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
