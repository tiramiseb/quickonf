const completion = require('modules/completion')
const diagnostic = require('modules/diagnostic')
const hover = require('modules/hover')
const semantictokens = require('modules/semantictokens')

function activate(ctx) {
	completion()
	diagnostic(ctx)
	hover()
	semantictokens()
}

module.exports = { activate }
