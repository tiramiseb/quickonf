const readConf = require('./readconf')
const vscode = require('vscode')

let diagnosticCollection

function parse(doc) {
	// Only process quickonf documents
	if (doc.languageId !== "quickonf") {
		return
	}
	readConf(doc).then((result) => {
		const diags = result.errors.map((diag) => {
			return new vscode.Diagnostic(
				new vscode.Range(diag.line, diag.start, diag.line, diag.end),
				diag.content + ": " + diag.message, diag.severity
			)
		})
		diagnosticCollection.set(doc.uri, diags)
	})
}

function parseOnChange(change) {
	parse(change.document)
}
function parseOnLoad(document) {
	parse(document)
}

function activate(ctx) {
	diagnosticCollection = vscode.languages.createDiagnosticCollection('quickonf')
	ctx.subscriptions.push(
		vscode.workspace.onDidOpenTextDocument(parseOnLoad),
		vscode.workspace.onDidChangeTextDocument(parseOnChange),
		diagnosticCollection,
	)
}

module.exports = activate
