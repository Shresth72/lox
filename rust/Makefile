tokenize:
	cargo r -- tokenize <(echo -ne '(){};,+-*!====<=>=!=<>/.45.43+43.32.24==! .123 123.and!=or+if_else+else "Hello World"\n//comment\n/"Yo"and')

error:
	cargo r -- tokenize <(echo -ne '\n\n#!!="hello\nworld\nactually\nbanana')

comments:
	cargo r -- tokenize <(echo -ne '\n//and\nand//\n\n/\nand\n/ /')
