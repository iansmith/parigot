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
    ;


Num: ( '0' .. '9')+;
Id: ('a' .. 'z' | 'A' .. 'Z')+;
Whitespace : ( ' ' | '\r' '\n' | '\n' | '\t' ) {Skip();};

