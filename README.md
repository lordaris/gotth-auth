# GOTTH Stack stateful authentication boilerplate

A web application boilerplate with stateful authentication using the GOTTH (golang, templ, tailwind, HTMX) stack:

- **Go**
- **Chi** Lightweight and idiomatic router for Go
- **Templ** - Type-safe templates for Go
- **Tailwind CSS** - CSS framework
- **HTMX** - Modern frontend interactivity without JavaScript frameworks
- **PostgreSQL + SQLx** - Database and SQL handling
- **Air** - Live reload for development

## Requirements

Ensure you have the following installed before starting:

- **Go 1.20+**
- **Node.js and npm**
- **PostgreSQL**

## Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/lordaris/gotth-boilerplate
   cd gotth-boilerplate
   ```

2. **Install Go dependencies**

   ```bash
   go mod tidy
   ```

3. **Install Node dependencies**

   ```bash
   npm install
   ```

4. **Configure the PostgreSQL database**
   - Set up your database and update the connection details in the environment configuration.
   - Create a .env file in the main folder and add the following content:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gotth_boilerplate

# Server Configuration
PORT=8080
```

This .env file is essential for storing sensitive configuration values and should not be committed to version control.

5. **Run database migrations**

```bash
make db-migrate
```

6. **Start the development server**

```bash
   make dev
```

7. Visit `http://localhost:8080` in your browser.

## Makefile Commands

This project includes a `Makefile` to simplify common operations, just remove the "copy" text from the file name and modify the project and database variables to your own:

```bash
# Install required tools
make install-tools

# Start development environment (Air + Tailwind)
make dev

# Database operations
make db-create     # Create the database
make db-migrate    # Run migrations
make db-reset      # Reset the database

# Generate Templ files
make generate-templ

# See all available commands
make help
```

Stateful Login System

In this boilerplate, the login system has been upgraded to use a stateful approach, meaning sessions are maintained with cookies on the client-side. Upon logging in, a session token is stored in an HTTP-only cookie, which is used for subsequent requests to authenticate the user. This is a more secure approach compared to stateless systems like JWTs because the session token remains on the client and is tied to a session on the server-side.
Improvements with State Management:

Session management: Instead of relying on tokens passed with each request, the server maintains sessions tied to cookies.
Logout functionality: Logs out the user by clearing the session cookie, ensuring no leftover sessions are active.
Redirect after logout: After logging out, users are automatically redirected to the home page or the login page as configured.

Session Flow:
 Login: When a user logs in, the server creates a session and stores the session ID in an HTTP-only cookie.
Authentication: For every request, the server retrieves the session cookie to authenticate the user.
Logout: The session cookie is cleared, and the user is logged out.


## Improvement Areas & Future Plans

- **Better HTMX usage**: As I'm still learning HTMX, the implementation isn't perfect. Expect improvements in the future.
- **General refinements**: UI, error handling, and overall experience will be enhanced over time.
- Add CRUD elements

For now, if you modify a `templ` file, run:

```bash
make generate-templ
```

and restart the server using:

```bash
make run  # or make dev
```

---

Feel free to improve on this template and make it your own. Happy coding!
