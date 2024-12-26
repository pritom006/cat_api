# CatAPI Project

CatAPI is a web service built with Go and the Beego framework. It allows users to interact with cat images and voting functionality. The project integrates with [The Cat API](https://thecatapi.com/) for fetching cat images, voting on them, and adding them to a favorites list.

## Table of Contents

1. [Project Setup](#project-setup)
    - [Installation Steps](#installation-steps)
2. [API Endpoints](#api-endpoints)
3. [Testing](#testing)


## Project Setup

Follow these steps to set up and run the project locally.

### Installation Steps

1. **Create a folder for the Go workspace:**
    ```bash
    mkdir go
    ```

2. **Inside the `go` folder, create another folder called `src`:**
    ```bash
    mkdir go/src
    ```

3. **In the `go` directory, run the following commands:**
    ```bash
    go get github.com/beego/beego/v2@latest
    go install github.com/beego/bee/v2@latest
    ```

4. **Change to the `src` directory:**
    ```bash
    cd go/src
    ```

    - To check the Bee version, run:
      ```bash
      bee version
      ```

    - To create a new Beego project, run:
      ```bash
      bee new catapigo
      ```

5. **Add Go path:**
    - Change to the `go` directory:
      ```bash
      cd go
      ```

    - Edit the `~/.bashrc` file:
      ```bash
      nano ~/.bashrc
      ```

    - Add the following lines to the file:
      ```bash
      export GOPATH=$HOME/go
      export PATH=$PATH:$GOPATH/bin
      ```

    - Save and close the file, then run:
      ```bash
      source ~/.bashrc
      ```

6. **Install project dependencies:**
    ```bash
    go mod tidy
    ```

7. **Run the application:**
    ```bash
    bee run
    ```

---

Alternatively, if you want to clone the repository:

**Clone the repository:**

   ```bash
   git clone https://github.com/pritom006/cat_api.git
   cd catapigo
 ```

To install the required dependencies for the project, follow these steps:

1. Run `go mod tidy` to clean up and synchronize your Go modules:
    ```bash
    go mod tidy
    ```

2. Install the Beego framework:
    ```bash
    go get github.com/beego/beego/v2@latest
    ```

3. Install the Bee tool, which is used to manage the Beego application:
    ```bash
    go install github.com/beego/bee/v2@latest
    ```

These commands will install the necessary dependencies and tools to run the project.

### Set up The Cat API Key

To interact with The Cat API, you need an API key. Follow these steps:

1. Go to [The Cat API](https://thecatapi.com/) and sign up to get your API key.
2. In the `app.conf` file located at the root of the project directory, add the following line with your API key:
    ```bash
    CATAPI_KEY=your-api-key-here
    ```

### Run the Application

Once you've set up your API key, you can run the application using the Bee tool:

```bash
bee run
```

## API Endpoints

### 1. **GET /**

- **Description**: Serves the frontend of the application.
- **Handler**: `ServeFrontend` in `MainController`
- **Method**: GET
- **Response**: Returns the frontend content (HTML, JS, etc.)

### 2. **GET /fetch-breeds**

- **Description**: Fetches the list of cat breeds.
- **Handler**: `FetchCatBreeds` in `MainController`
- **Method**: GET
- **Response**: A JSON array containing cat breeds information.

### 3. **POST /vote**

- **Description**: Submits a vote for a cat image (like/dislike).
- **Handler**: `VoteForCat` in `MainController`
- **Method**: POST
- **Request Body**: 
    ```json
    {
      "ImageID": "string",   // The ID of the image to vote on
      "Value": 1             // Vote value, 1 for like, 0 for dislike
    }
    ```
- **Response**: A JSON object with a message confirming the vote submission or an error message.

### 4. **GET /favorites**

- **Description**: Fetches the list of favorite cat images.
- **Handler**: `FetchFavorites` in `MainController`
- **Method**: GET
- **Response**: A JSON array containing the user's favorite cat images.

### 5. **POST /addToFavorites**

- **Description**: Adds a cat image to the user's favorites.
- **Handler**: `AddToFavorites` in `MainController`
- **Method**: POST
- **Request Body**:
    ```json
    {
      "ImageID": "string"    // The ID of the image to add to favorites
    }
    ```
- **Response**: A JSON object confirming the image was added to favorites or an error message.

### 6. **GET /fetch-new-cat**

- **Description**: Fetches a new random cat image.
- **Handler**: `FetchNewCatImage` in `MainController`
- **Method**: GET
- **Response**: A JSON object containing the URL of the new cat image:
    ```json
    {
      "url": "string" // The URL of the new cat image
    }
    ```

### 7. **GET /fetch-breed-images**

- **Description**: Fetches images of cats of a specific breed.
- **Handler**: `FetchBreedImages` in `MainController`
- **Method**: GET
- **Query Parameter**:
    - `breed_id`: The breed ID of the cat to fetch images for.
- **Response**: A JSON array containing cat images for the specified breed.


## Testing

To run tests for this project, you can use Go's testing tools. Below are the commands for different test scenarios.

### 1. **Run Tests for Controllers**

To run tests for the controllers, execute the following command:

```bash
go test ./tests -v
```

### Run Tests with Coverage

To run tests with code coverage, use the following commands:

```bash
go test -coverprofile coverage.out ./...
go tool cover -html coverage.out
```

### Run Tests for Routers

To run tests specifically for the routers, use the following command:

```bash
go test -coverprofile=coverage.out ./routers/...
go tool cover -html=coverage.out
