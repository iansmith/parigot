lexer grammar wcllex;

@lexer::header {
// at top of file
import "log"
}
// keywords
Text: '@text';
CSS: '@css';
Import: '@preamble';
Doc: '@doc';
Local: '@local';
Global: '@global';
Extern: '@extern';
Mvc: '@mvc';
Pre: '@pre';
Post: '@post';
Wcl: '@wcl';
Event: '@event';
Model: '@model';
View: '@view';
ViewCollection: '@collection';
Controller: '@controller';


//ids
//TypeId: (TypeStarter+)? IdentFirst (IdentAfter)*;
Id: IdentFirst (IdentAfter)*;

TypeStarter: '[' | ']'|'*';

// consistent def of Ident
fragment IdentFirst: ('a' .. 'z' | 'A' .. 'Z' | '.' | '_' | '-');

fragment IdentAfter: (
		'a' .. 'z'
		| 'A' .. 'Z'
		| '_'
		| Digit
	);

Version: Digit+ Dot Digit+ Dot Digit+;
fragment Digit: '0' ..'9';

DoubleLess: '<<' -> pushMode(GrabText);
Arrow: '->';
LCurly: '{' ;
RCurly: '}';
LParen: '(';
RParen: ')';
Dollar: '$';
Comma: ',';
Colon: ':';
LessThan: '<';
GreaterThan: '>';
Dot: '.';
Hash: '#';
Dash: '-';
Caret: '^';
Semi: ';';
Plus: '+';
BackTick: '`';
StringLit: '"' ( Esc | ~[\\"] )* '"';
fragment Esc : '\\"' | '\\\\' ;

DoubleSlashComment: '//' .+? [\n\r] -> skip;
Whitespace: [ \n\r\t\u000B\u000C\u0000]+ -> skip;

mode GrabText;
GrabDollar: '$';
GrabLCurly: '{' ;
GrabRCurly: '}';
GrabColon: ':';
GrabComma: ',';
GrabDot: '.';
GrabLParen: '(';
GrabRParen: ')';
GrabDoubleGreater: '>>' -> popMode;
GrabId: IdentFirst (IdentAfter)*;
RawText: 
	~[${}()>]+ {log.Printf("xxx got '%s' for raw text",l.GetText())}
	;

