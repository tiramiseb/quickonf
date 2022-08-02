const completion = require('./completion')
const diagnostic = require('./diagnostic')
const hover = require('./hover')
const semantictokens = require('./semantictokens')

function activate(ctx) {
	completion()
	diagnostic(ctx)
	hover()
	semantictokens()
}

module.exports = { activate }
