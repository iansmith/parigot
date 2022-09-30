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
    Ident Offset? Align?
    | Num
    | QuotedString
    | BlockAnnotation
    | ConstAnnotation
    | HexPointer
    | ConstValue
    ;

Whitespace: ( ' ' | '\r' '\n' | '\n' | '\t' ) -> skip;
Comment: ';;' ~( '\n' | '\r')* ('\r' '\n'| '\n') -> skip;

// need to put these simple ones ahead of the complex ones
Lparen: '(';
Rparen: ')';
Quote: '"';

// LineComment: ';;' ~('\r' | '\n')*;
Num: ('-')? ( '0' .. '9')+;
fragment IdentFirst: ('a' .. 'z' | 'A' .. 'Z' | '.' | '$' | '_' | ';' | '/' | '*' | '@') ;
fragment IdentAfter: ('a' .. 'z' | 'A' .. 'Z' | '.' | '$' | '_' | ';' | '/' | '*' | '@'| '0'..'9');
Ident: IdentFirst IdentAfter* ;

fragment IntConst: ('-')?  ('0' .. '9')+ ( '.' ('0' .. '9')+ ( 'e' ('+' | '-') ('0' .. '9')+)?)? ;
fragment FloatConst: ('-')? ('0x')? ('0' .. '9' | 'a'..'f')+ '.' ('0' .. '9' | 'a' .. 'f')+ 'p' ('+' | '-') ('0' .. '9') ;
HexPointer: ('-')? '0x' ( '0' .. '9')+ 'p+' ('0'..'9')+;
Offset: 'offset=' ( '0' .. '9')+;
Align: 'align=' ( '0' .. '9')+;
ConstAnnotation: ';' '=' IntConst ';' ;
BlockAnnotation: ';' '@' ( '0' .. '9')+ ';' ;
ConstValue: IntConst | FloatConst;

fragment HexByteValue: '\\' ( '0' .. '9' | 'a' .. 'f') ( '0' .. '9' | 'a' .. 'f');
QuotedString: '"' ( HexByteValue | ~('"') )* '"';
