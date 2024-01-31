# Gallery App

The Gallery App is a web application written in Go that provides basic functionality for user authentication, password reset via email, and the ability for users to create galleries and upload images to them.

## Features

- User Authentication: Users can sign up, log in, and securely authenticate themselves.
- Password Reset via Email: Users can request a password reset via email, providing a secure way to regain access to their accounts.
- Gallery Creation: Authenticated users can create galleries to organize and showcase their images.
- Image Upload: Users can upload images to their galleries to share with others.

## Technologies Used

- Go: The backend is written in Go with server-side rendering.
- Database: PostgreSQL is used to store user information, images are stored on a filesystem.
- Email Service: An email service is used to send password reset links to users.

## Getting Started

To run the Gallery App locally, follow these steps:

1. Clone the repository.
2. Setup .env file with your SMTP account credential. Example is provided in .env.template file.
3. TODO

