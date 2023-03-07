// Generated from /Users/iansmith/parigot/parigot-ui/parser/wcl.g4 by ANTLR 4.9.2
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast"})
public class wcl extends Parser {
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
		RULE_program = 0, RULE_import_section = 1, RULE_text_section = 2, RULE_text_decl = 3, 
		RULE_text_top = 4, RULE_text_content = 5, RULE_var_subs = 6, RULE_sub = 7, 
		RULE_uninterp = 8, RULE_nestedUninterp = 9, RULE_uninterpVar = 10, RULE_param_spec = 11, 
		RULE_param_seq = 12, RULE_param_rest = 13, RULE_css_section = 14, RULE_css_file = 15, 
		RULE_css_decl = 16;
	private static String[] makeRuleNames() {
		return new String[] {
			"program", "import_section", "text_section", "text_decl", "text_top", 
			"text_content", "var_subs", "sub", "uninterp", "nestedUninterp", "uninterpVar", 
			"param_spec", "param_seq", "param_rest", "css_section", "css_file", "css_decl"
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

	@Override
	public String getGrammarFileName() { return "wcl.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public wcl(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	public static class ProgramContext extends ParserRuleContext {
		public *Program p;
		public Import_sectionContext i;
		public Text_sectionContext t;
		public TerminalNode EOF() { return getToken(wcl.EOF, 0); }
		public Css_sectionContext css_section() {
			return getRuleContext(Css_sectionContext.class,0);
		}
		public Import_sectionContext import_section() {
			return getRuleContext(Import_sectionContext.class,0);
		}
		public Text_sectionContext text_section() {
			return getRuleContext(Text_sectionContext.class,0);
		}
		public ProgramContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_program; }
	}

	public final ProgramContext program() throws RecognitionException {
		ProgramContext _localctx = new ProgramContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_program);

				_localctx.p=NewProgram()	
			
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(35);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==Import) {
				{
				setState(34);
				((ProgramContext)_localctx).i = import_section();
				}
			}

			 _localctx.p.ImportSection = localctx.GetI().GetSection()
				
			setState(39);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==Text) {
				{
				setState(38);
				((ProgramContext)_localctx).t = text_section();
				}
			}

			 _localctx.p.TextSection = localctx.GetT().GetSection()
					
			setState(43);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==CSS) {
				{
				setState(42);
				css_section();
				}
			}

			setState(45);
			match(EOF);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Import_sectionContext extends ParserRuleContext {
		public *ImportSectionNode section;
		public UninterpContext u;
		public TerminalNode Import() { return getToken(wcl.Import, 0); }
		public UninterpContext uninterp() {
			return getRuleContext(UninterpContext.class,0);
		}
		public Import_sectionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_import_section; }
	}

	public final Import_sectionContext import_section() throws RecognitionException {
		Import_sectionContext _localctx = new Import_sectionContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_import_section);

				_localctx.section = NewImportSectionNode()
			
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(47);
			match(Import);
			setState(48);
			((Import_sectionContext)_localctx).u = uninterp();

					_localctx.section.Text = localctx.GetU().GetItem()[0];
				
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Text_sectionContext extends ParserRuleContext {
		public *TextSectionNode section;
		public Text_declContext d;
		public TerminalNode Text() { return getToken(wcl.Text, 0); }
		public List<Text_declContext> text_decl() {
			return getRuleContexts(Text_declContext.class);
		}
		public Text_declContext text_decl(int i) {
			return getRuleContext(Text_declContext.class,i);
		}
		public Text_sectionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_text_section; }
	}

	public final Text_sectionContext text_section() throws RecognitionException {
		Text_sectionContext _localctx = new Text_sectionContext(_ctx, getState());
		enterRule(_localctx, 4, RULE_text_section);

				_localctx.section = NewTextSectionNode()
			
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(51);
			match(Text);
			setState(57);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==Id) {
				{
				{
				setState(52);
				((Text_sectionContext)_localctx).d = text_decl();
				_localctx.section.Func=append(_localctx.section.Func,localctx.GetD().GetF())
				}
				}
				setState(59);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Text_declContext extends ParserRuleContext {
		public *TextFuncNode f;
		public Token i;
		public Text_topContext t;
		public TerminalNode Id() { return getToken(wcl.Id, 0); }
		public Text_topContext text_top() {
			return getRuleContext(Text_topContext.class,0);
		}
		public Param_specContext param_spec() {
			return getRuleContext(Param_specContext.class,0);
		}
		public Text_declContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_text_decl; }
	}

	public final Text_declContext text_decl() throws RecognitionException {
		Text_declContext _localctx = new Text_declContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_text_decl);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(60);
			((Text_declContext)_localctx).i = match(Id);
			setState(62);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==LParen) {
				{
				setState(61);
				param_spec();
				}
			}

			setState(64);
			((Text_declContext)_localctx).t = text_top();
			 
					_localctx.f=NewTextFuncNode(localctx.GetI().GetText(),localctx.GetT().GetItem())
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Text_topContext extends ParserRuleContext {
		public []TextItem item;
		public Text_contentContext content;
		public TerminalNode DoubleLeftCurly() { return getToken(wcl.DoubleLeftCurly, 0); }
		public TerminalNode DoubleRightCurly() { return getToken(wcl.DoubleRightCurly, 0); }
		public Text_contentContext text_content() {
			return getRuleContext(Text_contentContext.class,0);
		}
		public Text_topContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_text_top; }
	}

	public final Text_topContext text_top() throws RecognitionException {
		Text_topContext _localctx = new Text_topContext(_ctx, getState());
		enterRule(_localctx, 8, RULE_text_top);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(67);
			match(DoubleLeftCurly);
			setState(72);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,5,_ctx) ) {
			case 1:
				{
				setState(68);
				((Text_topContext)_localctx).content = text_content();
				_localctx.item = localctx.GetContent().GetItem()
				}
				break;
			case 2:
				{
				}
				break;
			}
			setState(74);
			match(DoubleRightCurly);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Text_contentContext extends ParserRuleContext {
		public []TextItem item;
		public Token c;
		public Var_subsContext v;
		public UninterpContext u;
		public List<TerminalNode> ContentRawText() { return getTokens(wcl.ContentRawText); }
		public TerminalNode ContentRawText(int i) {
			return getToken(wcl.ContentRawText, i);
		}
		public List<Var_subsContext> var_subs() {
			return getRuleContexts(Var_subsContext.class);
		}
		public Var_subsContext var_subs(int i) {
			return getRuleContext(Var_subsContext.class,i);
		}
		public List<UninterpContext> uninterp() {
			return getRuleContexts(UninterpContext.class);
		}
		public UninterpContext uninterp(int i) {
			return getRuleContext(UninterpContext.class,i);
		}
		public Text_contentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_text_content; }
	}

	public final Text_contentContext text_content() throws RecognitionException {
		Text_contentContext _localctx = new Text_contentContext(_ctx, getState());
		enterRule(_localctx, 10, RULE_text_content);

				_localctx.item = []TextItem{}
			
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(86);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << ContentRawText) | (1L << ContentDollar) | (1L << ContentLCurly))) != 0)) {
				{
				setState(84);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case ContentRawText:
					{
					setState(76);
					((Text_contentContext)_localctx).c = match(ContentRawText);
					 _localctx.item = append(_localctx.item, NewTextConstant(localctx.GetC().GetText()))
					}
					break;
				case ContentDollar:
					{
					setState(78);
					((Text_contentContext)_localctx).v = var_subs();
					 
					}
					break;
				case ContentLCurly:
					{
					setState(81);
					((Text_contentContext)_localctx).u = uninterp();
					 _localctx.item = append(_localctx.item, localctx.GetU().GetItem()...)
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				}
				setState(88);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Var_subsContext extends ParserRuleContext {
		public TerminalNode ContentDollar() { return getToken(wcl.ContentDollar, 0); }
		public SubContext sub() {
			return getRuleContext(SubContext.class,0);
		}
		public Var_subsContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_var_subs; }
	}

	public final Var_subsContext var_subs() throws RecognitionException {
		Var_subsContext _localctx = new Var_subsContext(_ctx, getState());
		enterRule(_localctx, 12, RULE_var_subs);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(89);
			match(ContentDollar);
			setState(90);
			sub();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SubContext extends ParserRuleContext {
		public TerminalNode VarLCurly() { return getToken(wcl.VarLCurly, 0); }
		public TerminalNode VarId() { return getToken(wcl.VarId, 0); }
		public TerminalNode VarRCurly() { return getToken(wcl.VarRCurly, 0); }
		public SubContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_sub; }
	}

	public final SubContext sub() throws RecognitionException {
		SubContext _localctx = new SubContext(_ctx, getState());
		enterRule(_localctx, 14, RULE_sub);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(92);
			match(VarLCurly);
			setState(93);
			match(VarId);
			setState(94);
			match(VarRCurly);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class UninterpContext extends ParserRuleContext {
		public []TextItem item;
		public Token c;
		public NestedUninterpContext u;
		public TerminalNode ContentLCurly() { return getToken(wcl.ContentLCurly, 0); }
		public TerminalNode UninterpRCurly() { return getToken(wcl.UninterpRCurly, 0); }
		public List<TerminalNode> UninterpDollar() { return getTokens(wcl.UninterpDollar); }
		public TerminalNode UninterpDollar(int i) {
			return getToken(wcl.UninterpDollar, i);
		}
		public List<SubContext> sub() {
			return getRuleContexts(SubContext.class);
		}
		public SubContext sub(int i) {
			return getRuleContext(SubContext.class,i);
		}
		public List<TerminalNode> UninterpRawText() { return getTokens(wcl.UninterpRawText); }
		public TerminalNode UninterpRawText(int i) {
			return getToken(wcl.UninterpRawText, i);
		}
		public List<NestedUninterpContext> nestedUninterp() {
			return getRuleContexts(NestedUninterpContext.class);
		}
		public NestedUninterpContext nestedUninterp(int i) {
			return getRuleContext(NestedUninterpContext.class,i);
		}
		public UninterpContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_uninterp; }
	}

	public final UninterpContext uninterp() throws RecognitionException {
		UninterpContext _localctx = new UninterpContext(_ctx, getState());
		enterRule(_localctx, 16, RULE_uninterp);

				_localctx.item=[]TextItem{}
			
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(96);
			match(ContentLCurly);
			setState(104); 
			_errHandler.sync(this);
			_la = _input.LA(1);
			do {
				{
				setState(104);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case UninterpRawText:
					{
					setState(97);
					((UninterpContext)_localctx).c = match(UninterpRawText);
					 _localctx.item = append(_localctx.item, NewTextConstant(localctx.GetC().GetText()))
					}
					break;
				case UninterpLCurly:
					{
					setState(99);
					((UninterpContext)_localctx).u = nestedUninterp();
					 _localctx.item = append(_localctx.item, localctx.GetU().GetItem()...)
					}
					break;
				case UninterpDollar:
					{
					setState(102);
					match(UninterpDollar);
					setState(103);
					sub();
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				}
				setState(106); 
				_errHandler.sync(this);
				_la = _input.LA(1);
			} while ( (((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << UninterpRawText) | (1L << UninterpLCurly) | (1L << UninterpDollar))) != 0) );
			setState(108);
			match(UninterpRCurly);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class NestedUninterpContext extends ParserRuleContext {
		public []TextItem item;
		public Token c;
		public NestedUninterpContext u;
		public TerminalNode UninterpLCurly() { return getToken(wcl.UninterpLCurly, 0); }
		public TerminalNode UninterpRCurly() { return getToken(wcl.UninterpRCurly, 0); }
		public List<TerminalNode> UninterpDollar() { return getTokens(wcl.UninterpDollar); }
		public TerminalNode UninterpDollar(int i) {
			return getToken(wcl.UninterpDollar, i);
		}
		public List<SubContext> sub() {
			return getRuleContexts(SubContext.class);
		}
		public SubContext sub(int i) {
			return getRuleContext(SubContext.class,i);
		}
		public List<TerminalNode> UninterpRawText() { return getTokens(wcl.UninterpRawText); }
		public TerminalNode UninterpRawText(int i) {
			return getToken(wcl.UninterpRawText, i);
		}
		public List<NestedUninterpContext> nestedUninterp() {
			return getRuleContexts(NestedUninterpContext.class);
		}
		public NestedUninterpContext nestedUninterp(int i) {
			return getRuleContext(NestedUninterpContext.class,i);
		}
		public NestedUninterpContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_nestedUninterp; }
	}

	public final NestedUninterpContext nestedUninterp() throws RecognitionException {
		NestedUninterpContext _localctx = new NestedUninterpContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_nestedUninterp);

				_localctx.item=[]TextItem{}
			
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(110);
			match(UninterpLCurly);

			setState(119); 
			_errHandler.sync(this);
			_la = _input.LA(1);
			do {
				{
				setState(119);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case UninterpRawText:
					{
					setState(112);
					((NestedUninterpContext)_localctx).c = match(UninterpRawText);
					 _localctx.item = append(_localctx.item, NewTextConstant(localctx.GetC().GetText()))
					}
					break;
				case UninterpDollar:
					{
					setState(114);
					match(UninterpDollar);
					setState(115);
					sub();
					}
					break;
				case UninterpLCurly:
					{
					setState(116);
					((NestedUninterpContext)_localctx).u = nestedUninterp();
					 _localctx.item = append(_localctx.item, localctx.GetU().GetItem()...)
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				}
				setState(121); 
				_errHandler.sync(this);
				_la = _input.LA(1);
			} while ( (((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << UninterpRawText) | (1L << UninterpLCurly) | (1L << UninterpDollar))) != 0) );
			setState(123);
			match(UninterpRCurly);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class UninterpVarContext extends ParserRuleContext {
		public TerminalNode UninterpDollar() { return getToken(wcl.UninterpDollar, 0); }
		public TerminalNode Id() { return getToken(wcl.Id, 0); }
		public UninterpVarContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_uninterpVar; }
	}

	public final UninterpVarContext uninterpVar() throws RecognitionException {
		UninterpVarContext _localctx = new UninterpVarContext(_ctx, getState());
		enterRule(_localctx, 20, RULE_uninterpVar);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(125);
			match(UninterpDollar);
			setState(126);
			match(Id);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Param_specContext extends ParserRuleContext {
		public TerminalNode LParen() { return getToken(wcl.LParen, 0); }
		public Param_seqContext param_seq() {
			return getRuleContext(Param_seqContext.class,0);
		}
		public TerminalNode RParen() { return getToken(wcl.RParen, 0); }
		public Param_specContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_param_spec; }
	}

	public final Param_specContext param_spec() throws RecognitionException {
		Param_specContext _localctx = new Param_specContext(_ctx, getState());
		enterRule(_localctx, 22, RULE_param_spec);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(128);
			match(LParen);
			setState(129);
			param_seq();
			setState(130);
			match(RParen);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Param_seqContext extends ParserRuleContext {
		public List<TerminalNode> Id() { return getTokens(wcl.Id); }
		public TerminalNode Id(int i) {
			return getToken(wcl.Id, i);
		}
		public List<TerminalNode> Comma() { return getTokens(wcl.Comma); }
		public TerminalNode Comma(int i) {
			return getToken(wcl.Comma, i);
		}
		public Param_seqContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_param_seq; }
	}

	public final Param_seqContext param_seq() throws RecognitionException {
		Param_seqContext _localctx = new Param_seqContext(_ctx, getState());
		enterRule(_localctx, 24, RULE_param_seq);
		try {
			int _alt;
			setState(140);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case Id:
				enterOuterAlt(_localctx, 1);
				{
				setState(134); 
				_errHandler.sync(this);
				_alt = 1+1;
				do {
					switch (_alt) {
					case 1+1:
						{
						{
						setState(132);
						match(Id);
						setState(133);
						match(Comma);
						}
						}
						break;
					default:
						throw new NoViableAltException(this);
					}
					setState(136); 
					_errHandler.sync(this);
					_alt = getInterpreter().adaptivePredict(_input,12,_ctx);
				} while ( _alt!=1 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER );
				setState(138);
				match(Id);
				}
				break;
			case EOF:
			case RParen:
				enterOuterAlt(_localctx, 2);
				{
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Param_restContext extends ParserRuleContext {
		public TerminalNode Comma() { return getToken(wcl.Comma, 0); }
		public Param_seqContext param_seq() {
			return getRuleContext(Param_seqContext.class,0);
		}
		public Param_restContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_param_rest; }
	}

	public final Param_restContext param_rest() throws RecognitionException {
		Param_restContext _localctx = new Param_restContext(_ctx, getState());
		enterRule(_localctx, 26, RULE_param_rest);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(142);
			match(Comma);
			setState(143);
			param_seq();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Css_sectionContext extends ParserRuleContext {
		public TerminalNode CSS() { return getToken(wcl.CSS, 0); }
		public List<Css_fileContext> css_file() {
			return getRuleContexts(Css_fileContext.class);
		}
		public Css_fileContext css_file(int i) {
			return getRuleContext(Css_fileContext.class,i);
		}
		public Css_sectionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_css_section; }
	}

	public final Css_sectionContext css_section() throws RecognitionException {
		Css_sectionContext _localctx = new Css_sectionContext(_ctx, getState());
		enterRule(_localctx, 28, RULE_css_section);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(145);
			match(CSS);
			setState(149);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==Id) {
				{
				{
				setState(146);
				css_file();
				}
				}
				setState(151);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Css_fileContext extends ParserRuleContext {
		public TerminalNode Id() { return getToken(wcl.Id, 0); }
		public TerminalNode LCurly() { return getToken(wcl.LCurly, 0); }
		public TerminalNode RCurly() { return getToken(wcl.RCurly, 0); }
		public List<Css_declContext> css_decl() {
			return getRuleContexts(Css_declContext.class);
		}
		public Css_declContext css_decl(int i) {
			return getRuleContext(Css_declContext.class,i);
		}
		public Css_fileContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_css_file; }
	}

	public final Css_fileContext css_file() throws RecognitionException {
		Css_fileContext _localctx = new Css_fileContext(_ctx, getState());
		enterRule(_localctx, 30, RULE_css_file);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(152);
			match(Id);
			setState(153);
			match(LCurly);
			setState(157);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==Id) {
				{
				{
				setState(154);
				css_decl();
				}
				}
				setState(159);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(160);
			match(RCurly);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Css_declContext extends ParserRuleContext {
		public TerminalNode Id() { return getToken(wcl.Id, 0); }
		public Css_declContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_css_decl; }
	}

	public final Css_declContext css_decl() throws RecognitionException {
		Css_declContext _localctx = new Css_declContext(_ctx, getState());
		enterRule(_localctx, 32, RULE_css_decl);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(162);
			match(Id);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static final String _serializedATN =
		"\3\u608b\ua72a\u8133\ub9ed\u417c\u3be7\u7786\u5964\3\33\u00a7\4\2\t\2"+
		"\4\3\t\3\4\4\t\4\4\5\t\5\4\6\t\6\4\7\t\7\4\b\t\b\4\t\t\t\4\n\t\n\4\13"+
		"\t\13\4\f\t\f\4\r\t\r\4\16\t\16\4\17\t\17\4\20\t\20\4\21\t\21\4\22\t\22"+
		"\3\2\5\2&\n\2\3\2\3\2\5\2*\n\2\3\2\3\2\5\2.\n\2\3\2\3\2\3\3\3\3\3\3\3"+
		"\3\3\4\3\4\3\4\3\4\7\4:\n\4\f\4\16\4=\13\4\3\5\3\5\5\5A\n\5\3\5\3\5\3"+
		"\5\3\6\3\6\3\6\3\6\3\6\5\6K\n\6\3\6\3\6\3\7\3\7\3\7\3\7\3\7\3\7\3\7\3"+
		"\7\7\7W\n\7\f\7\16\7Z\13\7\3\b\3\b\3\b\3\t\3\t\3\t\3\t\3\n\3\n\3\n\3\n"+
		"\3\n\3\n\3\n\3\n\6\nk\n\n\r\n\16\nl\3\n\3\n\3\13\3\13\3\13\3\13\3\13\3"+
		"\13\3\13\3\13\3\13\6\13z\n\13\r\13\16\13{\3\13\3\13\3\f\3\f\3\f\3\r\3"+
		"\r\3\r\3\r\3\16\3\16\6\16\u0089\n\16\r\16\16\16\u008a\3\16\3\16\5\16\u008f"+
		"\n\16\3\17\3\17\3\17\3\20\3\20\7\20\u0096\n\20\f\20\16\20\u0099\13\20"+
		"\3\21\3\21\3\21\7\21\u009e\n\21\f\21\16\21\u00a1\13\21\3\21\3\21\3\22"+
		"\3\22\3\22\3\u008a\2\23\2\4\6\b\n\f\16\20\22\24\26\30\32\34\36 \"\2\2"+
		"\2\u00a8\2%\3\2\2\2\4\61\3\2\2\2\6\65\3\2\2\2\b>\3\2\2\2\nE\3\2\2\2\f"+
		"X\3\2\2\2\16[\3\2\2\2\20^\3\2\2\2\22b\3\2\2\2\24p\3\2\2\2\26\177\3\2\2"+
		"\2\30\u0082\3\2\2\2\32\u008e\3\2\2\2\34\u0090\3\2\2\2\36\u0093\3\2\2\2"+
		" \u009a\3\2\2\2\"\u00a4\3\2\2\2$&\5\4\3\2%$\3\2\2\2%&\3\2\2\2&\'\3\2\2"+
		"\2\')\b\2\1\2(*\5\6\4\2)(\3\2\2\2)*\3\2\2\2*+\3\2\2\2+-\b\2\1\2,.\5\36"+
		"\20\2-,\3\2\2\2-.\3\2\2\2./\3\2\2\2/\60\7\2\2\3\60\3\3\2\2\2\61\62\7\5"+
		"\2\2\62\63\5\22\n\2\63\64\b\3\1\2\64\5\3\2\2\2\65;\7\3\2\2\66\67\5\b\5"+
		"\2\678\b\4\1\28:\3\2\2\29\66\3\2\2\2:=\3\2\2\2;9\3\2\2\2;<\3\2\2\2<\7"+
		"\3\2\2\2=;\3\2\2\2>@\7\6\2\2?A\5\30\r\2@?\3\2\2\2@A\3\2\2\2AB\3\2\2\2"+
		"BC\5\n\6\2CD\b\5\1\2D\t\3\2\2\2EJ\7\7\2\2FG\5\f\7\2GH\b\6\1\2HK\3\2\2"+
		"\2IK\3\2\2\2JF\3\2\2\2JI\3\2\2\2KL\3\2\2\2LM\7\24\2\2M\13\3\2\2\2NO\7"+
		"\21\2\2OW\b\7\1\2PQ\5\16\b\2QR\b\7\1\2RW\3\2\2\2ST\5\22\n\2TU\b\7\1\2"+
		"UW\3\2\2\2VN\3\2\2\2VP\3\2\2\2VS\3\2\2\2WZ\3\2\2\2XV\3\2\2\2XY\3\2\2\2"+
		"Y\r\3\2\2\2ZX\3\2\2\2[\\\7\22\2\2\\]\5\20\t\2]\17\3\2\2\2^_\7\31\2\2_"+
		"`\7\33\2\2`a\7\32\2\2a\21\3\2\2\2bj\7\23\2\2cd\7\25\2\2dk\b\n\1\2ef\5"+
		"\24\13\2fg\b\n\1\2gk\3\2\2\2hi\7\30\2\2ik\5\20\t\2jc\3\2\2\2je\3\2\2\2"+
		"jh\3\2\2\2kl\3\2\2\2lj\3\2\2\2lm\3\2\2\2mn\3\2\2\2no\7\27\2\2o\23\3\2"+
		"\2\2pq\7\26\2\2qy\b\13\1\2rs\7\25\2\2sz\b\13\1\2tu\7\30\2\2uz\5\20\t\2"+
		"vw\5\24\13\2wx\b\13\1\2xz\3\2\2\2yr\3\2\2\2yt\3\2\2\2yv\3\2\2\2z{\3\2"+
		"\2\2{y\3\2\2\2{|\3\2\2\2|}\3\2\2\2}~\7\27\2\2~\25\3\2\2\2\177\u0080\7"+
		"\30\2\2\u0080\u0081\7\6\2\2\u0081\27\3\2\2\2\u0082\u0083\7\n\2\2\u0083"+
		"\u0084\5\32\16\2\u0084\u0085\7\13\2\2\u0085\31\3\2\2\2\u0086\u0087\7\6"+
		"\2\2\u0087\u0089\7\r\2\2\u0088\u0086\3\2\2\2\u0089\u008a\3\2\2\2\u008a"+
		"\u008b\3\2\2\2\u008a\u0088\3\2\2\2\u008b\u008c\3\2\2\2\u008c\u008f\7\6"+
		"\2\2\u008d\u008f\3\2\2\2\u008e\u0088\3\2\2\2\u008e\u008d\3\2\2\2\u008f"+
		"\33\3\2\2\2\u0090\u0091\7\r\2\2\u0091\u0092\5\32\16\2\u0092\35\3\2\2\2"+
		"\u0093\u0097\7\4\2\2\u0094\u0096\5 \21\2\u0095\u0094\3\2\2\2\u0096\u0099"+
		"\3\2\2\2\u0097\u0095\3\2\2\2\u0097\u0098\3\2\2\2\u0098\37\3\2\2\2\u0099"+
		"\u0097\3\2\2\2\u009a\u009b\7\6\2\2\u009b\u009f\7\b\2\2\u009c\u009e\5\""+
		"\22\2\u009d\u009c\3\2\2\2\u009e\u00a1\3\2\2\2\u009f\u009d\3\2\2\2\u009f"+
		"\u00a0\3\2\2\2\u00a0\u00a2\3\2\2\2\u00a1\u009f\3\2\2\2\u00a2\u00a3\7\t"+
		"\2\2\u00a3!\3\2\2\2\u00a4\u00a5\7\6\2\2\u00a5#\3\2\2\2\22%)-;@JVXjly{"+
		"\u008a\u008e\u0097\u009f";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}