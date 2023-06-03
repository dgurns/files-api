```sh
curl -X POST -H "Content-Type: multipart/form-data" \
	-F "file=@/Users/dangurney/Downloads/EngineerChallenge.pdf" \
	-F "metadata={\"source\":\"email\",\"author\":\"inscribe\"}" \
	-u "demo:password" \
	http://localhost:8080/files/upload
```
