lexer grammar wcllex;

// keywords
Text: '@text';
CSS: '@css';
Import: '@preamble';
Doc: '@doc';
Local: '@local';
Global: '@global';
Extern: '@extern';
Pre: '@pre';
Post: '@post';
Wcl: '@wcl';
Event: '@event';

//ids
//TypeId: (TypeStarter+)? IdentFirst (IdentAfter)*;
Id: (TypeStarter+)? IdentFirst (IdentAfter)*;

fragment TypeStarter: '[' | ']'|'*';

// consistent def of Ident
fragment IdentFirst: ('a' .. 'z' | 'A' .. 'Z' | '.' | '_' | '-');

fragment IdentAfter: (
		'a' .. 'z'
		| 'A' .. 'Z'
		| '.'
		| '_'
		| '-'
		| Digit
	);

Version: Digit+ Dot Digit+ Dot Digit+;
fragment Digit: '0' ..'9';

LCurly: '{' -> pushMode(UNINTERPRETED);
RCurly: '}';
LParen: '(';
RParen: ')';
Dollar: '${' -> pushMode(VAR);
Comma: ',';
LessThan: '<';
GreaterThan: '>';
Dot: '.';
Hash: '#';
Dash: '-';
Semi: ';';
Plus: '+';
BackTick: '`' -> pushMode(CONTENT);
StringLit: '"' ( Esc | ~[\\"] )* '"';
fragment Esc : '\\"' | '\\\\' ;

DoubleSlashComment: '//' .+? [\n\r] -> skip;
Whitespace: [ \n\r\t\u000B\u000C\u0000]+ -> skip;

mode CONTENT;
ContentRawText: ~[${`]+;
ContentDollar: '${' -> pushMode(VAR);
ContentBackTick: '`' -> popMode;

mode UNINTERPRETED;
UninterpRawText: ~[${}]+;
UninterpDollar: '${' -> pushMode(VAR) ;
UninterpLCurly: '{' -> pushMode(UNINTERPRETED);
UninterpRCurly: '}' -> popMode;

mode VAR;
VarRCurly: '}' -> popMode;
VarId: IdentFirst (IdentAfter)*;

