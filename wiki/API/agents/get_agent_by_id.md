# Get Registered Agent by ID

return Information about a Agent by it's ID

> [!WARNING]
> Sample Information - WIP

## URL

```GET https://instance-url.com/v1/agents```

## Authorization

Requires a [Access-Token][access-token]

[access-token]: access-token

## Request Query Parameters

|Parameter|Type|Required|Description|
|---------|----|--------|-----------|
|AgentID|String|Yes|The unique AgentID|

## Response Body

|Field|Type|Description|
|-----|----|-----------|
|name|String|The Name of the Agent|

## Response Codes

|Code|Description|
|----|-----------|
|200|OK|
