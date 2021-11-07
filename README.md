# Go Code Louisville, Go!

## Project Summary
This project is a simple API that manages the storage of files in an AWS S3 bucket. 
It abstracts away the details of setting up the AWS session, managing the file transactions, etc. by exposing just "upload" and "get" endpoints. 
It could be used as part of a larger application that needs to store, for example, a lot of images but doesn't want to be in the file-hosting business.
This endpoint would serve as a middleman that offloads that work to AWS.

## Installation and Setup
This project is built  almost entirely on the standard Go library, and the only outside dependency is the AWS SDK. 
This should be pulled down for you automatically when you build the project; alternatively you can `go get ./...` to pull them down manually.
No database is required for this project, though an obvious feature would be to persist some data about the files that have been uploaded. 
This could allow us to seemlessly handle potential filename collisions within the bucket and offer a more robust upload mechanism.

## Usage
- **Password Required**: To avoid posting plaintext AWS IDs and secrets to a public GitHub repo, they have been encrypted in the source code and you will be prompted for a password when you run the program. That password is "CodeLouisville" without the quotes. The program will terminate if you type an incorrect password.
- **API Endpoints**: The server will be listening on `localhost:5050` and has two endpoints:
    - `/upload` A POST request with a form-data body. It must have a form field called "file" with the file data to upload as the value.
    The response body will contain a JSON with a succes boolean field as well as other diagnostic information.
    Example:
    ```
  curl --location --request POST 'http://localhost:5050/upload' \ 
  --form 'file=@"/path/to/file/upload.png"'
  ```
    - `\get\{fileToGet}` A GET request with no required body. The filename you want to retrieve should be included as part of the URL path.
    The response body will contain the request file data (if found; if not, you'll get a 404 response status code).
    Example:
    ```
  curl --location --request GET 'http://localhost:5050/get/upload.png'
    ```
- **Please be kind**: This is an active AWS S3 bucket, so please don't upload gigantic files. Not only will the request run for a very long time while it uploads to AWS, I'll also get charged for the data transfer. It's not much, but still...
Also, it would be fairly trivial to extract the unencrypted AWS S3 credentials; but they don't have much in the way of permissions, so it wouldn't be worth the effort.
- **Performance**: The speed of the response will largely depend on the speed of your internet uplink.
An upload of a 1MB file to this endpoint took about 1 second on my connection. YMMV.

## Project Requirements
- "Connect to an external/3rd party API and read data into your app"
    - This project connects to the AWS API via their Go SDK.
    It takes the response from that API to make certain decisions, such as to check whether an upload would cause a filename collision.
    All code that interacts with AWS is in `/server/aws.go`.
- "Create 3 or more unit tests for your application"
    - There is unit test coverage for most of the major components of this program.
    They are located in `aws_test.go`, `encryption_test.go`, `handlers_test.go`, and `server_test.go`.
- **TODO** Figure out how to squeeze existing code into one of these requirements or add something goofy to satisfy them.