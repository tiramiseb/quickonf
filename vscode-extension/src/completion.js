const readConf = require('./readconf')
const vscode = require('vscode')

const kinds = {
	"command": vscode.CompletionItemKind.Function,
	"variable": vscode.CompletionItemKind.Variable,
}

function tokenAt(tokens, line, col) {
	return tokens.find(tok => tok.line == line && tok.start <= col && tok.end >= col)
}

completionProvider = {
	provideCompletionItems(doc, position) {
		return readConf(doc).then(result => {
			const tok = tokenAt(result.tokens, position.line, position.character)
			if (!tok) {
				return
			}
			const range = new vscode.Range(tok.line, tok.start, tok.line, tok.end)
			return tok?.completion?.map(entry => {
				item = new vscode.CompletionItem(entry.label, kinds[entry.kind])
				item.range = range
				return item
			})
		})
	}
}

function activate() {
	vscode.languages.registerCompletionItemProvider('quickonf', completionProvider)
}

module.exports = activate
