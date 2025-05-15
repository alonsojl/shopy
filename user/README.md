## User Lambda Function
This Lambda function manages users in the application.

## Configuration
The following table lists the environment variables required by this function. For example values, check the `.env.example` file.

| Name        | Type   | Description                                                                 |
|-------------|--------|-----------------------------------------------------------------------------|
| TOKEN_KEY   | STRING | Secret key for generating the token after a successful login.               |
| TOKEN_EXP   | INT    | Number of hours after which the token expires.                              |