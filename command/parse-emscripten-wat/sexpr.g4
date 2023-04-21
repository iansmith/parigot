grammar sexpr;

sexpr
   returns [[]*Item item_]
   : item* EOF
   ;

item
   returns [*Item item_]
   : atom
   | list_
   | LPAREN item DOT item RPAREN { panic("was not expecting dotted pair")}
   ;

list_
   returns [[]*Item list ]
   : LPAREN item* RPAREN
   ;

atom
   returns [*Atom atom_]
   : STRING
   | SYMBOL
   | NUMBER
   | DOT
   | COMMENT_NUM
   ;

STRING
   : '"' ('\\' . | ~ ('\\' | '"'))* '"'
   ;

WHITESPACE
   : (' ' | '\n' | '\t' | '\r')+ -> skip
   ;

NUMBER
   : ('+' | '-')? (DIGIT)+ ('.' (DIGIT)+)?
   ;

SYMBOL
   : SYMBOL_START (SYMBOL_START | DIGIT)*
   ;

COMMENT_NUM
   : ';' DIGIT+ ';'
   ;


LPAREN
   : '('
   ;

RPAREN
   : ')'
   ;

DOT
   : '.'
   ;

fragment SYMBOL_START
   : ('a' .. 'z')
   | ('A' .. 'Z')
   | '+'
   | '-'
   | '*'
   | '/'
   | '.'
   ;

fragment DIGIT
   : ('0' .. '9')
   ;
