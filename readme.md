# File Upload server

A simple server for upload file and return a hashed url

## Upload

```
let bodyContent = new FormData();
bodyContent.append("file", "<file_path>");

let response = await fetch("localhost:8081/upload", {
  method: "POST",
  body: bodyContent,
  headers: headersList
});

```
