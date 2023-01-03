const { spawnSync } = require('child_process')
const vscode = require('vscode')

const parsedDocs = {}

function callQuickonf(doc) {
	// Do not process the same doc version twice
	if (!(parsedDocs[doc.uri]) || doc.version !== parsedDocs[doc.uri].version) {
		parsedDocs[doc.uri] = {
			version: doc.version,
			data: new Promise((resolve) => {
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
			}),
		}
	}
	return parsedDocs[doc.uri].data
}

module.exports = callQuickonf
