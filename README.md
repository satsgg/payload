# Payload - Video Streaming Platform

A modern video streaming platform built with React, TypeScript, and Go. This application provides a robust solution for video hosting and streaming with an admin interface for content management.

## Features

- 🎥 Video streaming with HLS support
- 🔒 Secure authentication system
- 👤 Admin dashboard for content management
- 📱 Responsive design with modern UI components
- 🎨 Built with Tailwind CSS and Radix UI
- 🔄 Real-time video streaming
- 📊 Video metadata management

## Tech Stack

### Frontend

- React 18
- TypeScript
- Vite
- Tailwind CSS
- ShadCN
- Media-Chrome
- React Router
- React Hook Form
- Zod for validation

### Backend

- Go (Golang)
- Gin Web Framework
- SQLite Database
- JWT Authentication
- FFmpeg

## Prerequisites

- Node.js (v18 or higher)
- Go (v1.20 or higher)
- SQLite3

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd payload
```

2. Install frontend dependencies:

```bash
npm install
```

3. Install Go dependencies:

```bash
go mod download
```

## Development

1. Start the frontend development server:

```bash
npm run dev
```

2. Start the backend server:

```bash
npm run start:backend
```

The application will be available at `http://localhost:5173/`.

Log into the admin panel at `http://localhost:5173/admin` to get started uploading videos.

## Building for Production

1. Build the frontend:

```bash
npm run build
```

2. Start the production server:

```bash
npm run start
```

## Configuration

The backend can be configured using command-line flags:

- `--port`: Server port (default: 8080)
- `--db`: SQLite database path (default: videos.db)
- `--videos`: Video storage path (default: videos)
- `--lnd-host`: LND gRPC host (default: localhost:10009)
- `--tls-cert`: Path to LND TLS certificate
- `--macaroon`: Path to LND macaroon file

## Default Admin Account

A default admin account is created during database initialization:

- Username: admin
- Password: admin123

## API Endpoints

### Public Endpoints

- `GET /api/videos` - List all videos
- `GET /api/videos/:id` - Get video details
- `GET /api/videos/:id/*filepath` - Stream video content

### Admin Endpoints

- `POST /api/admin/login` - Admin login
- `GET /api/admin/verify` - Verify admin token
- `POST /api/admin/videos` - Upload new video
- `PUT /api/admin/videos/:id` - Update video
- `DELETE /api/admin/videos/:id` - Delete video

## Security

- JWT-based authentication
- Password hashing with bcrypt
- Secure video streaming
- Admin-only access to sensitive operations

## License

This project is licensed under the terms of the license included in the repository.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
