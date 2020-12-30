---
title: Web
---

| Instruction      | Action                      | Arguments                       | Dry run |
| ---------------- | --------------------------- | ------------------------------- | ------- |
| `parse-web-page` | Find a string in a web page | Parameters as a map (see below) | -       |
| `http-get`       | Execute a GET request       | Parameters as a map (see below) | -       |
| `http-post`      | Execute a POST request      | Parameters as a map (see below) | -       |

## Parse web page

When parsing a web page, the two following parameters are mandatory:

- `url`: URL of the webpage to parse
- `regexp`: the regexp to search

The following parameter is optional:

- `store`: key name in the store for the string matching the regexp (if not provided, the string is not stored)

You may use named capturing groups in the regexp. In that case, you simply have th add a parameter with the group name prefixed with `store-`. See the [TeamViewer]({% link _recipes/teamviewer.md %}) recipe for an example.

## GET

When executing a GET request, the following parameter is mandatory:

- `url`: URL to request with the GET method

The following parameter is optional:

- `store`: key name in the store for the response (if not provided, the response is not stored)

## POST

When executing a POST request, the two following parameters are mandatory:

- `url`: URL to request with the POST method
- `payload`: the payload to send as the body of the POST method

The following parameter is optional:

- `store`: key name in the store for the response (if not provided, the response is not stored)
