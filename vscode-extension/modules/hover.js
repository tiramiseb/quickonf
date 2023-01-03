const readConf = require('./readconf')
const vscode = require('vscode')

function tokenAt(tokens, line, col) {
	return tokens.find(tok => tok.line == line && tok.start <= col && tok.end >= col)
}

hoverProvider = {
	provideHover(doc, position) {
		return readConf(doc).then(result => {
			const tok = tokenAt(result.tokens, position.line, position.character)
			if (tok?.help) {
				return { contents: [ tok.help ] }
			}
		})
	}
}

function activate() {
	vscode.languages.registerHoverProvider('quickonf', hoverProvider)
}

module.exports = activate
