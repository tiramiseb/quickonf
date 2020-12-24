---
title: JSON
---

| Instruction  | Action                                 | Arguments                      |
| ------------ | -------------------------------------- | ------------------------------ |
| `json-build` | Build a JSON structure from parameters | List of parameters (see below) |
| `json-get`   | Get a data from a JSON structure       | Map of parameters (see below)  |

Example for `json-build` and `json-get` may be found in the [Cura for Dagoma]({% link _recipes/curafordagoma.md %}) recipe.

## JSON build

The parameters are simply a list of strings like:

- `key=value`: would set the given key to the given value

The (optional) special key "store" may be used to give a key in the store.

This instruction is very simple for the moment. It can only define string, and only at the root of the JSON object. No embedded object, no int, no float, no list, etc.

## JSON get

Mandatory parameters are:

- `from`: the data to read
- `key`: the key to extract from the JSON (uses the [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md))

Optional parameters are:

- `store`: the store key, where to store the result
