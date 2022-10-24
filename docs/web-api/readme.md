# 0. Specification

## 0.1 http protocols

1. Frontend `POST``application/json`
2. Encoding: utf-8
3. Response includes `code`,`message`,`sever_time`,`data`

## 0.2 response

### description for common results

fieldName | type    | comment
------ |---------| ----
`code` | int     | 200: response success
`message` | string  | 
`sever_time` | string  | unixmills
`data` | jsonObj | response data

## 0.3 required params

### **attention**

- case sensitive,if a parameter is empty, make `string` "", make `int` `0

params | type | comment
------ | -------- | -------
`trace_id` | string(32) | unique trace_id ,for tracking
`access_token` | string(32) | after login,the access_token is required;
`request_time` | string | request_time for each request
`app_sn` | int | client type

# `web-api` 

## 1. login

`/web-api/oauth-login`

### params

params | type   | required | comment
------ |--------|----------| ---
`open_oauth_plt` | int    | yes      | 1: wechat platform
`open_user_id` | string | yes      | unique user_id for the platform
`nickname` | string | no       | 
`open_avatar` | string | no       | header
`open_token` | string | no       | third party platform token
`email` | string | no       | 

### response example

```json
{
  "code": 200,
  "message": "",
  "sever_time": "1599121007098",
  "data": {
    "access_token": "4e6d7b1b2301d8b6001c72c075f4c5de"
  }
}
```

