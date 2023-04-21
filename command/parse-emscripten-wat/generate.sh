#!/bin/sh

alias antlr4='java -Xmx500M -cp "../../ui/antlr-4.11.1-complete.jar:$CLASSPATH" org.antlr.v4.Tool'
antlr4 -Dlanguage=Go -listener -visitor -package main sexpr.g4
