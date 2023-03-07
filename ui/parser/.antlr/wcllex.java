// Generated from /Users/iansmith/parigot/parigot-ui/parser/wcllex.g4 by ANTLR 4.9.2
import org.antlr.v4.runtime.Lexer;
import org.antlr.v4.runtime.CharStream;
import org.antlr.v4.runtime.Token;
import org.antlr.v4.runtime.TokenStream;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.misc.*;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast"})
public class wcllex extends Lexer {
	static { RuntimeMetaData.checkVersion("4.9.2", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		Text=1, CSS=2, Import=3, Id=4, DoubleLeftCurly=5, LCurly=6, RCurly=7, 
		LParen=8, RParen=9, Dollar=10, Comma=11, PoundComment=12, DoubleSlashComment=13, 
		Whitespace=14, ContentRawText=15, ContentDollar=16, ContentLCurly=17, 
		DoubleRightCurly=18, UninterpRawText=19, UninterpLCurly=20, UninterpRCurly=21, 
		UninterpDollar=22, VarLCurly=23, VarRCurly=24, VarId=25;
	public static final int
		CONTENT=1, UNINTERPRETED=2, VAR=3;
	public static String[] channelNames = {
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN"
	};

	public static String[] modeNames = {
		"DEFAULT_MODE", "CONTENT", "UNINTERPRETED", "VAR"
	};

	private static String[] makeRuleNames() {
		return new String[] {
			"Text", "CSS", "Import", "Id", "IdentFirst", "IdentAfter", "Digit", "DoubleLeftCurly", 
			"LCurly", "RCurly", "LParen", "RParen", "Dollar", "Comma", "PoundComment", 
			"DoubleSlashComment", "Whitespace", "ContentRawText", "ContentDollar", 
			"ContentLCurly", "DoubleRightCurly", "UninterpRawText", "UninterpLCurly", 
			"UninterpRCurly", "UninterpDollar", "VarLCurly", "VarRCurly", "VarId"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'text'", "'css'", "'import'", null, "'{{'", null, null, "'('", 
			"')'", null, "','", null, null, null, null, null, null, "'}}'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, "Text", "CSS", "Import", "Id", "DoubleLeftCurly", "LCurly", "RCurly", 
			"LParen", "RParen", "Dollar", "Comma", "PoundComment", "DoubleSlashComment", 
			"Whitespace", "ContentRawText", "ContentDollar", "ContentLCurly", "DoubleRightCurly", 
			"UninterpRawText", "UninterpLCurly", "UninterpRCurly", "UninterpDollar", 
			"VarLCurly", "VarRCurly", "VarId"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}


	public wcllex(CharStream input) {
		super(input);
		_interp = new LexerATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@Override
	public String getGrammarFileName() { return "wcllex.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public String[] getChannelNames() { return channelNames; }

	@Override
	public String[] getModeNames() { return modeNames; }

	@Override
	public ATN getATN() { return _ATN; }

	public static final String _serializedATN =
		"\3\u608b\ua72a\u8133\ub9ed\u417c\u3be7\u7786\u5964\2\33\u00bd\b\1\b\1"+
		"\b\1\b\1\4\2\t\2\4\3\t\3\4\4\t\4\4\5\t\5\4\6\t\6\4\7\t\7\4\b\t\b\4\t\t"+
		"\t\4\n\t\n\4\13\t\13\4\f\t\f\4\r\t\r\4\16\t\16\4\17\t\17\4\20\t\20\4\21"+
		"\t\21\4\22\t\22\4\23\t\23\4\24\t\24\4\25\t\25\4\26\t\26\4\27\t\27\4\30"+
		"\t\30\4\31\t\31\4\32\t\32\4\33\t\33\4\34\t\34\4\35\t\35\3\2\3\2\3\2\3"+
		"\2\3\2\3\3\3\3\3\3\3\3\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\5\3\5\7\5"+
		"S\n\5\f\5\16\5V\13\5\3\6\3\6\3\7\3\7\5\7\\\n\7\3\b\3\b\3\t\3\t\3\t\3\t"+
		"\3\t\3\n\3\n\3\13\3\13\3\f\3\f\3\r\3\r\3\16\3\16\3\17\3\17\3\20\3\20\6"+
		"\20s\n\20\r\20\16\20t\3\20\3\20\3\20\3\20\3\21\3\21\3\21\3\21\6\21\177"+
		"\n\21\r\21\16\21\u0080\3\21\3\21\3\21\3\21\3\22\6\22\u0088\n\22\r\22\16"+
		"\22\u0089\3\22\3\22\3\23\6\23\u008f\n\23\r\23\16\23\u0090\3\24\3\24\3"+
		"\24\3\24\3\25\3\25\3\25\3\25\3\26\3\26\3\26\3\26\3\26\3\27\6\27\u00a1"+
		"\n\27\r\27\16\27\u00a2\3\30\3\30\3\30\3\30\3\31\3\31\3\31\3\31\3\32\3"+
		"\32\3\32\3\32\3\33\3\33\3\34\3\34\3\34\3\34\3\35\3\35\7\35\u00b9\n\35"+
		"\f\35\16\35\u00bc\13\35\4t\u0080\2\36\6\3\b\4\n\5\f\6\16\2\20\2\22\2\24"+
		"\7\26\b\30\t\32\n\34\13\36\f \r\"\16$\17&\20(\21*\22,\23.\24\60\25\62"+
		"\26\64\27\66\308\31:\32<\33\6\2\3\4\5\7\6\2\60\60C\\aac|\4\2\f\f\17\17"+
		"\5\2\2\2\13\17\"\"\5\2&&}}\177\177\4\2}}\177\177\2\u00be\2\6\3\2\2\2\2"+
		"\b\3\2\2\2\2\n\3\2\2\2\2\f\3\2\2\2\2\24\3\2\2\2\2\26\3\2\2\2\2\30\3\2"+
		"\2\2\2\32\3\2\2\2\2\34\3\2\2\2\2\36\3\2\2\2\2 \3\2\2\2\2\"\3\2\2\2\2$"+
		"\3\2\2\2\2&\3\2\2\2\3(\3\2\2\2\3*\3\2\2\2\3,\3\2\2\2\3.\3\2\2\2\4\60\3"+
		"\2\2\2\4\62\3\2\2\2\4\64\3\2\2\2\4\66\3\2\2\2\58\3\2\2\2\5:\3\2\2\2\5"+
		"<\3\2\2\2\6>\3\2\2\2\bC\3\2\2\2\nG\3\2\2\2\fP\3\2\2\2\16W\3\2\2\2\20["+
		"\3\2\2\2\22]\3\2\2\2\24_\3\2\2\2\26d\3\2\2\2\30f\3\2\2\2\32h\3\2\2\2\34"+
		"j\3\2\2\2\36l\3\2\2\2 n\3\2\2\2\"p\3\2\2\2$z\3\2\2\2&\u0087\3\2\2\2(\u008e"+
		"\3\2\2\2*\u0092\3\2\2\2,\u0096\3\2\2\2.\u009a\3\2\2\2\60\u00a0\3\2\2\2"+
		"\62\u00a4\3\2\2\2\64\u00a8\3\2\2\2\66\u00ac\3\2\2\28\u00b0\3\2\2\2:\u00b2"+
		"\3\2\2\2<\u00b6\3\2\2\2>?\7v\2\2?@\7g\2\2@A\7z\2\2AB\7v\2\2B\7\3\2\2\2"+
		"CD\7e\2\2DE\7u\2\2EF\7u\2\2F\t\3\2\2\2GH\7k\2\2HI\7o\2\2IJ\7r\2\2JK\7"+
		"q\2\2KL\7t\2\2LM\7v\2\2MN\3\2\2\2NO\b\4\2\2O\13\3\2\2\2PT\5\16\6\2QS\5"+
		"\20\7\2RQ\3\2\2\2SV\3\2\2\2TR\3\2\2\2TU\3\2\2\2U\r\3\2\2\2VT\3\2\2\2W"+
		"X\t\2\2\2X\17\3\2\2\2Y\\\t\2\2\2Z\\\5\22\b\2[Y\3\2\2\2[Z\3\2\2\2\\\21"+
		"\3\2\2\2]^\4\62;\2^\23\3\2\2\2_`\7}\2\2`a\7}\2\2ab\3\2\2\2bc\b\t\2\2c"+
		"\25\3\2\2\2de\7}\2\2e\27\3\2\2\2fg\7\177\2\2g\31\3\2\2\2hi\7*\2\2i\33"+
		"\3\2\2\2jk\7+\2\2k\35\3\2\2\2lm\7&\2\2m\37\3\2\2\2no\7.\2\2o!\3\2\2\2"+
		"pr\7%\2\2qs\13\2\2\2rq\3\2\2\2st\3\2\2\2tu\3\2\2\2tr\3\2\2\2uv\3\2\2\2"+
		"vw\t\3\2\2wx\3\2\2\2xy\b\20\3\2y#\3\2\2\2z{\7\61\2\2{|\7\61\2\2|~\3\2"+
		"\2\2}\177\13\2\2\2~}\3\2\2\2\177\u0080\3\2\2\2\u0080\u0081\3\2\2\2\u0080"+
		"~\3\2\2\2\u0081\u0082\3\2\2\2\u0082\u0083\t\3\2\2\u0083\u0084\3\2\2\2"+
		"\u0084\u0085\b\21\3\2\u0085%\3\2\2\2\u0086\u0088\t\4\2\2\u0087\u0086\3"+
		"\2\2\2\u0088\u0089\3\2\2\2\u0089\u0087\3\2\2\2\u0089\u008a\3\2\2\2\u008a"+
		"\u008b\3\2\2\2\u008b\u008c\b\22\3\2\u008c\'\3\2\2\2\u008d\u008f\n\5\2"+
		"\2\u008e\u008d\3\2\2\2\u008f\u0090\3\2\2\2\u0090\u008e\3\2\2\2\u0090\u0091"+
		"\3\2\2\2\u0091)\3\2\2\2\u0092\u0093\7&\2\2\u0093\u0094\3\2\2\2\u0094\u0095"+
		"\b\24\4\2\u0095+\3\2\2\2\u0096\u0097\7}\2\2\u0097\u0098\3\2\2\2\u0098"+
		"\u0099\b\25\5\2\u0099-\3\2\2\2\u009a\u009b\7\177\2\2\u009b\u009c\7\177"+
		"\2\2\u009c\u009d\3\2\2\2\u009d\u009e\b\26\6\2\u009e/\3\2\2\2\u009f\u00a1"+
		"\n\6\2\2\u00a0\u009f\3\2\2\2\u00a1\u00a2\3\2\2\2\u00a2\u00a0\3\2\2\2\u00a2"+
		"\u00a3\3\2\2\2\u00a3\61\3\2\2\2\u00a4\u00a5\7}\2\2\u00a5\u00a6\3\2\2\2"+
		"\u00a6\u00a7\b\30\5\2\u00a7\63\3\2\2\2\u00a8\u00a9\7\177\2\2\u00a9\u00aa"+
		"\3\2\2\2\u00aa\u00ab\b\31\6\2\u00ab\65\3\2\2\2\u00ac\u00ad\7&\2\2\u00ad"+
		"\u00ae\3\2\2\2\u00ae\u00af\b\32\4\2\u00af\67\3\2\2\2\u00b0\u00b1\7}\2"+
		"\2\u00b19\3\2\2\2\u00b2\u00b3\7\177\2\2\u00b3\u00b4\3\2\2\2\u00b4\u00b5"+
		"\b\34\6\2\u00b5;\3\2\2\2\u00b6\u00ba\5\16\6\2\u00b7\u00b9\5\20\7\2\u00b8"+
		"\u00b7\3\2\2\2\u00b9\u00bc\3\2\2\2\u00ba\u00b8\3\2\2\2\u00ba\u00bb\3\2"+
		"\2\2\u00bb=\3\2\2\2\u00bc\u00ba\3\2\2\2\16\2\3\4\5T[t\u0080\u0089\u0090"+
		"\u00a2\u00ba\7\7\3\2\b\2\2\7\5\2\7\4\2\6\2\2";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}