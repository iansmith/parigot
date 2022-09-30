grammar jsstrip;

program: (sexpr)*;

sexpr: list
    |  atom            {fmt.Printf("matched sexpr\n");}
    ;

list:
   '('')'              {fmt.Printf("matched empty list\n");}
   | '(' members ')'   {fmt.Printf("matched list\n");}

    ;

members: (sexpr)+      {fmt.Printf("members  1\n");};

atom: Id               {fmt.Printf("ID ");}
    | Num              {fmt.Printf("NUM ");}
    | String_           {fmt.Printf("STRING ");}
    ;


Num: ( '0' .. '9')+;
HexDigit: ( '0' .. '9' | 'a' .. 'f')+;
Id: ('a' .. 'z' | 'A' .. 'Z' | '.' | '$' | '_')+;
Whitespace : ( ' ' | '\r' '\n' | '\n' | '\t' ) -> skip;
String_ : '"' (~["\\] | HexByteValue )* '"';
HexByteValue: '\\' HexDigit HexDigit;
LineComment: ';' ~[\r\n]*      -> skip;
