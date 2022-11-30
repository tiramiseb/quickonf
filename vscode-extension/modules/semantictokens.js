const readConf = require('./readconf')
const vscode = require('vscode')

const legend = new vscode.SemanticTokensLegend([
	'keyword', 'function', 'variable', 'operator', 'unknown'
], [])

const quickonfToVscodeTokens = {
	"keyword": 0,
	"groupname": 1,
	"variable": 2,
	"function": 1,
	"operator": 3,
	"unknown": 4
}

const provider = {
	provideDocumentSemanticTokens(doc) {
		return readConf(doc).then(result => {
			builder = new vscode.SemanticTokensBuilder(legend)
			result.tokens.forEach(tok => {
				builder.push(tok.line, tok.start, tok.length, quickonfToVscodeTokens[tok.type])
			})
			return builder.build()
		})
	}
}

function activate() {
	vscode.languages.registerDocumentSemanticTokensProvider('quickonf', provider, legend)
}

module.exports = activate
