grammar jsstrip;

sexpr:
   list
   |
   atom
   ;

list:
    Lparen Rparen
    | Lparen members Rparen
    ;

members:
    (sexpr)+
    ;

atom:
    Ident
    | Num
    | QuotedString
    ;

Whitespace: ( ' ' | '\r' '\n' | '\n' | '\t' ) -> skip;

// need to put these simple ones ahead of the complex ones
Lparen: '(';
Rparen: ')';
Quote: '"';

// LineComment: ';;' ~('\r' | '\n')*;
Num: ( '0' .. '9')+;
fragment IdentFirst: ('a' .. 'z' | 'A' .. 'Z' | '.' | '$' | '_' | ';');
fragment IdentAfter: ('a' .. 'z' | 'A' .. 'Z' | '.' | '$' | '_' | ';' | '0'..'9');
Ident: IdentFirst IdentAfter* ;

fragment HexByteValue: '\\' ( '0' .. '9' | 'a' .. 'f') ( '0' .. '9' | 'a' .. 'f');
QuotedString: '"' ( HexByteValue | ~('"') )* '"';
