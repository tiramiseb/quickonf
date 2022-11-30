const { spawnSync } = require('child_process')
const vscode = require('vscode')

const docVersion = {}
const docsParsers = {}

function callQuickonf(doc) {
	// Do not process the same doc version twice
	if (doc.version !== docVersion[doc.uri] || docsParsers[doc.uri]) {
		docsParsers[doc.uri] = new Promise((resolve) => {
			docVersion[doc.uri] = doc.version
			output = spawnSync(
				vscode.workspace.getConfiguration("quickonf").get("executablePath"),
				["--check-stdin"],
				{ input: doc.getText() }
			)
			if (output.error) {
				let message = output.error.message
				if (output.error.code === "ENOENT") {
					message = "Error: wrong path for quickonf executable"
				}
				vscode.window.showErrorMessage(message, "Open settings").then((opt) => {
					if (opt === "Open settings") {
						vscode.commands.executeCommand("workbench.action.openSettings", "quickonf.executablePath")
					}
				})
			}
			resolve(JSON.parse(output.stdout))
		})
	}
	return docsParsers[doc.uri]
}

module.exports = callQuickonf
